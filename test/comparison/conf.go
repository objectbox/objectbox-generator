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
	// prepareTempDir prepares tempDir contents (already a copy of the srcDir) with any language specific setup.
	// Returns an errorTransformer.
	prepareTempDir(t *testing.T, conf testSpec, srcDir, tempDir, tempRoot string) func(err error) error

	// build compiles the code in the given directory
	build(t *testing.T, conf testSpec, dir string, expectedError error, errorTransformer func(err error) error)

	// args returns additional configuration. TBD whether this shouldn't be replaced with a generator.Options-updating function
	args(t *testing.T, sourceFile string) map[string]string
}

// Generator configurations for all supported languages. The map index is the top level directory.
type testSpec struct {
	sourceExt    string
	generatedExt string
	generator    generator.CodeGenerator
	helper       testHelper
}

var confs = map[string]testSpec{
	"c": {".fbs", ".obx.h", &cgenerator.CGenerator{PlainC: true}, cTestHelper{cpp: false}},
	// TODO "cpp": {".fbs", "-cpp.obx.h", &cgenerator.CGenerator{cpp: false}, cTestHelper{cpp: true}},
	// TODO "go":  {".go", ".obx.go", &gogenerator.GoGenerator{}, goTestHelper{}},
}
