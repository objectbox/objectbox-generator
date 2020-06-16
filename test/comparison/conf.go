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

package comparison

import (
	"testing"

	"github.com/objectbox/objectbox-generator/internal/generator"
	"github.com/objectbox/objectbox-generator/internal/generator/c"
)

type testHelper interface {
	prepareTempDir(t *testing.T, srcDir, tempDir, tempRoot string) func(err error) error
	build(t *testing.T, dir string, errorTransformer func(err error) error)
}

// this containing module name - used for test case modules
const moduleName = "github.com/objectbox/objectbox-go"

// Generator configurations for all supported languages. The map index is the top level directory.
type testSpec struct {
	sourceExt    string
	generatedExt string
	generator    generator.CodeGenerator
	helper       testHelper
}

var confs = map[string]testSpec{
	"c": {".fbs", ".obx.h", &cgenerator.CGenerator{PlainC: true}, nil},
	// TODO "cpp": {".fbs", "-cpp.obx.h", &cgenerator.CGenerator{PlainC: false}, nil},
	// TODO "go":  {".go", ".obx.go", &gogenerator.GoGenerator{}, goTestHelper{}},
}
