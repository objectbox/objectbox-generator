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

package generator

import (
	"errors"
	"fmt"
	"log"

	"github.com/objectbox/objectbox-generator/v4/internal/generator/model"
)

func mergeBindingWithModelInfo(currentModel *model.ModelInfo, storedModel *model.ModelInfo) error {
	// we need to first prepare all entities - otherwise relations wouldn't be able to find them in the model
	var models = make([]*model.Entity, len(currentModel.Entities))
	var err error
	for k, entity := range currentModel.Entities {
		models[k], err = getModelEntity(entity, storedModel)
		if err != nil {
			return fmt.Errorf("entity %s: %s", entity.Name, err)
		}
	}

	for k, entity := range currentModel.Entities {
		if err := mergeModelEntity(entity, models[k], storedModel); err != nil {
			return fmt.Errorf("merging entity %s: %s", entity.Name, err)
		}
	}

	currentModel.LastEntityId = storedModel.LastEntityId
	currentModel.LastIndexId = storedModel.LastIndexId
	currentModel.LastRelationId = storedModel.LastRelationId

	return nil
}

func getModelEntity(currentEntity *model.Entity, storedModel *model.ModelInfo) (*model.Entity, error) {
	if uid, err := currentEntity.Id.GetUidAllowZero(); err != nil {
		return nil, err
	} else if uid != 0 {
		return storedModel.FindEntityByUid(uid)
	}

	// we don't care about this error = either the entity is found or we create it
	entity, _ := storedModel.FindEntityByName(currentEntity.Name)

	// handle uid request
	if currentEntity.UidRequest {
		var errInfo string
		if entity != nil {
			uid, err := entity.Id.GetUid()
			if err != nil {
				return nil, err
			}
			errInfo = fmt.Sprintf("model entity UID = %d", uid)
		} else {
			errInfo = "entity not found in the model"
		}
		return nil, fmt.Errorf("uid annotation value must not be empty (%s) on entity %s", errInfo, currentEntity.Name)
	}

	if entity == nil {
		return storedModel.CreateEntity(currentEntity.Name)
	}

	return entity, nil
}

func mergeModelEntity(currentEntity *model.Entity, storedEntity *model.Entity, storedModel *model.ModelInfo) (err error) {
	storedEntity.Name = currentEntity.Name
	storedEntity.Flags = currentEntity.Flags
	storedEntity.Comments = currentEntity.Comments
	storedEntity.ExternalName = currentEntity.ExternalName

	if currentEntity.Meta != nil {
		storedEntity.Meta = currentEntity.Meta.Merge(storedEntity)
	} else {
		storedEntity.Meta = nil
	}

	// TODO not sure we need this check
	if _, _, err := storedEntity.Id.Get(); err != nil {
		return err
	} else {
		currentEntity.Id = storedEntity.Id
	}

	{ // region Properties

		// add all properties from the bindings to the model and update/rename the changed ones
		for _, currentProperty := range currentEntity.Properties {
			if modelProperty, err := getModelProperty(currentProperty, storedEntity, storedModel); err != nil {
				return fmt.Errorf("property %s: %s", currentProperty.Name, err)
			} else if err := mergeModelProperty(currentProperty, modelProperty); err != nil {
				return fmt.Errorf("merging property %s: %s", currentProperty.Name, err)
			}
		}

		// remove the missing (removed) properties
		removedProperties := make([]*model.Property, 0)
		for _, modelProperty := range storedEntity.Properties {
			if !bindingPropertyExists(modelProperty, currentEntity) {
				removedProperties = append(removedProperties, modelProperty)
			}
		}

		for _, property := range removedProperties {
			if err := storedEntity.RemoveProperty(property); err != nil {
				return fmt.Errorf("removing property %s: %s", property.Name, err)
			}
		}

		currentEntity.LastPropertyId = storedEntity.LastPropertyId
	} // endregion

	{ // region Relations

		// add all standalone relations from the bindings to the model and update/rename the changed ones
		for _, currentRelation := range currentEntity.Relations {
			if modelRelation, err := getModelRelation(currentRelation, storedEntity); err != nil {
				return fmt.Errorf("relation %s: %s", currentRelation.Name, err)
			} else if err := mergeModelRelation(currentRelation, modelRelation, storedModel); err != nil {
				return fmt.Errorf("merging relation %s: %s", currentRelation.Name, err)
			}
		}

		// remove the missing (removed) relations
		removedRelations := make([]*model.StandaloneRelation, 0)
		for _, modelRelation := range storedEntity.Relations {
			if !bindingRelationExists(modelRelation, currentEntity) {
				removedRelations = append(removedRelations, modelRelation)
			}
		}

		for _, relation := range removedRelations {
			if err := storedEntity.RemoveRelation(relation); err != nil {
				return fmt.Errorf("removing relation %s: %s", relation.Name, err)
			}
		}
	} // endregion

	return nil
}

func getModelProperty(currentProperty *model.Property, storedEntity *model.Entity, storedModel *model.ModelInfo) (*model.Property, error) {
	if uid, err := currentProperty.Id.GetUidAllowZero(); err != nil {
		return nil, err
	} else if uid != 0 {
		property, err := storedEntity.FindPropertyByUid(uid)
		if err == nil {
			return property, nil
		}

		// handle "reset property data" use-case - adding a new UID to an existing property
		property, err2 := storedEntity.FindPropertyByName(currentProperty.Name)
		if err2 != nil {
			return nil, fmt.Errorf("%v; %v", err, err2)
		}

		log.Printf("Notice - new UID was specified for the same property name '%s' - resetting value (recreating the property)", currentProperty.Name)
		return property, nil
	}

	// we don't care about this error, either the property is found or we create it
	property, _ := storedEntity.FindPropertyByName(currentProperty.Name)

	// This effectively check for duplicate property names in the entity source definition by checking that a model
	// property with this name has already been merged with another "currentProperty".
	if property != nil && property.Meta != nil {
		return nil, fmt.Errorf("duplicate property name (note that property names are case insensitive)")
	}

	// handle uid request
	if currentProperty.UidRequest {
		if property != nil {
			uid, err := property.Id.GetUid()
			if err != nil {
				return nil, err
			}
			newUid, err := storedModel.GenerateUid()
			if err != nil {
				return nil, err
			}

			// handle "reset property data" use-case - adding a new UID to an existing property
			return nil, fmt.Errorf(`uid annotation value must not be empty:
    [rename] apply the current UID %d
    [change/reset] apply a new UID %d`,
				uid, newUid)
		}
		return nil, errors.New("uid annotation value must not be empty, the property isn't present in the persisted model")
	}

	if property == nil {
		return storedEntity.CreateProperty()
	}

	return property, nil
}

func mergeModelProperty(currentProperty *model.Property, storedProperty *model.Property) error {
	storedProperty.Name = currentProperty.Name
	storedProperty.Comments = currentProperty.Comments

	if currentProperty.Meta != nil {
		storedProperty.Meta = currentProperty.Meta.Merge(storedProperty)
	} else {
		storedProperty.Meta = nil
	}

	// handle "reset property data" use-case - adding a new UID to an existing property
	if curUid, err := currentProperty.Id.GetUidAllowZero(); err != nil {
		return err
	} else if oldUid, err := storedProperty.Id.GetUidAllowZero(); err != nil {
		return err
	} else if curUid != 0 && oldUid != curUid {
		highestId, _, err := storedProperty.Entity.LastPropertyId.Get()
		if err != nil {
			return err
		}

		storedProperty.Id = model.CreateIdUid(highestId+1, curUid)
		storedProperty.Entity.LastPropertyId = storedProperty.Id
	}

	// TODO not sure we need this check
	if _, _, err := storedProperty.Id.Get(); err != nil {
		return err
	} else {
		currentProperty.Id = storedProperty.Id
	}

	if currentProperty.IndexId == nil {
		// if there shouldn't be an index
		if storedProperty.IndexId != nil {
			// if there originally was an index, remove it
			if err := storedProperty.RemoveIndex(); err != nil {
				return err
			}
		}
	} else {
		// if there should be an index, create it (or reuse an existing one)
		if storedProperty.IndexId == nil {
			if err := storedProperty.CreateIndex(); err != nil {
				return err
			}
		}

		if id, uid, err := storedProperty.IndexId.Get(); err != nil {
			return err
		} else {
			var idUid = model.CreateIdUid(id, uid)
			currentProperty.IndexId = &idUid
		}
	}

	storedProperty.RelationTarget = currentProperty.RelationTarget
	storedProperty.Type = currentProperty.Type
	storedProperty.Flags = currentProperty.Flags
	storedProperty.HnswParams = currentProperty.HnswParams
	storedProperty.ExternalName = currentProperty.ExternalName
	storedProperty.ExternalType = currentProperty.ExternalType

	return nil
}

func bindingPropertyExists(modelProperty *model.Property, bindingEntity *model.Entity) bool {
	for _, bindingProperty := range bindingEntity.Properties {
		if bindingProperty.Name == modelProperty.Name {
			return true
		}
	}

	return false
}

func getModelRelation(currentRelation *model.StandaloneRelation, storedEntity *model.Entity) (*model.StandaloneRelation, error) {
	if uid, err := currentRelation.Id.GetUidAllowZero(); err != nil {
		return nil, err
	} else if uid != 0 {
		return storedEntity.FindRelationByUid(uid)
	}

	// we don't care about this error, either the relation is found or we create it
	relation, _ := storedEntity.FindRelationByName(currentRelation.Name)

	// handle uid request
	if currentRelation.UidRequest {
		var errInfo string
		if relation != nil {
			uid, err := relation.Id.GetUid()
			if err != nil {
				return nil, err
			}
			errInfo = fmt.Sprintf("model relation UID = %d", uid)
		} else {
			errInfo = "relation not found in the model"
		}
		return nil, fmt.Errorf("uid annotation value must not be empty (%s)", errInfo)
	}

	if relation == nil {
		return storedEntity.CreateRelation()
	}

	return relation, nil
}

func mergeModelRelation(currentRelation *model.StandaloneRelation, storedRelation *model.StandaloneRelation, storedModel *model.ModelInfo) (err error) {
	storedRelation.Name = currentRelation.Name
	storedRelation.ExternalName = currentRelation.ExternalName
	storedRelation.ExternalType = currentRelation.ExternalType

	if currentRelation.Meta != nil {
		storedRelation.Meta = currentRelation.Meta.Merge(storedRelation)
	} else {
		storedRelation.Meta = nil
	}

	if _, _, err = storedRelation.Id.Get(); err != nil {
		return err
	} else {
		currentRelation.Id = storedRelation.Id
	}

	// find the target entity & read it's ID/UID for the binding code
	if targetEntity, err := storedModel.FindEntityByName(currentRelation.Target.Name); err != nil {
		return err
	} else if _, _, err = targetEntity.Id.Get(); err != nil {
		return err
	} else {
		currentRelation.Target.Id = targetEntity.Id
		storedRelation.SetTarget(targetEntity)
	}

	return nil
}

func bindingRelationExists(modelRelation *model.StandaloneRelation, bindingEntity *model.Entity) bool {
	for _, bindingRelation := range bindingEntity.Relations {
		if bindingRelation.Name == modelRelation.Name {
			return true
		}
	}

	return false
}
