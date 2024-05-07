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

import "fmt"

// CheckRelationCycles finds relations cycles
func (model *ModelInfo) CheckRelationCycles() error {
	// DFS cycle check, storing relation path in the recursion stack
	var recursionStack = make(map[*Entity]bool)
	for _, entity := range model.Entities {
		if err := entity.checkRelationCycles(&recursionStack, entity.Name); err != nil {
			return err
		}
	}

	return nil
}

func (entity *Entity) checkRelationCycles(recursionStack *map[*Entity]bool, path string) error {
	(*recursionStack)[entity] = true

	// to-many relations
	for _, rel := range entity.Relations {
		if err := checkRelationCycle(recursionStack, path+"."+rel.Name, rel.Target); err != nil {
			return err
		}
	}

	// to-one relations
	for _, prop := range entity.Properties {
		if prop.RelationTarget == "" {
			continue
		}

		relTarget, _ := entity.Model.FindEntityByName(prop.RelationTarget)

		if err := checkRelationCycle(recursionStack, path+"."+prop.Name, relTarget); err != nil {
			return err
		}
	}

	delete(*recursionStack, entity)
	return nil
}

func checkRelationCycle(recursionStack *map[*Entity]bool, path string, relTarget *Entity) error {
	// this happens if the entity containing this relation haven't been defined in this file
	if relTarget == nil {
		return nil
	}

	if (*recursionStack)[relTarget] {
		return fmt.Errorf("relation cycle detected: %s (%s)", path, relTarget.Name)
	}

	return relTarget.checkRelationCycles(recursionStack, path)
}
