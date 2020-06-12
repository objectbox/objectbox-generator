/*
 * Copyright 2019 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
}

var supportedPropertyAnnotations = map[string]bool{
	"transient":    true,
	"date":         true,
	"id":           true,
	"index":        true,
	"link":         true,
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
		if isAnnotation, err := parseAnnotations(comment, &annotations, supportedEntityAnnotations); err != nil {
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
		if isAnnotation, err := parseAnnotations(comment, &annotations, supportedPropertyAnnotations); err != nil {
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

// NOTE this is a copy of gogenerator.parseAnnotations with changes to accommodate a different format (not
func parseAnnotations(comment string, annotations *map[string]*binding.Annotation, supportedAnnotations map[string]bool) (bool, error) {
	if strings.HasPrefix(comment, "objectbox:") || strings.HasPrefix(comment, "ObjectBox:") {
		comment = strings.TrimSpace(comment[len("objectbox:"):])
		if len(comment) == 0 {
			return true, nil
		}
	} else {
		return false, nil
	}

	// sample content at this point:
	// name="name",index
	// id

	type state struct {
		name          string
		value         *binding.Annotation
		valueQuoted   bool
		valueFinished bool
	}
	var s state

	var finishAnnotation = func() error {
		s.name = strings.TrimSpace(s.name)
		if s.value == nil {
			s.value = &binding.Annotation{} // empty value
		} else {
			s.value.Value = strings.TrimSpace(s.value.Value)
		}
		if (*annotations)[s.name] != nil {
			return fmt.Errorf("duplicate annotation %s", s.name)
		} else if !supportedAnnotations[s.name] {
			return fmt.Errorf("unknown annotation %s", s.name)
		} else {
			(*annotations)[s.name] = s.value
		}
		s = state{}
		return nil
	}

	for i, char := range comment {
		if char == '=' && !s.valueQuoted { // start a value
			if len(s.name) == 0 {
				return true, fmt.Errorf("invalid annotation format: name must precede equal sign at position %d in `%s` ", i, comment)
			}
			s.value = &binding.Annotation{}
		} else if char == ',' && !s.valueQuoted { // finish an annotation
			if err := finishAnnotation(); err != nil {
				return true, err
			}
		} else if s.value != nil { // continue a value (set contents)
			if char == '"' {
				if len(s.value.Value) == 0 {
					s.valueQuoted = true
				} else {
					s.valueQuoted = false
					s.valueFinished = true
				}
			} else if s.valueFinished {
				return true, fmt.Errorf("invalid annotation format: no more characters may follow after a quoted value at position %d in `%s`", i, comment)
			} else {
				s.value.Value += string(char)
			}
		} else { // continue a name
			s.name += string(char)
		}
	}

	if err := finishAnnotation(); err != nil {
		return true, err
	}

	return true, nil
}
