/*
 * Copyright (C) 2020 ObjectBox Ltd. All rights reserved.
 * https://objectbox.io
 *
 * This file is part of ObjectBox Generator.
 *
 * ObjectBox Generator is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * ObjectBox Generator is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with ObjectBox Generator.  If not, see <http://www.gnu.org/licenses/>.
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
