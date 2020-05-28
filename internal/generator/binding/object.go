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

	"github.com/objectbox/objectbox-go/internal/generator/go"
	"github.com/objectbox/objectbox-go/internal/generator/model"
)

// Object holds common entity information used by specialized code parsers/generators.
// Additionally, it groups some shared logic, e.g. annotation processing
type Object struct {
	Name      string
	Namespace string
	IsSkipped bool

	entity *model.Entity
}

func CreateObject(entity *model.Entity) *Object {
	return &Object{entity: entity}
}

func (object *Object) SetName(name string) {
	// look for namespace separators
	var lastDot = strings.LastIndex(name, ".")
	if lastDot > 0 {
		object.Namespace = name[:lastDot]
		name = name[lastDot+1:]
	}

	object.Name = name
	if len(object.entity.Name) == 0 {
		object.entity.Name = strings.ToLower(name)
	}
}

// ProcessAnnotations checks all set annotations for any inconsistencies and sets local/entity properties (uid, name, ...)
// TODO move generator.Annotation to this package
func (object *Object) ProcessAnnotations(a map[string]*gogenerator.Annotation) error {
	if a["-"] != nil {
		if len(a) != 1 || a["-"].Value != "" {
			return errors.New("to ignore the entity, use only `objectbox:\"-\"` as an annotation")
		}
		object.IsSkipped = true
		return nil
	}

	if a["name"] != nil {
		if len(a["name"].Value) == 0 {
			return fmt.Errorf("name annotation value must not be empty - it's the entity name in DB")
		}
		object.entity.Name = strings.ToLower(a["name"].Value)
	}

	if a["uid"] != nil {
		if len(a["uid"].Value) == 0 {
			// in case the user doesn't provide `objectbox:"uid"` value, it's considered in-process of setting up UID
			// this flag is handled by the merge mechanism and prints the UID of the already existing entity
			object.entity.UidRequest = true
		} else if uid, err := strconv.ParseUint(a["uid"].Value, 10, 64); err != nil {
			return fmt.Errorf("can't parse uid - %s", err)
		} else if id, err := object.entity.Id.GetIdAllowZero(); err != nil {
			return fmt.Errorf("can't parse entity Id - %s", err)
		} else {
			object.entity.Id = model.CreateIdUid(id, uid)
		}
	}

	return nil
}
