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

package comparison

import (
	"testing"

	"github.com/objectbox/objectbox-generator/v4/internal/generator"
	cgenerator "github.com/objectbox/objectbox-generator/v4/internal/generator/c"
	gogenerator "github.com/objectbox/objectbox-generator/v4/internal/generator/go"
)

type testHelper interface {
	// init sets up the helper before the very first execution
	init(t *testing.T, conf testSpec)

	// generatorFor constructs and configures a code generator for the given source file
	generatorFor(t *testing.T, conf testSpec, sourceFile string, genDir string) generator.CodeGenerator

	// prepareTempDir prepares tempDir contents (already a copy of the srcDir) with any language specific setup.
	// Returns an errorTransformer.
	prepareTempDir(t *testing.T, conf testSpec, srcDir, tempDir, tempRoot string) func(err error) error

	// build compiles the code in the given directory
	build(t *testing.T, conf testSpec, dir string, expectedError error, errorTransformer func(err error) error)
}

// Generator configurations for all supported languages. The map index is the top level directory.
type testSpec struct {
	targetLang   string
	sourceExt    string
	generatedExt []string
	generator    generator.CodeGenerator
	helper       testHelper
}

var confs = map[string]testSpec{
	"fbs-c":     {"c", ".fbs", []string{".obx.h"}, &cgenerator.CGenerator{PlainC: true, LangVersion: -1}, &cTestHelper{cpp: false}},
	"fbs-cpp":   {"cpp", ".fbs", []string{".obx.hpp", ".obx.cpp"}, &cgenerator.CGenerator{PlainC: false, LangVersion: 14}, &cTestHelper{cpp: true}},
	"fbs-cpp11": {"cpp11", ".fbs", []string{".obx.hpp", ".obx.cpp"}, &cgenerator.CGenerator{PlainC: false, LangVersion: 11}, &cTestHelper{cpp: true}},
	"go":        {"go", ".go", []string{".obx.go"}, &gogenerator.GoGenerator{}, &goTestHelper{}},
}
