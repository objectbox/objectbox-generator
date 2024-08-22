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

package flatbuffersc

// Note the cgo library path: loads compiled libraries from "flatbuffers-c-bridge/cmake-build". This is where the FB
// build script ./third_party/flatbuffers-c/build.sh outputs the static libraries unless other path is specified.

/*
#cgo LDFLAGS: -lstdc++ -lflatbuffers-c-bridge -lflatbuffers-c-bridge-flatc -lflatbuffers -lm
#cgo LDFLAGS: -L${SRCDIR}/../../../third_party/flatbuffers-c-bridge/cmake-build/
#cgo windows LDFLAGS: -static -static-libgcc -static-libstdc++ -lpthread

#include <stdlib.h>
#include "flatbuffersc.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/objectbox/objectbox-generator/v4/internal/generator/flatbuffersc/reflection"
)

func ParseSchemaFile(filename string) (*reflection.Schema, error) {
	var cFilename = C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	var cErr *C.char = nil
	defer C.fbs_error_free(cErr)

	var fbsBytes *C.FBS_bytes = C.fbs_schema_parse_file(cFilename, &cErr)
	if fbsBytes == nil {
		if cErr == nil {
			return nil, errors.New("unknown error")
		}
		return nil, fmt.Errorf(C.GoString(cErr))
	}
	defer C.fbs_schema_free(fbsBytes)

	// make a copy of the bytes because the source is deallocated by fbs_schema_free
	var bytes []byte = C.GoBytes(fbsBytes.data, C.int(fbsBytes.size))

	return reflection.GetRootAsSchema(bytes, 0), nil
}

// ExecuteFlatc runs flatc with the given arguments and returns its exit code and error, if any
func ExecuteFlatc(args []string) (int, error) {
	var cErr *C.char = nil
	defer C.fbs_error_free(cErr)

	cArgs := goStringArrayToC(args)
	defer cArgs.free()

	var code = int(C.fbs_flatc(cArgs.cArray, C.size_t(cArgs.size), &cErr))
	if cErr != nil {
		return code, errors.New("flatc execution failed: " + C.GoString(cErr))
	}
	if code != 0 {
		return code, errors.New("flatc execution failed with an unknown error")
	}
	return code, nil
}
