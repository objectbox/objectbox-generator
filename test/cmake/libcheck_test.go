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

package cmake_test

import (
	"testing"

	"github.com/objectbox/objectbox-generator/test/assert"
	"github.com/objectbox/objectbox-generator/test/cmake"
)

func TestLibExists(t *testing.T) {
	err := cmake.LibraryExists("nonsense", nil)
	assert.Err(t, err)

	err = cmake.LibraryExists("", []string{"non-existent-lib/include.h"})
	assert.Err(t, err)

	err = cmake.LibraryExists("stdc++", []string{"bits/stdc++.h"})
	assert.NoErr(t, err)
}
