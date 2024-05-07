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
