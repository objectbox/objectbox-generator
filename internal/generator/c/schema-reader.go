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

	"github.com/objectbox/objectbox-go/internal/generator/fbsparser/reflection"
	"github.com/objectbox/objectbox-go/internal/generator/model"
)

// fbSchemaReader reads FlatBuffers schema and populates a model
type fbSchemaReader struct {
	// model produced by reading the schema
	model *model.ModelInfo
}

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
	entity.Name = string(object.Name())

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
		return entity.Properties[i].Meta.(*fbsProperty).fbsField.Id() < entity.Properties[j].Meta.(*fbsProperty).fbsField.Id()
	})

	r.model.Entities = append(r.model.Entities, entity)
	return nil
}

func (r *fbSchemaReader) readObjectField(entity *model.Entity, field *reflection.Field) error {
	var property = model.CreateProperty(entity, 0, 0)
	property.Name = string(field.Name())
	property.Meta = &fbsProperty{property, field}

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

		// apply flags defined for this type (e.g. unsigned)
		property.Flags = property.Flags | fbsTypeToObxFlag[fbsBaseType]
	}

	entity.Properties = append(entity.Properties, property)
	return nil
}
