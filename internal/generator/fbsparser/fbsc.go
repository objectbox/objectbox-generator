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

/*
#cgo LDFLAGS: -lstdc++ -lflatbuffersc -lflatbuffers -lm
#include <stdlib.h>
#include "flatbuffersc.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"github.com/objectbox/objectbox-go/internal/generator/fbsparser/reflection"
	"reflect"
	"unsafe"
)

// TODO deduplicate, this is a copy from objectbox/c-arrays
func cVoidPtrToByteSlice(data unsafe.Pointer, size int, bytes *[]byte) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(bytes))
	header.Data = uintptr(data)
	header.Len = size
	header.Cap = size
}

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

	var bytes []byte
	cVoidPtrToByteSlice(fbsBytes.data, int(fbsBytes.size), &bytes)
	return reflection.GetRootAsSchema(bytes, 0), nil
}
