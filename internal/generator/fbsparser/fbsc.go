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

package fbsparser

// Note the cgo library path: load libraries from "build-artifacts" in the root of this repository. This is where
// the FlatBuffers C-API build script ./build/build-flatbuffersc.sh outputs the static libraries.

/*
#cgo LDFLAGS: -lstdc++ -lflatbuffersc -lflatbuffers -lm
#cgo LDFLAGS: -L${SRCDIR}/../../../build-artifacts
#include <stdlib.h>
#include "flatbuffersc.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/objectbox/objectbox-go/internal/generator/fbsparser/reflection"
)

func ParseSchemaFile(filename string) (*reflection.Schema, error) {
	var cFilename = C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	var errText *C.char
	var fbsBytes *C.FBS_bytes = C.fbs_schema_parse_file(&errText, cFilename)
	if fbsBytes == nil {
		if errText == nil {
			return nil, errors.New("unknown error")
		}
		return nil, fmt.Errorf(C.GoString(errText))
	}
	defer C.fbs_schema_free(fbsBytes)

	// make a copy of the bytes because the source is deallocated by fbs_schema_free
	var bytes []byte = C.GoBytes(fbsBytes.data, C.int(fbsBytes.size))

	return reflection.GetRootAsSchema(bytes, 0), nil
}
