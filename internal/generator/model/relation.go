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
)

// StandaloneRelation in a model
type StandaloneRelation struct {
	Id           IdUid                  `json:"id"`
	Name         string                 `json:"name"`
	Target       *Entity                `json:"-"` // TODO consider changing to TargetName, nothing else seems to be used.
	ExternalName string					`json:"externalName,omitempty"`
	ExternalType ExternalType			`json:"externalType,omitempty"`
	TargetId     IdUid                  `json:"targetId"`
	UidRequest   bool                   `json:"-"` // used when the user gives an empty uid annotation // TODO test
	Meta         StandaloneRelationMeta `json:"-"`
	entity       *Entity
}

// CreateStandaloneRelation creates a standalone relation
func CreateStandaloneRelation(entity *Entity, id IdUid) *StandaloneRelation {
	return &StandaloneRelation{entity: entity, Id: id}
}

// Validate performs initial validation of loaded data so that it doesn't have to be checked in each function
func (relation *StandaloneRelation) Validate() error {
	if err := relation.Id.Validate(); err != nil {
		return err
	}

	if len(relation.Name) == 0 {
		return errors.New("name is undefined")
	}

	if len(relation.TargetId) > 0 {
		if err := relation.TargetId.Validate(); err != nil {
			return err
		}

		for _, entity := range relation.entity.Model.Entities {
			if entity.Id == relation.TargetId {
				relation.Target = entity
			}
		}

		if relation.Target == nil {
			return fmt.Errorf("target entity ID %s not found", string(relation.TargetId))
		}
	}

	return nil
}

// SetTarget sets the relation target entity
func (relation *StandaloneRelation) SetTarget(entity *Entity) {
	relation.Target = entity
	relation.TargetId = entity.Id
}
