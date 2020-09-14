/*
 * Copyright (C) 2020 ObjectBox Ltd. All rights reserved.
 * https://objectbox.io
 *
 * This file is part of ObjectBox Generator.
 *
 * ObjectBox Generator is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * ObjectBox Generator is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with ObjectBox Generator.  If not, see <http://www.gnu.org/licenses/>.
 */

package cgenerator

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/objectbox/objectbox-generator/internal/generator/binding"
	"github.com/objectbox/objectbox-generator/internal/generator/flatbuffersc/reflection"
	"github.com/objectbox/objectbox-generator/internal/generator/model"
)

var supportedEntityAnnotations = map[string]bool{
	"transient": true,
	"name":      true,
	"uid":       true,
	"relation":  true, // to-many, standalone
}

var supportedPropertyAnnotations = map[string]bool{
	"transient":    true,
	"date":         true,
	"date-nano":    true,
	"id":           true,
	"index":        true,
	"relation":     true, // to-one
	"name":         true,
	"uid":          true,
	"unique":       true,
	"id-companion": true,
}

// fbSchemaReader reads FlatBuffers schema and populates a model
type fbSchemaReader struct {
	// model produced by reading the schema
	model *model.ModelInfo
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
		return true, parseAnnotations(comment, annotations, supportedAnnotations)
	}
	return false, nil
}

type annotationInProgress struct {
	name          string
	key           string
	value         *binding.Annotation
	valueQuoted   bool
	valueFinished bool
}

func (s *annotationInProgress) finishAnnotation(annotations *map[string]*binding.Annotation, supportedAnnotations map[string]bool) error {
	s.name = strings.TrimSpace(s.name)
	if len(s.name) == 0 {
		return nil
	}
	if s.value == nil {
		s.value = &binding.Annotation{} // empty value
	} else {
		s.value.Value = strings.TrimSpace(s.value.Value)
	}
	var key = s.key
	if len(key) == 0 {
		key = s.name
	}
	if (*annotations)[key] != nil {
		return fmt.Errorf("duplicate annotation %s", key)
	} else if !supportedAnnotations[s.name] {
		return fmt.Errorf("unknown annotation %s", s.name)
	} else {
		(*annotations)[key] = s.value
	}
	return nil
}

// counts all "relation-" prefixed annotations (standalone relations) - used to ensure consistent processing order
func relationsCount(annotations map[string]*binding.Annotation) uint {
	var count uint
	for key := range annotations {
		if strings.HasPrefix(key, "relation-") {
			count++
		}
	}
	return count
}

// parseAnnotations parses annotations in any of the following formats.
// name="name",index - creates two annotations, name and index, the former having a non-empty value
// relation(name=manyToManyRelName,to=TargetEntity) - creates a single annotation relation with two items as details
// id - creates a single annotation
// NOTE: this started as a very simple parser but it seems like the requirements are ever-increasing... maybe some form
//       of recursive tokenization would be better in case we decided to rework.
func parseAnnotations(str string, annotations *map[string]*binding.Annotation, supportedAnnotations map[string]bool) error {
	var s annotationInProgress
	for i := 0; i < len(str); i++ {
		var char = str[i]

		if (char == '=' || char == '(') && !s.valueQuoted { // start a value
			if len(s.name) == 0 {
				return fmt.Errorf("invalid annotation format: name expected before '%s' at position %d in `%s` ", string(char), i, str)
			}
			s.value = &binding.Annotation{}

			// special handling for "recursive" details (many-to-many relation)
			if char == '(' {
				// find the closing bracket
				var detailsStr string
				for j := i + 1; j < len(str); j++ {
					if str[j] == ')' { // NOTE we're ignoring potential closing brackets in quotes
						detailsStr += str[i+1 : j]
						i = j // skip up to this character in the parent loop
						break
					}
				}
				if len(detailsStr) == 0 {
					return fmt.Errorf("invalid annotation details format, closing bracket ')' not found in `%s`", str[i+1:])
				}
				s.name = strings.TrimSpace(s.name)
				if s.name != "relation" {
					return fmt.Errorf("invalid annotation format: details only supported for `relation` annotations, found `%s`", s.name)
				}
				s.value.Details = make(map[string]*binding.Annotation)
				if err := parseAnnotations(detailsStr, &s.value.Details, map[string]bool{"to": true, "name": true, "uid": true}); err != nil {
					return err
				}
				if s.value.Details["name"] == nil {
					return fmt.Errorf("invalid annotation format: relation name missing in `%s`", str)
				}
				s.key = fmt.Sprintf("relation-%10d-%s", relationsCount(*annotations), s.value.Details["name"].Value)
				if err := s.finishAnnotation(annotations, supportedAnnotations); err != nil {
					return err
				}
				s = annotationInProgress{} // reset
			}

		} else if char == ',' && !s.valueQuoted { // finish an annotation
			if err := s.finishAnnotation(annotations, supportedAnnotations); err != nil {
				return err
			}
			s = annotationInProgress{} // reset
		} else if s.value != nil { // continue a value (set contents)
			if char == '"' {
				if len(s.value.Value) == 0 {
					s.valueQuoted = true
				} else {
					s.valueQuoted = false
					s.valueFinished = true
				}
			} else if s.valueFinished {
				return fmt.Errorf("invalid annotation format: no more characters may follow after a quoted value at position %d in `%s`", i, str)
			} else {
				s.value.Value += string(char)
			}
		} else { // continue a name
			s.name += string(char)
		}
	}

	return s.finishAnnotation(annotations, supportedAnnotations)
}
