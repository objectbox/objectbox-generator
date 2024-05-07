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

package cmake_test

import (
	"runtime"
	"testing"

	"github.com/objectbox/objectbox-generator/test/assert"
	"github.com/objectbox/objectbox-generator/test/cmake"
)

func TestLibExists(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	err := cmake.LibraryExists("nonsense", nil, nil, nil, nil)
	assert.Err(t, err)

	err = cmake.LibraryExists("", []string{"non-existent-lib/include.h"}, nil, nil, nil)
	assert.Err(t, err)

	if runtime.GOOS == "windows" {
		err = cmake.LibraryExists("", []string{"array"}, nil, nil, nil)
	} else {
		err = cmake.LibraryExists("stdc++", []string{"array"}, nil, nil, nil)
	}
	assert.NoErr(t, err)
}
