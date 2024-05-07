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

// EntityMeta provides a way for bindings to provide additional information to other users of Entity
type EntityMeta interface {
	// Merge produces new EntityMeta based on its internal state and given entity
	Merge(entity *Entity) EntityMeta
}

// PropertyMeta provides a way for bindings to provide additional information to other users of Property
type PropertyMeta interface {
	// Merge produces new PropertyMeta based on its internal state and given property
	Merge(property *Property) PropertyMeta
}

// StandaloneRelationMeta provides a way for bindings to provide additional information to other users of StandaloneRelation
type StandaloneRelationMeta interface {
	// Merge produces new StandaloneRelationMeta based on its internal state and given standalone relation
	Merge(relation *StandaloneRelation) StandaloneRelationMeta
}
