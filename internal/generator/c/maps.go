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

package cgenerator

import (
	"github.com/google/flatbuffers/go"

	"github.com/objectbox/objectbox-go/internal/generator/fbsparser/reflection"
	"github.com/objectbox/objectbox-go/internal/generator/model"
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
	reflection.BaseTypeObj:    0, // TODO, also in other maps
	reflection.BaseTypeUnion:  0, // TODO, also in other maps
	reflection.BaseTypeArray:  0, // TODO, also in other maps
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

var fbsTypeToFlatccType = map[reflection.BaseType]string{
	reflection.BaseTypeNone:   "",
	reflection.BaseTypeUType:  "",
	reflection.BaseTypeBool:   "bool",
	reflection.BaseTypeByte:   "int8",
	reflection.BaseTypeUByte:  "uint8",
	reflection.BaseTypeShort:  "int16",
	reflection.BaseTypeUShort: "uint16",
	reflection.BaseTypeInt:    "int32",
	reflection.BaseTypeUInt:   "uint32",
	reflection.BaseTypeLong:   "int64",
	reflection.BaseTypeULong:  "uint64",
	reflection.BaseTypeFloat:  "float",
	reflection.BaseTypeDouble: "double",
	reflection.BaseTypeString: "", // TODO
	reflection.BaseTypeVector: "", // TODO
	reflection.BaseTypeObj:    "",
	reflection.BaseTypeUnion:  "",
	reflection.BaseTypeArray:  "",
}
