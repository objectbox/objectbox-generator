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

package modelinfo

const (
	PropertyFlagId                   = 1
	PropertyFlagNonPrimitiveType     = 2
	PropertyFlagNotNull              = 4
	PropertyFlagIndexed              = 8
	PropertyFlagReserved             = 16
	PropertyFlagUnique               = 32
	PropertyFlaIdMonotonicSequence   = 64
	PropertyFlagIdSelfAssignable     = 128
	PropertyFlagIndexPartialSkipNull = 256
	PropertyFlagIndexPartialSkipZero = 512
	PropertyFlagVirtual              = 1024
	PropertyFlagIndexHash            = 2048
	PropertyFlagIndexHash64          = 4096
	PropertyFlagUnsigned             = 8192
)

const (
	PropertyTypeBool         = 1
	PropertyTypeByte         = 2
	PropertyTypeShort        = 3
	PropertyTypeChar         = 4
	PropertyTypeInt          = 5
	PropertyTypeLong         = 6
	PropertyTypeFloat        = 7
	PropertyTypeDouble       = 8
	PropertyTypeString       = 9
	PropertyTypeDate         = 10
	PropertyTypeRelation     = 11
	PropertyTypeByteVector   = 23
	PropertyTypeStringVector = 30
)
