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

package flatbuffersc

// Note the cgo library path: loads compiled libraries from "flatbuffers-c-bridge/cmake-build". This is where the FB
// build script ./third_party/flatbuffers-c/build.sh outputs the static libraries unless other path is specified.

/*
#cgo LDFLAGS: -lstdc++ -lflatbuffers-c-bridge -lflatbuffers-c-bridge-flatc -lflatbuffers -lm
#cgo LDFLAGS: -L${SRCDIR}/../../../third_party/flatbuffers-c-bridge/cmake-build/
#include <stdlib.h>
#include "flatbuffersc.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/objectbox/objectbox-generator/internal/generator/flatbuffersc/reflection"
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
