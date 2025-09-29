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
	"fmt"
)

type HnswDistanceType string

const (
	HnswDistanceType_Unknown                 = "Unknown"
	HnswDistanceType_Euclidean               = "Euclidean"
	HnswDistanceType_Cosine                  = "Cosine"
	HnswDistanceType_DotProduct              = "DotProduct"
	HnswDistanceType_DotProductNonNormalized = "DotProductNonNormalized"
	HnswDistanceType_Geo                     = "Geo"
)

type HnswParams struct {
	Dimensions                    *uint64    `json:"dimensions,omitempty"`
	DistanceType                  string     `json:"distance-type,omitempty"`
	NeighborsPerNode              *uint32    `json:"neighbors-per-node,omitempty"`
	IndexingSearchCount           *uint32    `json:"indexing-search-count,omitempty"`
	ReparationBacklinkProbability *float32   `json:"reparation-backlink-probability,omitempty"`
	VectorCacheHintSizeKb         *uint64    `json:"vector-cache-hint-size-kb,omitempty"`
	Flags                         *HnswFlags `json:"flags,omitempty"`
}

// Property in a model
type Property struct {
	Id             IdUid         `json:"id"`
	Name           string        `json:"name"`
	IndexId        *IdUid        `json:"indexId,omitempty"` // a pointer because it may be nil
	Type           PropertyType  `json:"type"`
	ExternalName   string        `json:"externalName,omitempty"`
	ExternalType   ExternalType  `json:"externalType,omitempty"`
	Flags          PropertyFlags `json:"flags,omitempty"`
	RelationTarget string        `json:"relationTarget,omitempty"`
	Entity         *Entity       `json:"-"`
	UidRequest     bool          `json:"-"` // used when the user gives an empty uid annotation
	HnswParams     *HnswParams   `json:"hnswParams,omitempty"`
	Meta           PropertyMeta  `json:"-"`
	Comments       []string      `json:"-"`
}

// CreateProperty creates a property
func CreateProperty(entity *Entity, id Id, uid Uid) *Property {
	return &Property{
		Entity: entity,
		Id:     CreateIdUid(id, uid),
	}
}

func (property *Property) CreateHnswParams() error {
	if property.HnswParams != nil {
		return fmt.Errorf("Double set")
	}
	property.HnswParams = &HnswParams{}
	return nil
}

func (property *Property) CheckHnswParams() error {
	if property.HnswParams == nil {
		return fmt.Errorf("Need to set index='hnsw'")
	}
	return nil
}

// Validate performs initial validation of loaded data so that it doesn't have to be checked in each function
func (property *Property) Validate() error {
	if property.Entity == nil {
		return fmt.Errorf("undefined parent entity")
	}

	if err := property.Id.Validate(); err != nil {
		return err
	}

	if property.IndexId != nil {
		if err := property.IndexId.Validate(); err != nil {
			return fmt.Errorf("indexId: %s", err)
		}
	}

	if len(property.Name) == 0 {
		return fmt.Errorf("name is undefined")
	}

	// NOTE type can't be validated because entities are update one-by-one and so
	// on the second one, validate() during load would failonly check this
	// if property.Type == 0 {
	//	return fmt.Errorf("type is undefined")
	// }

	// IDs must not be tagged unsigned for compatibility reasons
	if property.IsIdProperty() {
		if !property.hasValidTypeAsId(nil) {
			return fmt.Errorf("invalid type on property marked as ID: %d", property.Type)
		}
	}

	return nil
}

func (property *Property) finalize() error {
	if property.Type == PropertyTypeRelation || property.IsIdProperty() {
		// IDs must not be tagged unsigned for model compatibility (across language bindings)
		property.Flags = property.Flags & ^PropertyFlagUnsigned
	}

	if property.IsIdProperty() {
		// always stored IDs as Long, regardless of their type in the code/declaration
		property.Type = PropertyTypeLong
	}

	return property.Validate()
}

func (property *Property) IsIdProperty() bool {
	return property.Flags&PropertyFlagId != 0
}

func (property *Property) hasValidTypeAsId(acceptedTypes []PropertyType) bool {
	if acceptedTypes == nil {
		return property.Type == PropertyTypeLong
	} else {
		for _, t := range acceptedTypes {
			if property.Type == t {
				return true
			}
		}
		return false
	}
}

// CreateIndex creates an index
func (property *Property) CreateIndex() error {
	if property.IndexId != nil {
		return fmt.Errorf("can't create an index - it already exists")
	}

	indexId, err := property.Entity.Model.createIndexId()
	if err != nil {
		return err
	}

	property.IndexId = &indexId
	return nil
}

// RemoveIndex removes an index
func (property *Property) RemoveIndex() error {
	if property.IndexId == nil {
		return fmt.Errorf("can't remove index - it's not defined")
	}

	property.Entity.Model.RetiredIndexUids = append(property.Entity.Model.RetiredIndexUids, property.IndexId.getUidSafe())

	property.IndexId = nil

	return nil
}

// AddFlag flags the property
func (property *Property) AddFlag(flag PropertyFlags) {
	property.Flags = property.Flags | flag
}

// SetIndex defines an index on the property
func (property *Property) SetIndex() error {
	if property.IndexId != nil {
		return fmt.Errorf("index is already defined")
	}
	var blank IdUid
	property.IndexId = &blank
	return nil
}

// containsUid recursively checks whether given Uid is present in the model
func (property *Property) containsUid(searched Uid) bool {
	if property.Id.getUidSafe() == searched {
		return true
	}

	if property.IndexId != nil && property.IndexId.getUidSafe() == searched {
		return true
	}

	return false
}

// FbvTableOffset calculates flatbuffers vTableOffset.
func (property *Property) FbvTableOffset() (uint16, error) {
	// derived from the FB generated code & https://google.github.io/flatbuffers/md__internals.html
	var result = 4 + 2*uint32(property.FbSlot())

	if uint32(uint16(result)) != result {
		return 0, fmt.Errorf("can't calculate FlatBuffers VTableOffset: property %s ID %s is too large",
			property.Name, property.Id)
	}

	return uint16(result), nil
}

// FbSlot is called from the template. It calculates flatbuffers slot number.
func (property *Property) FbSlot() int {
	return int(property.Id.getIdSafe() - 1)
}
