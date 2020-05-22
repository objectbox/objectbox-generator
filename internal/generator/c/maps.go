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
	"github.com/objectbox/objectbox-go/internal/generator/fbsparser/reflection"
	"github.com/objectbox/objectbox-go/internal/generator/modelinfo"
)

var fbsTypeToObxType = map[reflection.BaseType]int{
	reflection.BaseTypeNone:   0,
	reflection.BaseTypeUType:  0,
	reflection.BaseTypeBool:   modelinfo.PropertyTypeBool,
	reflection.BaseTypeByte:   modelinfo.PropertyTypeByte,
	reflection.BaseTypeUByte:  modelinfo.PropertyTypeByte,
	reflection.BaseTypeShort:  modelinfo.PropertyTypeShort,
	reflection.BaseTypeUShort: modelinfo.PropertyTypeShort,
	reflection.BaseTypeInt:    modelinfo.PropertyTypeInt,
	reflection.BaseTypeUInt:   modelinfo.PropertyTypeInt,
	reflection.BaseTypeLong:   modelinfo.PropertyTypeLong,
	reflection.BaseTypeULong:  modelinfo.PropertyTypeLong,
	reflection.BaseTypeFloat:  modelinfo.PropertyTypeFloat,
	reflection.BaseTypeDouble: modelinfo.PropertyTypeDouble,
	reflection.BaseTypeString: modelinfo.PropertyTypeString,
	reflection.BaseTypeVector: 0, // handled in schema-reader
	reflection.BaseTypeObj:    0, // TODO
	reflection.BaseTypeUnion:  0, // TODO
	reflection.BaseTypeArray:  0, // TODO
}

var fbsTypeToObxFlag = map[reflection.BaseType]int{
	reflection.BaseTypeUByte:  modelinfo.PropertyFlagUnsigned,
	reflection.BaseTypeUShort: modelinfo.PropertyFlagUnsigned,
	reflection.BaseTypeUInt:   modelinfo.PropertyFlagUnsigned,
	reflection.BaseTypeULong:  modelinfo.PropertyFlagUnsigned,
}
