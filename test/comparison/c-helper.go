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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/objectbox/objectbox-generator/v4/internal/generator"
	cgenerator "github.com/objectbox/objectbox-generator/v4/internal/generator/c"
	"github.com/objectbox/objectbox-generator/v4/test/assert"
	"github.com/objectbox/objectbox-generator/v4/test/build"
	"github.com/objectbox/objectbox-generator/v4/test/cmake"
)

type cTestHelper struct {
	cpp        bool
	canCompile bool
}

func (h *cTestHelper) init(t *testing.T, conf testSpec) {
	if !testing.Short() {
		var mandatory = h.cpp // we require cpp compilation to be available at the moment
		if _, isCI := os.LookupEnv("CI"); isCI {
			t.Log("CI environment variable defined. Compilation support is mandatory.")
			mandatory = true
		}

		h.canCompile = build.CanCompileObjectBoxCCpp(t, repoRoot(t), h.cpp, mandatory)
	}
}

func (h cTestHelper) generatorFor(t *testing.T, conf testSpec, sourceFile string, genDir string) generator.CodeGenerator {
	// make a copy of the default generator
	var gen = *conf.generator.(*cgenerator.CGenerator)
	return &gen
}

func (cTestHelper) prepareTempDir(t *testing.T, conf testSpec, srcDir, tempDir, tempRoot string) func(err error) error {
	return nil
}

func (h cTestHelper) build(t *testing.T, conf testSpec, dir string, expectedError error, errorTransformer func(err error) error) {
	if !h.canCompile {
		t.Skip("Compilation not available")
	}

	includeDir, err := filepath.Abs(dir) // main.c/cpp will include generated headers from here
	assert.NoErr(t, err)

	cmak := cmake.Cmake{
		Name:        "compilation-test",
		IsCpp:       h.cpp,
		IncludeDirs: append(build.IncludeDirs(repoRoot(t)), includeDir),
		LinkDirs:    build.LibDirs(repoRoot(t)),
		LinkLibs:    []string{"objectbox"},
	}
	assert.NoErr(t, cmak.CreateTempDirs())
	defer cmak.RemoveTempDirs()

	var mainFile string
	if cmak.IsCpp {
		if conf.targetLang == "cpp11" {
			cmak.Standard = 11
		} else {
			cmak.Standard = 14
		}
		mainFile = path.Join(cmak.ConfDir, "main.cpp")
	} else {
		cmak.Standard = 99
		mainFile = path.Join(cmak.ConfDir, "main.c")
		cmak.LinkLibs = append(cmak.LinkLibs, "flatccrt")
	}

	cmak.Files = append(cmak.Files, mainFile)

	assert.NoErr(t, cmak.WriteCMakeListsTxt())
	if testing.Verbose() {
		cml, err := cmak.GetCMakeListsTxt()
		assert.NoErr(t, err)
		t.Logf("Using CMakeLists.txt: %s", cml)
	}

	{ // write main.c/cpp to the conf dir - a simple one, just include all sources
		var mainSrc = ""
		if cmak.IsCpp {
			mainSrc = mainSrc + "#define OBX_CPP_FILE\n"
			mainSrc = mainSrc + "#include \"objectbox.hpp\"\n"
		} else {
			mainSrc = mainSrc + "#include \"objectbox.h\"\n"
		}

		files, err := ioutil.ReadDir(includeDir)
		assert.NoErr(t, err)
		for _, file := range files {
			if conf.generator.IsGeneratedFile(file.Name()) {
				mainSrc = mainSrc + "#include \"" + file.Name() + "\"\n"
			}
		}

		mainSrc = mainSrc + "int main(){ return 0; }\n\n"
		t.Logf("main.c/cpp file contents \n%s", mainSrc)
		assert.NoErr(t, ioutil.WriteFile(mainFile, []byte(mainSrc), 0600))
	}

	if stdOut, stdErr, err := cmak.Configure(); err != nil {
		assert.Failf(t, "cmake configuration failed: \n%s\n%s\n%s", stdOut, stdErr, err)
	} else {
		t.Logf("configuration output:\n%s", string(stdOut))
	}

	if stdOut, stdErr, err := cmak.BuildTarget(); err != nil {
		checkBuildError(t, errorTransformer, stdOut, stdErr, err, expectedError)
		assert.Failf(t, "cmake build failed: \n%s\n%s\n%s", stdOut, stdErr, err)
	} else {
		t.Logf("build output:\n%s", string(stdOut))
	}
}
