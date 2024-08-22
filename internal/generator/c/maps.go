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

package cgenerator

import (
	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/objectbox/objectbox-generator/v4/internal/generator/flatbuffersc/reflection"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/model"
)

var fbsTypeToObxType = map[reflection.BaseType]model.PropertyType{
	reflection.BaseTypeNone:   0,
	reflection.BaseTypeUType:  0,
	reflection.BaseTypeBool:   model.PropertyTypeBool,
	reflection.BaseTypeByte:   model.PropertyTypeByte,
	reflection.BaseTypeUByte:  model.PropertyTypeByte,
	reflection.BaseTypeShort:  model.PropertyTypeShort,
	reflection.BaseTypeUShort: model.PropertyTypeShort,
	reflection.BaseTypeInt:    model.PropertyTypeInt,
	reflection.BaseTypeUInt:   model.PropertyTypeInt,
	reflection.BaseTypeLong:   model.PropertyTypeLong,
	reflection.BaseTypeULong:  model.PropertyTypeLong,
	reflection.BaseTypeFloat:  model.PropertyTypeFloat,
	reflection.BaseTypeDouble: model.PropertyTypeDouble,
	reflection.BaseTypeString: model.PropertyTypeString,
	reflection.BaseTypeVector: 0, // handled in schema-reader
	reflection.BaseTypeObj:    0, // not supported
	reflection.BaseTypeUnion:  0, // not supported
	reflection.BaseTypeArray:  0, // not supported
}

var fbsTypeToObxFlag = map[reflection.BaseType]model.PropertyFlags{
	reflection.BaseTypeUByte:  model.PropertyFlagUnsigned,
	reflection.BaseTypeUShort: model.PropertyFlagUnsigned,
	reflection.BaseTypeUInt:   model.PropertyFlagUnsigned,
	reflection.BaseTypeULong:  model.PropertyFlagUnsigned,
}

var fbsTypeToCppType = map[reflection.BaseType]string{
	reflection.BaseTypeNone:   "",
	reflection.BaseTypeUType:  "",
	reflection.BaseTypeBool:   "bool",
	reflection.BaseTypeByte:   "int8_t",
	reflection.BaseTypeUByte:  "uint8_t",
	reflection.BaseTypeShort:  "int16_t",
	reflection.BaseTypeUShort: "uint16_t",
	reflection.BaseTypeInt:    "int32_t",
	reflection.BaseTypeUInt:   "uint32_t",
	reflection.BaseTypeLong:   "int64_t",
	reflection.BaseTypeULong:  "uint64_t",
	reflection.BaseTypeFloat:  "float",
	reflection.BaseTypeDouble: "double",
	reflection.BaseTypeString: "std::string",
	reflection.BaseTypeVector: "std::vector", // Note: additional handling in fbsField
	reflection.BaseTypeObj:    "",
	reflection.BaseTypeUnion:  "",
	reflection.BaseTypeArray:  "",
}

var fbsTypeSize = map[reflection.BaseType]uint8{
	reflection.BaseTypeNone:   0,
	reflection.BaseTypeUType:  0,
	reflection.BaseTypeBool:   flatbuffers.SizeBool,
	reflection.BaseTypeByte:   flatbuffers.SizeByte,
	reflection.BaseTypeUByte:  flatbuffers.SizeByte,
	reflection.BaseTypeShort:  flatbuffers.SizeInt16,
	reflection.BaseTypeUShort: flatbuffers.SizeUint16,
	reflection.BaseTypeInt:    flatbuffers.SizeInt32,
	reflection.BaseTypeUInt:   flatbuffers.SizeUint32,
	reflection.BaseTypeLong:   flatbuffers.SizeInt64,
	reflection.BaseTypeULong:  flatbuffers.SizeUint64,
	reflection.BaseTypeFloat:  flatbuffers.SizeFloat32,
	reflection.BaseTypeDouble: flatbuffers.SizeFloat64,
	reflection.BaseTypeString: flatbuffers.SizeUOffsetT,
	reflection.BaseTypeVector: flatbuffers.SizeUOffsetT,
	reflection.BaseTypeObj:    0,
	reflection.BaseTypeUnion:  0,
	reflection.BaseTypeArray:  0,
}

var fbsTypeToFlatccFnPrefix = map[reflection.BaseType]string{
	reflection.BaseTypeNone:   "",
	reflection.BaseTypeUType:  "",
	reflection.BaseTypeBool:   "flatbuffers_bool",
	reflection.BaseTypeByte:   "flatbuffers_int8",
	reflection.BaseTypeUByte:  "flatbuffers_uint8",
	reflection.BaseTypeShort:  "flatbuffers_int16",
	reflection.BaseTypeUShort: "flatbuffers_uint16",
	reflection.BaseTypeInt:    "flatbuffers_int32",
	reflection.BaseTypeUInt:   "flatbuffers_uint32",
	reflection.BaseTypeLong:   "flatbuffers_int64",
	reflection.BaseTypeULong:  "flatbuffers_uint64",
	reflection.BaseTypeFloat:  "flatbuffers_float",
	reflection.BaseTypeDouble: "flatbuffers_double",
	reflection.BaseTypeString: "__flatbuffers_soffset",
	reflection.BaseTypeVector: "__flatbuffers_soffset",
	reflection.BaseTypeObj:    "",
	reflection.BaseTypeUnion:  "",
	reflection.BaseTypeArray:  "",
}
