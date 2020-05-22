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

package generator

import (
	"fmt"
	"log"

	"github.com/objectbox/objectbox-go/internal/generator/modelinfo"
)

func mergeBindingWithModelInfo(currentModel *modelinfo.ModelInfo, storedModel *modelinfo.ModelInfo) error {
	// we need to first prepare all entities - otherwise relations wouldn't be able to find them in the model
	var models = make([]*modelinfo.Entity, len(currentModel.Entities))
	var err error
	for k, bindingEntity := range currentModel.Entities {
		models[k], err = getModelEntity(bindingEntity, storedModel)
		if err != nil {
			return err
		}
	}

	for k, bindingEntity := range currentModel.Entities {
		if err := mergeModelEntity(bindingEntity, models[k], storedModel); err != nil {
			return err
		}
	}

	// NOTE this is not ideal as there could be models across multiple packages
	storedModel.Package = currentModel.Package

	return nil
}

func getModelEntity(currentEntity *modelinfo.Entity, storedModel *modelinfo.ModelInfo) (*modelinfo.Entity, error) {
	if uid, err := currentEntity.Id.GetUid(); err != nil {
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

func mergeModelEntity(currentEntity *modelinfo.Entity, storedEntity *modelinfo.Entity, storedModel *modelinfo.ModelInfo) (err error) {
	storedEntity.Name = currentEntity.Name

	// TODO not sure we need this check
	if _, _, err := storedEntity.Id.Get(); err != nil {
		return err
	} else {
		currentEntity.Id = storedEntity.Id
	}

	{ //region Properties

		// add all properties from the bindings to the model and update/rename the changed ones
		for _, bindingProperty := range currentEntity.Properties {
			if modelProperty, err := getModelProperty(bindingProperty, storedEntity, storedModel); err != nil {
				return err
			} else if err := mergeModelProperty(bindingProperty, modelProperty); err != nil {
				return err
			}
		}

		// remove the missing (removed) properties
		removedProperties := make([]*modelinfo.Property, 0)
		for _, modelProperty := range storedEntity.Properties {
			if !bindingPropertyExists(modelProperty, currentEntity) {
				removedProperties = append(removedProperties, modelProperty)
			}
		}

		for _, property := range removedProperties {
			if err := storedEntity.RemoveProperty(property); err != nil {
				return err
			}
		}

		currentEntity.LastPropertyId = storedEntity.LastPropertyId
	} //endregion

	{ //region Relations

		// add all standalone relations from the bindings to the model and update/rename the changed ones
		for _, bindingRelation := range currentEntity.Relations {
			if modelRelation, err := getModelRelation(bindingRelation, storedEntity); err != nil {
				return err
			} else if err := mergeModelRelation(bindingRelation, modelRelation, storedModel); err != nil {
				return err
			}
		}

		// remove the missing (removed) relations
		removedRelations := make([]*modelinfo.StandaloneRelation, 0)
		for _, modelRelation := range storedEntity.Relations {
			if !bindingRelationExists(modelRelation, currentEntity) {
				removedRelations = append(removedRelations, modelRelation)
			}
		}

		for _, relation := range removedRelations {
			if err := storedEntity.RemoveRelation(relation); err != nil {
				return err
			}
		}
	} //endregion

	return nil
}

func getModelProperty(currentProperty *modelinfo.Property, storedEntity *modelinfo.Entity, storedModel *modelinfo.ModelInfo) (*modelinfo.Property, error) {
	if uid, err := currentProperty.Id.GetUid(); err != nil {
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

		log.Printf("Notice - new UID was specified for the same property name '%s' - resetting value (recreating the property)", currentProperty.Path)
		return property, nil
	}

	// we don't care about this error, either the property is found or we create it
	property, _ := storedEntity.FindPropertyByName(currentProperty.Name)

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
			return nil, fmt.Errorf(`uid annotation value must not be empty on property %s, entity %s:
    [rename] apply the current UID %d
    [change/reset] apply a new UID %d`,
				currentProperty.Name, currentProperty.Entity.Name, uid, newUid)
		}
		return nil, fmt.Errorf("uid annotation value must not be empty on an unknown property %s, entity %s",
			currentProperty.Name, currentProperty.Entity.Name)
	}

	if property == nil {
		return storedEntity.CreateProperty()
	}

	return property, nil
}

func mergeModelProperty(currentProperty *modelinfo.Property, storedProperty *modelinfo.Property) error {
	storedProperty.Name = currentProperty.Name

	// handle "reset property data" use-case - adding a new UID to an existing property
	if uid, err := currentProperty.Id.GetUid(); err != nil {
		return err
	} else if uid != 0 {
		id, _, err := storedProperty.Id.Get()
		if err != nil {
			return err
		}
		storedProperty.Id = modelinfo.CreateIdUid(id, uid)
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
			var idUid = modelinfo.CreateIdUid(id, uid)
			currentProperty.IndexId = &idUid
		}
	}

	storedProperty.RelationTarget = currentProperty.RelationTarget
	storedProperty.Type = currentProperty.Type
	storedProperty.Flags = currentProperty.Flags

	return nil
}

func bindingPropertyExists(modelProperty *modelinfo.Property, bindingEntity *modelinfo.Entity) bool {
	for _, bindingProperty := range bindingEntity.Properties {
		if bindingProperty.Name == modelProperty.Name {
			return true
		}
	}

	return false
}

func getModelRelation(currentRelation *modelinfo.StandaloneRelation, storedEntity *modelinfo.Entity) (*modelinfo.StandaloneRelation, error) {
	if uid, err := currentRelation.Id.GetUid(); err != nil {
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
		return nil, fmt.Errorf("uid annotation value must not be empty (%s) on relation %s, entity %s",
			errInfo, currentRelation.Name, storedEntity.Name)
	}

	if relation == nil {
		return storedEntity.CreateRelation()
	}

	return relation, nil
}

func mergeModelRelation(currentRelation *modelinfo.StandaloneRelation, storedRelation *modelinfo.StandaloneRelation, storedModel *modelinfo.ModelInfo) (err error) {
	storedRelation.Name = currentRelation.Name

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

func bindingRelationExists(modelRelation *modelinfo.StandaloneRelation, bindingEntity *modelinfo.Entity) bool {
	for _, bindingRelation := range bindingEntity.Relations {
		if bindingRelation.Name == modelRelation.Name {
			return true
		}
	}

	return false
}
