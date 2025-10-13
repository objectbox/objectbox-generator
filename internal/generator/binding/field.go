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

// Field holds common field/property information used by specialized code parsers/generators.
// Additionally, it groups some shared logic, e.g. annotation processing
type Field struct {
	ModelProperty *model.Property
	Name          string
	Optional      string
	IsSkipped     bool
}

func CreateField(prop *model.Property) *Field {
	return &Field{ModelProperty: prop}
}

func (field *Field) SetName(name string) {
	field.Name = name
	if len(field.ModelProperty.Name) == 0 {
		field.ModelProperty.Name = name
	}
}
func (field *Field) PreProcessAnnotations(a map[string]*Annotation) error {
	field.IsSkipped = false
	for _, alternative := range []string{"-", "transient"} {
		if a[alternative] != nil {
			if len(a) != 1 || a[alternative].Value != "" {
				return errors.New("to ignore the property, use only `objectbox:\"" + alternative + "\"` as an annotation")
			}
			field.IsSkipped = true
			return nil
		}
	}
	return nil
}

// ProcessAnnotations checks all set annotations for any inconsistencies and sets local/property fields (flags, name, ...)
func (field *Field) ProcessAnnotations(a map[string]*Annotation) error {
	if err := field.PreProcessAnnotations(a); err != nil {
		return err
	}

	if field.IsSkipped {
		return nil
	}

	if a["id"] != nil {
		field.ModelProperty.AddFlag(model.PropertyFlagId)
		if hasDetail, err := HasBooleanDetail(a, "id", "assignable"); err != nil {
			return err
		} else if hasDetail {
			field.ModelProperty.AddFlag(model.PropertyFlagIdSelfAssignable)
		}
	}

	if a["name"] != nil {
		if len(a["name"].Value) == 0 {
			return fmt.Errorf("name annotation value must not be empty - it's the field name in DB")
		}
		field.ModelProperty.Name = a["name"].Value
	}

	if a["date"] != nil || a["date-nano"] != nil {
		if a["date"] != nil && a["date-nano"] != nil {
			return errors.New("date and date-nano annotations cannot be used at the same time")
		}

		if field.ModelProperty.Type != model.PropertyTypeLong {
			return fmt.Errorf("invalid underlying type '%v' for date/date-nano field; expecting long", model.PropertyTypeNames[field.ModelProperty.Type])
		}

		if a["date"] != nil {
			field.ModelProperty.Type = model.PropertyTypeDate
		} else {
			field.ModelProperty.Type = model.PropertyTypeDateNano
		}
	}

	if a["id-companion"] != nil {
		if field.ModelProperty.Type != model.PropertyTypeDate && field.ModelProperty.Type != model.PropertyTypeDateNano {
			return fmt.Errorf("invalid underlying type '%v' for ID companion field; expecting date/date-nano", model.PropertyTypeNames[field.ModelProperty.Type])
		}
		field.ModelProperty.AddFlag(model.PropertyFlagIdCompanion)
	}

	if a["unique"] != nil {
		field.ModelProperty.AddFlag(model.PropertyFlagUnique)

		// add a default index type, unless specified otherwise
		if a["index"] == nil {
			a["index"] = &Annotation{}
		}
	}

	if a["external-name"] != nil {
		field.ModelProperty.ExternalName = a["external-name"].Value
	}
	if a["external-type"] != nil {
		externalTypeValue := a["external-type"].Value
		externalType, exists := model.ExternalTypeValues[externalTypeValue]
		if !exists {
			return fmt.Errorf("invalid external-type '%s' for property %s - must be one of the supported external types", externalTypeValue, field.Name)
		}
		field.ModelProperty.ExternalType = externalType
	}

	if a["index"] != nil {
		switch strings.ToLower(a["index"].Value) {
		case "":
			// if the user doesn't define index type use the default based on the data-type
			if field.ModelProperty.Type == model.PropertyTypeString {
				field.ModelProperty.AddFlag(model.PropertyFlagIndexHash)
			} else {
				field.ModelProperty.AddFlag(model.PropertyFlagIndexed)
			}
		case "value":
			field.ModelProperty.AddFlag(model.PropertyFlagIndexed)
		case "hash":
			field.ModelProperty.AddFlag(model.PropertyFlagIndexHash)
		case "hash64":
			field.ModelProperty.AddFlag(model.PropertyFlagIndexHash64)
		case "hnsw":
			if field.ModelProperty.Type == model.PropertyTypeFloatVector {
				field.ModelProperty.CreateHnswParams()
				field.ModelProperty.AddFlag(model.PropertyFlagIndexed)
			} else {
				return fmt.Errorf("index type 'hnsw' only supported for float vectors")
			}
		default:
			return fmt.Errorf("unknown index type %s", a["index"].Value)
		}

		if err := field.ModelProperty.SetIndex(); err != nil {
			return err
		}
	}

	if a["uid"] != nil {
		if len(a["uid"].Value) == 0 {
			// in case the user doesn't provide `objectbox:"uid"` value, it's considered in-process of setting up UID
			// this flag is handled by the merge mechanism and prints the UID of the already existing property
			field.ModelProperty.UidRequest = true
		} else if uid, err := strconv.ParseUint(a["uid"].Value, 10, 64); err != nil {
			return fmt.Errorf("can't parse uid - %s", err)
		} else if id, err := field.ModelProperty.Id.GetIdAllowZero(); err != nil {
			return fmt.Errorf("can't parse property Id - %s", err)
		} else {
			field.ModelProperty.Id = model.CreateIdUid(id, uid)
		}
	}

	var toOneRelation = a["relation"]
	if toOneRelation == nil {
		toOneRelation = a["link"]
	}
	if toOneRelation != nil && field.ModelProperty.Type != model.PropertyTypeRelation {
		if field.ModelProperty.Type != model.PropertyTypeLong {
			return fmt.Errorf("invalid underlying type (PropertyType %v) for relation field; expecting long", model.PropertyTypeNames[field.ModelProperty.Type])
		}
		if len(toOneRelation.Value) == 0 {
			return errors.New("unknown link target entity, define by changing the `link` annotation to the `link=Entity` format")
		}
		field.ModelProperty.Type = model.PropertyTypeRelation
		field.ModelProperty.RelationTarget = toOneRelation.Value
		field.ModelProperty.AddFlag(model.PropertyFlagIndexed)
		field.ModelProperty.AddFlag(model.PropertyFlagIndexPartialSkipZero)

		if err := field.ModelProperty.SetIndex(); err != nil {
			return err
		}
	}

	if a["optional"] != nil {
		field.Optional = a["optional"].Value
	}

	if a["hnsw-dimensions"] != nil {
		if err := field.ModelProperty.CheckHnswParams(); err != nil {
			return err
		}
		dimensions, err := strconv.ParseUint(a["hnsw-dimensions"].Value, 10, 64)
		if err != nil {
			return fmt.Errorf("Annotation 'hnsw-dimensions' value type mismatch: %s", err)
		}
		field.ModelProperty.HnswParams.Dimensions = &dimensions
	}

	if a["hnsw-distance-type"] != nil {
		if err := field.ModelProperty.CheckHnswParams(); err != nil {
			return err
		}
		distanceType := a["hnsw-distance-type"].Value
		switch distanceType {
		case "Unknown", "Euclidean", "Cosine", "DotProduct", "DotProductNonNormalized", "Geo":
			break
		default:
			return fmt.Errorf("Annotation 'hnsw-distance-type' value type mismatch: must be one of 'Unknown', 'Euclidean', 'Cosine', 'DotProduct', 'DotProductNonNormalized', 'Geo'")
		}
		field.ModelProperty.HnswParams.DistanceType = distanceType
	}

	if a["hnsw-neighbors-per-node"] != nil {
		if err := field.ModelProperty.CheckHnswParams(); err != nil {
			return err
		}
		value, err := strconv.ParseUint(a["hnsw-neighbors-per-node"].Value, 10, 32)
		if err != nil {
			return fmt.Errorf("Annotation 'hnsw-neighbors-per-node' value type mismatch: %s", err)
		}
		var neighborsPerNode uint32 = uint32(value)
		field.ModelProperty.HnswParams.NeighborsPerNode = &neighborsPerNode
	}
	if a["hnsw-indexing-search-count"] != nil {
		if err := field.ModelProperty.CheckHnswParams(); err != nil {
			return err
		}
		value, err := strconv.ParseUint(a["hnsw-indexing-search-count"].Value, 10, 32)
		if err != nil {
			return fmt.Errorf("Annotation 'hnsw-neighbors-per-node' value type mismatch: %s", err)
		}
		var indexingSearchCount uint32 = uint32(value)
		field.ModelProperty.HnswParams.IndexingSearchCount = &indexingSearchCount
	}
	if a["hnsw-reparation-backlink-probability"] != nil {
		if err := field.ModelProperty.CheckHnswParams(); err != nil {
			return err
		}
		value, err := strconv.ParseFloat(a["hnsw-reparation-backlink-probability"].Value, 32)
		if err != nil {
			return fmt.Errorf("Annotation 'hnsw-reparation-backlink-probability' value type mismatch: %s", err)
		}
		var reparationBacklinkProbability float32 = float32(value)
		field.ModelProperty.HnswParams.ReparationBacklinkProbability = &reparationBacklinkProbability
	}
	if a["hnsw-vector-cache-hint-size-kb"] != nil {
		if err := field.ModelProperty.CheckHnswParams(); err != nil {
			return err
		}
		vectorCacheHintSizeKb, err := strconv.ParseUint(a["hnsw-vector-cache-hint-size-kb"].Value, 10, 64)
		if err != nil {
			return fmt.Errorf("Annotation 'hnsw-dimensions' value type mismatch: %s", err)
		}
		field.ModelProperty.HnswParams.VectorCacheHintSizeKb = &vectorCacheHintSizeKb

	}
	if a["hnsw-flags"] != nil {
		if err := field.ModelProperty.CheckHnswParams(); err != nil {
			return err
		}
		flagsString := a["hnsw-flags"].Value
		flagsVector := strings.Split(flagsString, "|")
		var flags model.HnswFlags = 0
		for _, v := range flagsVector {
			flagName := strings.TrimSpace(v)
			flagValue, ok := model.HnswFlagValues[flagName]
			if ok {
				flags |= flagValue
			} else {
				keys := make([]string, 0, len(model.HnswFlagValues))
				for k := range model.HnswFlagValues {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				availableFlags := strings.Join(keys, ", ")
				return fmt.Errorf("HNSW Flag unknown: '%s' (Available flags: %s)", flagName, availableFlags)
			}
		}
		field.ModelProperty.HnswParams.Flags = &flags
	}

	return nil
}
