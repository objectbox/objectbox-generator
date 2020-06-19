/*
 * Copyright (C) 2020 ObjectBox Ltd. All rights reserved.
 * https://objectbox.io
 *
 * This file is part of ObjectBox Generator.
 *
 * ObjectBox Generator is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * ObjectBox Generator is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with ObjectBox Generator.  If not, see <http://www.gnu.org/licenses/>.
 */

package compilers

import (
	"testing"

	"github.com/objectbox/objectbox-generator/test/assert"
	"github.com/objectbox/objectbox-generator/test/cmake"
)

// Check verifies the C/C++ objectbox test code can be compiled - whether the required libraries are available.
func CanCompileC(t *testing.T, cpp, required bool) bool {
	{ // check objectbox lib
		var includeFiles = []string{"objectbox.h"}
		if cpp {
			includeFiles = append(includeFiles, "objectbox-cpp.h")
		}
		assert.NoErr(t, cmake.LibraryExists("objectbox", includeFiles))
	}

	// check flatbuffers library availability
	var err error
	if cpp {
		// Note: we don't need flatbuffers library explicitly, it's part of objectbox at the moment.
		err = cmake.LibraryExists("", []string{"flatbuffers/flatbuffers.h"})
	} else {
		err = cmake.LibraryExists("flatccrt", []string{"flatcc/flatcc.h", "flatcc/flatcc_builder.h"})
	}

	if required {
		assert.NoErr(t, err)
	} else if err != nil {
		t.Logf("C/C++ compilation not available because %s", err)
		return false
	}
	return true
}
