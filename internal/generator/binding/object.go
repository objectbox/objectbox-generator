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

package binding

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/objectbox/objectbox-generator/v4/internal/generator/model"
)

// Object holds common entity information used by specialized code parsers/generators.
// Additionally, it groups some shared logic, e.g. annotation processing
type Object struct {
	ModelEntity *model.Entity
	Name        string
	Namespace   string
	IsSkipped   bool
}

func CreateObject(entity *model.Entity) *Object {
	return &Object{ModelEntity: entity}
}

func (object *Object) SetName(name string) {
	// look for namespace separators
	var lastDot = strings.LastIndex(name, ".")
	if lastDot > 0 {
		object.Namespace = name[:lastDot]
		name = name[lastDot+1:]
	}

	object.Name = name
	if len(object.ModelEntity.Name) == 0 {
		object.ModelEntity.Name = name
	}
}

// ProcessAnnotations checks all set annotations for any inconsistencies and sets local/entity properties (uid, name, ...)
func (object *Object) ProcessAnnotations(a map[string]*Annotation) error {
	for _, alternative := range []string{"-", "transient"} {
		if a[alternative] != nil {
			if len(a) != 1 || a[alternative].Value != "" {
				return errors.New("to ignore the entity, use only `objectbox:\"" + alternative + "\"` as an annotation")
			}
			object.IsSkipped = true
			return nil
		}
	}

	if a["name"] != nil {
		if len(a["name"].Value) == 0 {
			return fmt.Errorf("name annotation value must not be empty - it's the entity name in DB")
		}
		object.ModelEntity.Name = a["name"].Value
	}

	if a["sync"] != nil {
		if len(a["sync"].Value) != 0 {
			return fmt.Errorf("sync annotation value must be empty")
		}
		object.ModelEntity.AddFlag(model.EntityFlagSyncEnabled)
		if hasDetail, err := HasBooleanDetail(a, "sync", "sharedglobalids"); err != nil {
			return err
		} else if hasDetail {
			object.ModelEntity.AddFlag(model.EntityFlagSharedGlobalIds)
		}
	}

	if a["external-name"] != nil {
		object.ModelEntity.ExternalName = a["external-name"].Value
	}

	if a["uid"] != nil {
		if len(a["uid"].Value) == 0 {
			// in case the user doesn't provide `objectbox:"uid"` value, it's considered in-process of setting up UID
			// this flag is handled by the merge mechanism and prints the UID of the already existing entity
			object.ModelEntity.UidRequest = true
		} else if uid, err := strconv.ParseUint(a["uid"].Value, 10, 64); err != nil {
			return fmt.Errorf("can't parse uid - %s", err)
		} else if id, err := object.ModelEntity.Id.GetIdAllowZero(); err != nil {
			return fmt.Errorf("can't parse entity Id - %s", err)
		} else {
			object.ModelEntity.Id = model.CreateIdUid(id, uid)
		}
	}

	// Always process standalone relations in the same order by gathering the keys and sorting them, instead of relying
	// on the random order of map keys. We're doing this to avoid unintended order changes in the generated code/model.
	var relationKeys []string
	for key := range a {
		if strings.HasPrefix(key, "relation-") {
			relationKeys = append(relationKeys, key)
		}
	}
	sort.Strings(relationKeys)
	for _, key := range relationKeys {
		if _, err := object.AddRelation(a[key].Details); err != nil {
			return err
		}
	}

	return nil
}

func (object *Object) AddRelation(details map[string]*Annotation) (*model.StandaloneRelation, error) {
	var relation = model.CreateStandaloneRelation(object.ModelEntity, model.CreateIdUid(0, 0))
	if details["name"] == nil || len(details["name"].Value) == 0 {
		return nil, fmt.Errorf("name annotation value must not be empty on relation - it's the relation name in DB")
	}
	relation.Name = details["name"].Value

	if details["to"] == nil || len(details["to"].Value) == 0 {
		return nil, fmt.Errorf("to annotation value must not be empty on relation %s - specify target entity", relation.Name)
	}

	if details["external-name"] != nil {
		relation.ExternalName = details["external-name"].Value
	}

	if details["external-type"] != nil {
		externalTypeValue := details["external-type"].Value
		externalType, exists := model.ExternalTypeValues[externalTypeValue]
		if !exists {
			return nil, fmt.Errorf("invalid external-type '%s' for relation %s - must be one of the supported external types", externalTypeValue, relation.Name)
		}
		relation.ExternalType = externalType
	}

	// NOTE: we don't need an actual entity pointer, it's resolved during stored model merging.
	relation.Target = &model.Entity{Name: details["to"].Value}

	if details["uid"] != nil {
		if len(details["uid"].Value) == 0 {
			// in case the user doesn't provide `objectbox:"uid"` value, it's considered in-process of setting up UID
			// this flag is handled by the merge mechanism and prints the UID of the already existing entity
			relation.UidRequest = true
		} else if uid, err := strconv.ParseUint(details["uid"].Value, 10, 64); err != nil {
			return nil, fmt.Errorf("can't parse uid on relation %s - %s", relation.Name, err)
		} else {
			relation.Id = model.CreateIdUid(0, uid)
		}
	}
	object.ModelEntity.Relations = append(object.ModelEntity.Relations, relation)
	return relation, nil
}
