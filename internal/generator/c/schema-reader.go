/*
 * ObjectBox Generator - a build time tool for ObjectBox
 * Copyright (C) 2018-2024 ObjectBox Ltd. All rights reserved.
 * https://objectbox.io
 *
 * This file is part of ObjectBox Generator.
 *
 * ObjectBox Generator is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * ObjectBox Generator is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with ObjectBox Generator.  If not, see <http://www.gnu.org/licenses/>.
 */

package cgenerator

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/objectbox/objectbox-generator/v4/internal/generator/binding"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/flatbuffersc/reflection"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/model"
)

var supportedEntityAnnotations = map[string]bool{
	"name":          true,
	"relation":      true, // to-many, standalone
	"sync":          true,
	"transient":     true,
	"uid":           true,
	"external-name": true,
}

var supportedPropertyAnnotations = map[string]bool{
	"date":                                 true,
	"date-nano":                            true,
	"id":                                   true,
	"id-companion":                         true,
	"index":                                true,
	"name":                                 true,
	"optional":                             true,
	"relation":                             true, // to-one
	"transient":                            true,
	"uid":                                  true,
	"unique":                               true,
	"hnsw-dimensions":                      true,
	"hnsw-distance-type":                   true,
	"hnsw-neighbors-per-node":              true,
	"hnsw-indexing-search-count":           true,
	"hnsw-flags":                           true,
	"hnsw-reparation-backlink-probability": true,
	"hnsw-vector-cache-hint-size-kb":       true,
	"external-name":                        true,
	"external-type":                        true,
}

// fbSchemaReader reads FlatBuffers schema and populates a model
type fbSchemaReader struct {
	// model produced by reading the schema
	model *model.ModelInfo

	// see CGenerator.Optional
	optional string
}

// const annotationPrefix = "objectbox:"

func (r *fbSchemaReader) read(schema *reflection.Schema) error {
	for i := 0; i < schema.ObjectsLength(); i++ {
		var object reflection.Object
		if !schema.Objects(&object, i) {
			return fmt.Errorf("can't access object %d", i)
		}

		if err := r.readObject(&object); err != nil {
			return fmt.Errorf("object %d %s: %v", i, string(object.Name()), err)
		}
	}

	return nil
}

func (r *fbSchemaReader) readObject(object *reflection.Object) error {
	var entity = model.CreateEntity(r.model, 0, 0)
	var metaEntity = &fbsObject{binding.CreateObject(entity), object}
	entity.Meta = metaEntity
	metaEntity.SetName(string(object.Name()))

	// look for annotations: "/// objectbox:..."
	var annotations = make(map[string]*binding.Annotation)
	for i := 0; i < object.DocumentationLength(); i++ {
		var comment = strings.TrimSpace(string(object.Documentation(i)))
		if isAnnotation, err := parseCommentAsAnnotations(comment, &annotations, supportedEntityAnnotations); err != nil {
			return err
		} else if !isAnnotation {
			entity.Comments = append(entity.Comments, comment)
		}
	}

	if err := metaEntity.ProcessAnnotations(annotations); err != nil {
		return err
	}

	if metaEntity.IsSkipped {
		return nil
	}

	// attach "meta" objects to relations
	for _, rel := range entity.Relations {
		rel.Meta = &standaloneRel{ModelRelation: rel}
	}

	for i := 0; i < object.FieldsLength(); i++ {
		var field reflection.Field
		if !object.Fields(&field, i) {
			return fmt.Errorf("can't access field %d", i)
		}

		if err := r.readObjectField(entity, &field); err != nil {
			return fmt.Errorf("field %d %s: %v", i, string(field.Name()), err)
		}
	}

	// Schema reader provides fields ordered by name but we want them ordered by the order they appear in the input
	// file. While that's not available on reflection.Field, there's an alternative: FlatBufferSchema ID, which is,
	// unless explicitly overridden using an id attribute in the schema, the order in the input file.
	sort.Slice(entity.Properties, func(i, j int) bool {
		return entity.Properties[i].Meta.(*fbsField).fbsField.Id() < entity.Properties[j].Meta.(*fbsField).fbsField.Id()
	})

	r.model.Entities = append(r.model.Entities, entity)
	return nil
}

func (r *fbSchemaReader) readObjectField(entity *model.Entity, field *reflection.Field) error {
	var property = model.CreateProperty(entity, 0, 0)
	var metaProperty = &fbsField{binding.CreateField(property), field}
	property.Meta = metaProperty
	metaProperty.SetName(string(field.Name()))

	// look for annotations: "/// objectbox:..."
	var annotations = make(map[string]*binding.Annotation)
	for i := 0; i < field.DocumentationLength(); i++ {
		var comment = strings.TrimSpace(string(field.Documentation(i)))
		if isAnnotation, err := parseCommentAsAnnotations(comment, &annotations, supportedPropertyAnnotations); err != nil {
			return err
		} else if !isAnnotation {
			property.Comments = append(property.Comments, comment)
		}
	}

	if err := metaProperty.PreProcessAnnotations(annotations); err != nil {
		return err
	}

	if metaProperty.IsSkipped {
		return nil
	}

	if fbsType := field.Type(nil); fbsType == nil {
		return errors.New("can't access Type() from the source schema")
	} else {
		var fbsBaseType = fbsType.BaseType()
		if fbsBaseType == reflection.BaseTypeVector {
			var fbsElBaseType = fbsType.Element()
			switch fbsElBaseType {
			case reflection.BaseTypeString:
				property.Type = model.PropertyTypeStringVector
			case reflection.BaseTypeByte:
				fallthrough
			case reflection.BaseTypeUByte:
				property.Type = model.PropertyTypeByteVector
			case reflection.BaseTypeFloat:
				property.Type = model.PropertyTypeFloatVector
			default:
				return fmt.Errorf("unsupported vector element type: %s", reflection.EnumNamesBaseType[fbsElBaseType])
			}
		} else {
			property.Type = fbsTypeToObxType[fbsBaseType]
		}

		if property.Type == 0 {
			return fmt.Errorf("unsupported type: %s", reflection.EnumNamesBaseType[fbsBaseType])
		}

		// apply flags defined for this type (e.g.
		property.AddFlag(fbsTypeToObxFlag[fbsBaseType])
	}

	if annotations["optional"] != nil {
		if len(annotations["optional"].Value) != 0 {
			return errors.New("optional annotation value must be empty")
		}
		annotations["optional"].Value = r.optional
	}

	if err := metaProperty.ProcessAnnotations(annotations); err != nil {
		return err
	}

	entity.Properties = append(entity.Properties, property)
	return nil
}

// NOTE this is a copy of gogenerator.parseAnnotations with changes to accommodate a different format
func parseCommentAsAnnotations(comment string, annotations *map[string]*binding.Annotation, supportedAnnotations map[string]bool) (bool, error) {
	if strings.HasPrefix(comment, "objectbox:") || strings.HasPrefix(comment, "ObjectBox:") {
		comment = strings.TrimSpace(comment[len("objectbox:"):])
		if len(comment) == 0 {
			return true, nil
		}
		return true, binding.ParseAnnotations(comment, annotations, supportedAnnotations)
	}
	return false, nil
}
