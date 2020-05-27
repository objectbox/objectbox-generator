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
