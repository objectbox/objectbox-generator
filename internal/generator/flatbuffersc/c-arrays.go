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

/*
#include <stdlib.h>

char** newCharArray(int size) {
	return calloc(sizeof(char*), size);
}

void setArrayString(const char **array, size_t index, const char *value) {
    array[index] = value;
}

void freeCharArray(char **a, int size) {
    for (size_t i = 0; i < size; i++) {
    	free(a[i]);
    }
    free(a);
}
*/
import "C"

type stringArray struct {
	cArray **C.char
	size   int
}

func (array *stringArray) free() {
	if array.cArray != nil {
		C.freeCharArray(array.cArray, C.int(array.size))
		array.cArray = nil
	}
}

func goStringArrayToC(values []string) *stringArray {
	result := &stringArray{
		cArray: C.newCharArray(C.int(len(values))),
		size:   len(values),
	}
	for i, s := range values {
		C.setArrayString(result.cArray, C.size_t(i), C.CString(s))
	}
	return result
}
