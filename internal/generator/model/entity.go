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

package model

import (
	"errors"
	"fmt"
	"strings"
)

// Entity represents a DB entity
type Entity struct {
	Id               IdUid                 `json:"id"`
	LastPropertyId   IdUid                 `json:"lastPropertyId"`
	Name             string                `json:"name"`
	ExternalName     string                `json:"externalName,omitempty"`
	Flags            EntityFlags           `json:"flags,omitempty"`
	Properties       []*Property           `json:"properties"`
	Relations        []*StandaloneRelation `json:"relations,omitempty"`
	UidRequest       bool                  `json:"-"` // used when the user gives an empty uid annotation
	Meta             EntityMeta            `json:"-"`
	CurrentlyPresent bool                  `json:"-"`
	Comments         []string              `json:"-"`
	Model            *ModelInfo            `json:"-"`
}

// CreateEntity constructs an Entity
func CreateEntity(model *ModelInfo, id Id, uid Uid) *Entity {
	return &Entity{
		Model:      model,
		Id:         CreateIdUid(id, uid),
		Properties: make([]*Property, 0),
	}
}

// Validate performs validation of the entity model
func (entity *Entity) Validate() (err error) {
	if entity.Model == nil {
		return fmt.Errorf("undefined parent model")
	}

	if err = entity.Id.Validate(); err != nil {
		return err
	}

	if len(entity.Name) == 0 {
		return fmt.Errorf("name is undefined")
	}

	if len(entity.Properties) > 0 {
		if err = entity.LastPropertyId.Validate(); err != nil {
			return fmt.Errorf("lastPropertyId: %s", err)
		}

		var lastId = entity.LastPropertyId.getIdSafe()
		var lastUid = entity.LastPropertyId.getUidSafe()

		var propertiesByName = make(map[string]bool)

		var found = false
		for _, property := range entity.Properties {
			// ObjectBox core internally converts to lowercase so we should check it as this as well
			var realName = strings.ToLower(property.Name)
			if propertiesByName[realName] {
				return fmt.Errorf("duplicate property name '%s' (note that property names are case insensitive)", property.Name)
			}
			propertiesByName[realName] = true

			if property.Entity == nil {
				property.Entity = entity
			} else if property.Entity != entity {
				return fmt.Errorf("relation %s %s has incorrect parent entity reference",
					property.Name, property.Id)
			}

			if lastId == property.Id.getIdSafe() {
				if lastUid != property.Id.getUidSafe() {
					return fmt.Errorf("lastPropertyId %s doesn't match property %s %s",
						entity.LastPropertyId, property.Name, property.Id)
				}
				found = true
			} else if lastId < property.Id.getIdSafe() {
				return fmt.Errorf("lastPropertyId %s is lower than relation %s %s",
					entity.LastPropertyId, property.Name, property.Id)
			}
		}

		if !found && !searchSliceUid(entity.Model.RetiredPropertyUids, lastUid) {
			return fmt.Errorf("lastPropertyId %s doesn't match any relation", entity.LastPropertyId)
		}
	}

	if entity.Properties == nil {
		return fmt.Errorf("properties are not defined or not an array")
	}

	var idProp *Property
	for _, property := range entity.Properties {
		err = property.Validate()
		if err != nil {
			return fmt.Errorf("property %s %s is invalid: %s", property.Name, string(property.Id), err)
		}
		if property.IsIdProperty() {
			if idProp != nil {
				return fmt.Errorf("multiple properties marked as ID: %s (%s) and %s (%s)",
					idProp.Name, idProp.Id, property.Name, property.Id)
			}
			idProp = property
		}
	}

	for _, relation := range entity.Relations {
		if relation.entity == nil {
			relation.entity = entity
		} else if relation.entity != entity {
			return fmt.Errorf("relation %s %s has incorrect parent model reference", relation.Name, relation.Id)
		}

		err = relation.Validate()
		if err != nil {
			return fmt.Errorf("relation %s %s is invalid: %s", relation.Name, string(relation.Id), err)
		}
	}

	return nil
}

func (entity *Entity) finalize() error {
	for _, property := range entity.Properties {
		if err := property.finalize(); err != nil {
			return err
		}
	}
	if err := entity.AutosetIdProperty(nil); err != nil {
		return err
	}
	return entity.Validate()
}

func (entity *Entity) getIdProperty() *Property {
	for _, property := range entity.Properties {
		if property.IsIdProperty() {
			return property
		}
	}
	return nil
}

// AutosetIdProperty updates finds a property that's defined as an ID and if none is, tries to set one based on its name and type
func (entity *Entity) AutosetIdProperty(acceptedTypes []PropertyType) error {
	if entity.getIdProperty() == nil {
		// try to find an ID property automatically based on its name and type
		var idProp *Property
		for _, property := range entity.Properties {
			if strings.ToLower(property.Name) == "id" && property.hasValidTypeAsId(acceptedTypes) {
				if idProp != nil {
					return fmt.Errorf("multiple properties recognized as an ID: %s (%s) and %s (%s)",
						idProp.Name, idProp.Id, property.Name, property.Id)
				}
				idProp = property
			}
		}
		if idProp == nil {
			return errors.New("no property recognized as an ID")
		}

		idProp.Flags = idProp.Flags | PropertyFlagId

		// IDs must not be tagged unsigned for compatibility reasons
		idProp.Flags = idProp.Flags & ^PropertyFlagUnsigned
	}

	return nil
}

// AddFlag flags the entity
func (entity *Entity) AddFlag(flag EntityFlags) {
	entity.Flags = entity.Flags | flag
}

// IdProperty updates finds a property that's defined as an ID and if none is, tries to set one based on its name and type
func (entity *Entity) IdProperty() (*Property, error) {
	prop := entity.getIdProperty()
	if prop == nil {
		return nil, errors.New("ID property not defined")
	}
	return prop, nil
}

// FindPropertyByUid finds a property by Uid
func (entity *Entity) FindPropertyByUid(uid Uid) (*Property, error) {
	for _, property := range entity.Properties {
		propertyUid, _ := property.Id.GetUid()
		if propertyUid == uid {
			return property, nil
		}
	}

	return nil, fmt.Errorf("property with Uid %d not found in '%s'", uid, entity.Name)
}

// FindPropertyByName finds a property by name
func (entity *Entity) FindPropertyByName(name string) (*Property, error) {
	for _, property := range entity.Properties {
		if strings.ToLower(property.Name) == strings.ToLower(name) {
			return property, nil
		}
	}

	return nil, fmt.Errorf("property named '%s' not found in '%s'", name, entity.Name)
}

// CreateProperty creates a property
func (entity *Entity) CreateProperty() (*Property, error) {
	var id Id = 1
	if len(entity.Properties) > 0 {
		id = entity.LastPropertyId.getIdSafe() + 1
	}

	uniqueUid, err := entity.Model.GenerateUid()

	if err != nil {
		return nil, err
	}

	var property = CreateProperty(entity, id, uniqueUid)

	entity.Properties = append(entity.Properties, property)
	entity.LastPropertyId = property.Id

	return property, nil
}

// RemoveProperty removes a property
func (entity *Entity) RemoveProperty(property *Property) error {
	var indexToRemove = -1
	for index, prop := range entity.Properties {
		if prop == property {
			indexToRemove = index
			break
		}
	}

	if indexToRemove < 0 {
		return fmt.Errorf("can't remove property %s %s - not found", property.Name, property.Id)
	}

	// remove index from the property
	if property.IndexId != nil {
		if err := property.RemoveIndex(); err != nil {
			return err
		}
	}

	// remove from list
	entity.Properties = append(entity.Properties[:indexToRemove], entity.Properties[indexToRemove+1:]...)

	// store the UID in the "retired" list so that it's not reused in the future
	entity.Model.RetiredPropertyUids = append(entity.Model.RetiredPropertyUids, property.Id.getUidSafe())

	return nil
}

// FindRelationByUid Finds relation by Uid
func (entity *Entity) FindRelationByUid(uid Uid) (*StandaloneRelation, error) {
	for _, relation := range entity.Relations {
		relationUid, _ := relation.Id.GetUid()
		if relationUid == uid {
			return relation, nil
		}
	}

	return nil, fmt.Errorf("relation with Uid %d not found in '%s'", uid, entity.Name)
}

// FindRelationByName finds relation by name
func (entity *Entity) FindRelationByName(name string) (*StandaloneRelation, error) {
	for _, relation := range entity.Relations {
		if strings.ToLower(relation.Name) == strings.ToLower(name) {
			return relation, nil
		}
	}

	return nil, fmt.Errorf("relation named '%s' not found in '%s'", name, entity.Name)
}

// CreateRelation creates relation
func (entity *Entity) CreateRelation() (*StandaloneRelation, error) {
	id, err := entity.Model.createRelationId()
	if err != nil {
		return nil, err
	}

	var relation = CreateStandaloneRelation(entity, id)
	entity.Relations = append(entity.Relations, relation)
	return relation, nil
}

// RemoveRelation removes relation
func (entity *Entity) RemoveRelation(relation *StandaloneRelation) error {
	var indexToRemove = -1
	for index, rel := range entity.Relations {
		if rel == relation {
			indexToRemove = index
			break
		}
	}

	if indexToRemove < 0 {
		return fmt.Errorf("can't remove relation %s %s - not found", relation.Name, relation.Id)
	}

	// remove from list
	entity.Relations = append(entity.Relations[:indexToRemove], entity.Relations[indexToRemove+1:]...)

	// store the UID in the "retired" list so that it's not reused in the future
	entity.Model.RetiredRelationUids = append(entity.Model.RetiredRelationUids, relation.Id.getUidSafe())

	return nil
}

// containsUid recursively checks whether given Uid is present in the model
func (entity *Entity) containsUid(searched Uid) bool {
	if entity.Id.getUidSafe() == searched {
		return true
	}

	if entity.LastPropertyId.getUidSafe() == searched {
		return true
	}

	for _, property := range entity.Properties {
		if property.containsUid(searched) {
			return true
		}
	}

	for _, relation := range entity.Relations {
		if relation.Id.getUidSafe() == searched {
			return true
		}
	}

	return false
}
