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
	PropertyTypeByteVector   PropertyType = 23
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
	PropertyTypeByteVector:   "ByteVector",
	PropertyTypeStringVector: "StringVector",
}
