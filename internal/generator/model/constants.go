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

// EntityFlags is a bit combination of 0..n entity flags corresponding with objectbox-c
type EntityFlags int32

const (
	EntityFlagSyncEnabled     EntityFlags = 2
	EntityFlagSharedGlobalIds EntityFlags = 4
)

// EntityFlagNames assigns a name to each PropertyFlag
var EntityFlagNames = map[EntityFlags]string{
	EntityFlagSyncEnabled:     "SyncEnabled",
	EntityFlagSharedGlobalIds: "SharedGlobalIds",
}

// PropertyFlags is a bit combination of 0..n property flags corresponding with objectbox-c
type PropertyFlags int32

const (
	PropertyFlagId                   PropertyFlags = 1
	PropertyFlagNonPrimitiveType     PropertyFlags = 2
	PropertyFlagNotNull              PropertyFlags = 4
	PropertyFlagIndexed              PropertyFlags = 8
	PropertyFlagReserved             PropertyFlags = 16
	PropertyFlagUnique               PropertyFlags = 32
	PropertyFlagIdMonotonicSequence  PropertyFlags = 64
	PropertyFlagIdSelfAssignable     PropertyFlags = 128
	PropertyFlagIndexPartialSkipNull PropertyFlags = 256
	PropertyFlagIndexPartialSkipZero PropertyFlags = 512
	PropertyFlagVirtual              PropertyFlags = 1024
	PropertyFlagIndexHash            PropertyFlags = 2048
	PropertyFlagIndexHash64          PropertyFlags = 4096
	PropertyFlagUnsigned             PropertyFlags = 8192
	PropertyFlagIdCompanion          PropertyFlags = 16384
)

// PropertyFlagNames assigns a name to each PropertyFlag
var PropertyFlagNames = map[PropertyFlags]string{
	PropertyFlagId:                   "Id",
	PropertyFlagNonPrimitiveType:     "NonPrimitiveType",
	PropertyFlagNotNull:              "NotNull",
	PropertyFlagIndexed:              "Indexed",
	PropertyFlagReserved:             "Reserved",
	PropertyFlagUnique:               "Unique",
	PropertyFlagIdMonotonicSequence:  "IdMonotonicSequence",
	PropertyFlagIdSelfAssignable:     "IdSelfAssignable",
	PropertyFlagIndexPartialSkipNull: "IndexPartialSkipNull",
	PropertyFlagIndexPartialSkipZero: "IndexPartialSkipZero",
	PropertyFlagVirtual:              "Virtual",
	PropertyFlagIndexHash:            "IndexHash",
	PropertyFlagIndexHash64:          "IndexHash64",
	PropertyFlagUnsigned:             "Unsigned",
	PropertyFlagIdCompanion:          "IdCompanion",
}

// PropertyType is an identifier of a property type corresponding with objectbox-c
type PropertyType int8

const (
	PropertyTypeBool         PropertyType = 1
	PropertyTypeByte         PropertyType = 2
	PropertyTypeShort        PropertyType = 3
	PropertyTypeChar         PropertyType = 4
	PropertyTypeInt          PropertyType = 5
	PropertyTypeLong         PropertyType = 6
	PropertyTypeFloat        PropertyType = 7
	PropertyTypeDouble       PropertyType = 8
	PropertyTypeString       PropertyType = 9
	PropertyTypeDate         PropertyType = 10
	PropertyTypeRelation     PropertyType = 11
	PropertyTypeDateNano     PropertyType = 12
	PropertyTypeByteVector   PropertyType = 23
	PropertyTypeFloatVector  PropertyType = 28
	PropertyTypeStringVector PropertyType = 30
)

// PropertyTypeNames assigns a name to each PropertyType
var PropertyTypeNames = map[PropertyType]string{
	PropertyTypeBool:         "Bool",
	PropertyTypeByte:         "Byte",
	PropertyTypeShort:        "Short",
	PropertyTypeChar:         "Char",
	PropertyTypeInt:          "Int",
	PropertyTypeLong:         "Long",
	PropertyTypeFloat:        "Float",
	PropertyTypeDouble:       "Double",
	PropertyTypeString:       "String",
	PropertyTypeDate:         "Date",
	PropertyTypeRelation:     "Relation",
	PropertyTypeDateNano:     "DateNano",
	PropertyTypeByteVector:   "ByteVector",
	PropertyTypeFloatVector:  "FloatVector",
	PropertyTypeStringVector: "StringVector",
}

type ExternalType int32

const (
	ExternalTypeNone           ExternalType = 0
	ExternalTypeInt128         ExternalType = 100
	ExternalTypeUuid           ExternalType = 102
	ExternalTypeDecimal128     ExternalType = 103
	ExternalTypeUuidString     ExternalType = 104
	ExternalTypeUuidV4         ExternalType = 105
	ExternalTypeUuidV4String   ExternalType = 106
	ExternalTypeFlexMap        ExternalType = 107
	ExternalTypeFlexVector     ExternalType = 108
	ExternalTypeJson           ExternalType = 109
	ExternalTypeBson           ExternalType = 110
	ExternalTypeJavaScript     ExternalType = 111
	ExternalTypeJsonToNative   ExternalType = 112
	ExternalTypeInt128Vector   ExternalType = 116
	ExternalTypeUuidVector     ExternalType = 118
	ExternalTypeMongoId        ExternalType = 123
	ExternalTypeMongoIdVector  ExternalType = 124
	ExternalTypeMongoTimestamp ExternalType = 125
	ExternalTypeMongoBinary    ExternalType = 126
	ExternalTypeMongoRegex     ExternalType = 127
)

var ExternalTypeNames = map[ExternalType]string{
	ExternalTypeNone:           "None",
	ExternalTypeInt128:         "Int128",
	ExternalTypeUuid:           "Uuid",
	ExternalTypeDecimal128:     "Decimal128",
	ExternalTypeUuidString:     "UuidString",
	ExternalTypeUuidV4:         "UuidV4",
	ExternalTypeUuidV4String:   "UuidV4String",
	ExternalTypeFlexMap:        "FlexMap",
	ExternalTypeFlexVector:     "FlexVector",
	ExternalTypeJson:           "Json",
	ExternalTypeBson:           "Bson",
	ExternalTypeJavaScript:     "JavaScript",
	ExternalTypeJsonToNative:   "JsonToNative",
	ExternalTypeInt128Vector:   "Int128Vector",
	ExternalTypeUuidVector:     "UuidVector",
	ExternalTypeMongoId:        "MongoId",
	ExternalTypeMongoIdVector:  "MongoIdVector",
	ExternalTypeMongoTimestamp: "MongoTimestamp",
	ExternalTypeMongoBinary:    "MongoBinary",
	ExternalTypeMongoRegex:     "MongoRegex",
}

var ExternalTypeValues = map[string]ExternalType{
	"None":           ExternalTypeNone,
	"Int128":         ExternalTypeInt128,
	"Uuid":           ExternalTypeUuid,
	"Decimal128":     ExternalTypeDecimal128,
	"UuidString":     ExternalTypeUuidString,
	"UuidV4":         ExternalTypeUuidV4,
	"UuidV4String":   ExternalTypeUuidV4String,
	"FlexMap":        ExternalTypeFlexMap,
	"FlexVector":     ExternalTypeFlexVector,
	"Json":           ExternalTypeJson,
	"Bson":           ExternalTypeBson,
	"JavaScript":     ExternalTypeJavaScript,
	"JsonToNative":   ExternalTypeJsonToNative,
	"Int128Vector":   ExternalTypeInt128Vector,
	"UuidVector":     ExternalTypeUuidVector,
	"MongoId":        ExternalTypeMongoId,
	"MongoIdVector":  ExternalTypeMongoIdVector,
	"MongoTimestamp": ExternalTypeMongoTimestamp,
	"MongoBinary":    ExternalTypeMongoBinary,
	"MongoRegex":     ExternalTypeMongoRegex,
}

// HnswFlags is a bit combination of 0..n Hnsw flags corresponding with objectbox-c
type HnswFlags int32

const (
	HnswFlagNone                      HnswFlags = 0
	HnswFlagDebugLogs                 HnswFlags = 1
	HnswFlagDebugLogsDetailed         HnswFlags = 2
	HnswFlagVectorCacheSimdPaddingOff HnswFlags = 4
	HnswFlagReparationLimitCandidates HnswFlags = 8
)

// HnswFlagNames assigns a name to each PropertyFlag
var HnswFlagNames = map[HnswFlags]string{
	HnswFlagNone:                      "None",
	HnswFlagDebugLogs:                 "DebugLogs",
	HnswFlagDebugLogsDetailed:         "DebugLogsDetailed",
	HnswFlagVectorCacheSimdPaddingOff: "VectorCacheSimdPaddingOff",
	HnswFlagReparationLimitCandidates: "ReparationLimitCandidates",
}

var HnswFlagValues = map[string]HnswFlags{
	"None":                      HnswFlagNone,
	"DebugLogs":                 HnswFlagDebugLogs,
	"DebugLogsDetailed":         HnswFlagDebugLogsDetailed,
	"VectorCacheSimdPaddingOff": HnswFlagVectorCacheSimdPaddingOff,
	"ReparationLimitCandidates": HnswFlagReparationLimitCandidates,
}
