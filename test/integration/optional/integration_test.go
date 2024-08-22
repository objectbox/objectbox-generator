/*
 * ObjectBox Generator - a build time tool for ObjectBox
 * Copyright (C) 2020-2024 ObjectBox Ltd. All rights reserved.
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

package typeful

import (
	"testing"

	cgenerator "github.com/objectbox/objectbox-generator/v4/internal/generator/c"
	"github.com/objectbox/objectbox-generator/v4/test/integration"
)

const optionalSchemaFields = `
	id           : uint64	;
	/// objectbox:optional
	int          : int		;
	/// objectbox:optional
	int8         : int8		;
	/// objectbox:optional
	int16        : int16	;
	/// objectbox:optional
	int32        : int32	;
	/// objectbox:optional
	int64        : int64	;
	/// objectbox:optional
	uint         : uint		;
	/// objectbox:optional
	uint8        : uint8	;
	/// objectbox:optional
	uint16       : uint16	;
	/// objectbox:optional
	uint32       : uint32	;
	/// objectbox:optional
	uint64       : uint64	;
	/// objectbox:optional
	bool         : bool		;
	/// objectbox:optional
	string       : string	;
	/// objectbox:optional
	stringvector : [string]	;
	/// objectbox:optional
	byte         : byte		;
	/// objectbox:optional
	ubyte        : ubyte	;
	/// objectbox:optional
	bytevector   : [byte]	;
	/// objectbox:optional
	ubytevector  : [ubyte]	;
	/// objectbox:optional
	float32      : float32	;
	/// objectbox:optional
	float64      : float64	;
	/// objectbox:optional
	float        : float	;
	/// objectbox:optional
	floatvector  : [float]  ;
	/// objectbox:optional
	double       : double	;
	/// objectbox:relation=RelTarget, optional
	relId:ulong;
`

const asNullSchemaFields = `
	id           : uint64	;
	int          : int		;
	int8         : int8		;
	int16        : int16	;
	int32        : int32	;
	int64        : int64	;
	uint         : uint		;
	uint8        : uint8	;
	uint16       : uint16	;
	uint32       : uint32	;
	uint64       : uint64	;
	bool         : bool		;
	string       : string	;
	stringvector : [string]	;
	byte         : byte		;
	ubyte        : ubyte	;
	bytevector   : [byte]	;
	ubytevector  : [ubyte]	;
	float32      : float32	;
	float64      : float64	;
	float        : float	;
	floatvector  : [float]  ;
	double       : double	;
	/// objectbox:relation=RelTarget
	relId:ulong;
`

func TestCppAndC(t *testing.T) {
	conf := &integration.CCppTestConf{}
	defer conf.Cleanup()
	conf.CreateCMake(t, integration.Cpp17, "main.cpp")
	conf.Generate(t, map[string]string{"rel.fbs": "table RelTarget {id: uint64;}"})

	conf.Generator = &cgenerator.CGenerator{EmptyStringAsNull: true, NaNAsNull: true}
	conf.Generate(t, map[string]string{"as-null.fbs": "table AsNull {" + asNullSchemaFields + "}"})

	conf.Generator = &cgenerator.CGenerator{Optional: "std::optional"}
	conf.Generate(t, map[string]string{"std-optional.fbs": "table Optional {" + optionalSchemaFields + "}"})

	conf.Generator = &cgenerator.CGenerator{Optional: "std::optional", EmptyStringAsNull: true, NaNAsNull: true}
	conf.Generate(t, map[string]string{"std-optional-as-null.fbs": "table OptionalAsNull {" + optionalSchemaFields + "}"})

	conf.Generator = &cgenerator.CGenerator{Optional: "std::unique_ptr"}
	conf.Generate(t, map[string]string{"std-unique_ptr.fbs": "table UniquePtr {" + optionalSchemaFields + "}"})

	conf.Generator = &cgenerator.CGenerator{Optional: "std::unique_ptr", EmptyStringAsNull: true, NaNAsNull: true}
	conf.Generate(t, map[string]string{"std-unique_ptr-as-null.fbs": "table UniquePtrAsNull {" + optionalSchemaFields + "}"})

	conf.Generator = &cgenerator.CGenerator{Optional: "std::shared_ptr"}
	conf.Generate(t, map[string]string{"std-shared_ptr.fbs": "table SharedPtr {" + optionalSchemaFields + "}"})

	conf.Generator = &cgenerator.CGenerator{Optional: "std::shared_ptr", EmptyStringAsNull: true, NaNAsNull: true}
	conf.Generate(t, map[string]string{"std-shared_ptr-as-null.fbs": "table SharedPtrAsNull {" + optionalSchemaFields + "}"})

	conf.Generator = &cgenerator.CGenerator{PlainC: true, Optional: "ptr"}
	conf.Generate(t, map[string]string{"c-ptr.fbs": "table PlainCPtr {" + optionalSchemaFields + "}"})

	conf.Build(t)
	conf.Run(t, nil)
}
