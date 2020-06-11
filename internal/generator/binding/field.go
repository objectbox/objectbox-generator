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

package binding

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/objectbox/objectbox-generator/internal/generator/model"
)

// Field holds common field/property information used by specialized code parsers/generators.
// Additionally, it groups some shared logic, e.g. annotation processing
type Field struct {
	ModelProperty *model.Property
	Name          string
	IsSkipped     bool
}

func CreateField(prop *model.Property) *Field {
	return &Field{ModelProperty: prop}
}

func (field *Field) SetName(name string) {
	field.Name = name
	if len(field.ModelProperty.Name) == 0 {
		field.ModelProperty.Name = name
	}
}
func (field *Field) PreProcessAnnotations(a map[string]*Annotation) error {
	field.IsSkipped = false
	for _, alternative := range []string{"-", "transient"} {
		if a[alternative] != nil {
			if len(a) != 1 || a[alternative].Value != "" {
				return errors.New("to ignore the property, use only `objectbox:\"" + alternative + "\"` as an annotation")
			}
			field.IsSkipped = true
			return nil
		}
	}
	return nil
}

// ProcessAnnotations checks all set annotations for any inconsistencies and sets local/property fields (flags, name, ...)
// TODO move generator.Annotation to this package
func (field *Field) ProcessAnnotations(a map[string]*Annotation) error {
	if err := field.PreProcessAnnotations(a); err != nil {
		return err
	}

	if field.IsSkipped {
		return nil
	}

	if a["id"] != nil {
		field.ModelProperty.AddFlag(model.PropertyFlagId)
	}

	if a["name"] != nil {
		if len(a["name"].Value) == 0 {
			return fmt.Errorf("name annotation value must not be empty - it's the field name in DB")
		}
		field.ModelProperty.Name = a["name"].Value
	}

	if a["date"] != nil {
		if field.ModelProperty.Type != model.PropertyTypeLong {
			return fmt.Errorf("invalid underlying type '%v' for date field; expecting long", model.PropertyTypeNames[field.ModelProperty.Type])
		}
		field.ModelProperty.Type = model.PropertyTypeDate
	}

	if a["id-companion"] != nil {
		if field.ModelProperty.Type != model.PropertyTypeDate {
			return fmt.Errorf("invalid underlying type '%v' for ID companion field; expecting date", model.PropertyTypeNames[field.ModelProperty.Type])
		}
		field.ModelProperty.AddFlag(model.PropertyFlagIdCompanion)
	}

	if a["index"] != nil {
		switch strings.ToLower(a["index"].Value) {
		case "":
			// if the user doesn't define index type use the default based on the data-type
			if field.ModelProperty.Type == model.PropertyTypeString {
				field.ModelProperty.AddFlag(model.PropertyFlagIndexHash)
			} else {
				field.ModelProperty.AddFlag(model.PropertyFlagIndexed)
			}
		case "value":
			field.ModelProperty.AddFlag(model.PropertyFlagIndexed)
		case "hash":
			field.ModelProperty.AddFlag(model.PropertyFlagIndexHash)
		case "hash64":
			field.ModelProperty.AddFlag(model.PropertyFlagIndexHash64)
		default:
			return fmt.Errorf("unknown index type %s", a["index"].Value)
		}

		if err := field.ModelProperty.SetIndex(); err != nil {
			return err
		}
	}

	if a["unique"] != nil {
		field.ModelProperty.AddFlag(model.PropertyFlagUnique)

		if err := field.ModelProperty.SetIndex(); err != nil {
			return err
		}
	}

	if a["uid"] != nil {
		if len(a["uid"].Value) == 0 {
			// in case the user doesn't provide `objectbox:"uid"` value, it's considered in-process of setting up UID
			// this flag is handled by the merge mechanism and prints the UID of the already existing property
			field.ModelProperty.UidRequest = true
		} else if uid, err := strconv.ParseUint(a["uid"].Value, 10, 64); err != nil {
			return fmt.Errorf("can't parse uid - %s", err)
		} else if id, err := field.ModelProperty.Id.GetIdAllowZero(); err != nil {
			return fmt.Errorf("can't parse property Id - %s", err)
		} else {
			field.ModelProperty.Id = model.CreateIdUid(id, uid)
		}
	}

	// To-one relation
	if a["link"] != nil && field.ModelProperty.Type != model.PropertyTypeRelation {
		if field.ModelProperty.Type != model.PropertyTypeLong {
			return fmt.Errorf("invalid underlying type (PropertyType %v) for relation field; expecting long", model.PropertyTypeNames[field.ModelProperty.Type])
		}
		if len(a["link"].Value) == 0 {
			return errors.New("unknown link target entity, define by changing the `link` annotation to the `link=Entity` format")
		}
		field.ModelProperty.Type = model.PropertyTypeRelation
		field.ModelProperty.RelationTarget = a["link"].Value
		field.ModelProperty.AddFlag(model.PropertyFlagIndexed)
		field.ModelProperty.AddFlag(model.PropertyFlagIndexPartialSkipZero)

		if err := field.ModelProperty.SetIndex(); err != nil {
			return err
		}
	}

	return nil
}
