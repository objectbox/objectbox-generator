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
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/objectbox/objectbox-generator/internal/generator"
	cgenerator "github.com/objectbox/objectbox-generator/internal/generator/c"
	"github.com/objectbox/objectbox-generator/test/assert"
)

var cmakeListsTpl = template.Must(template.New("cmake").Parse(`
cmake_minimum_required(VERSION 3.0)
{{if .Cpp}}
set(CMAKE_CXX_STANDARD 11)
{{else}}
set(CMAKE_C_STANDARD 99)
{{end}}
project(compilation-test {{if not .Cpp}}C{{end}})

add_executable(${PROJECT_NAME} {{.Main}})
target_include_directories(${PROJECT_NAME} PRIVATE {{.Include}})
target_link_libraries(${PROJECT_NAME} objectbox)
`))

type cTestHelper struct {
	cpp bool
}

func (h cTestHelper) generatorFor(t *testing.T, conf testSpec, sourceFile string, genDir string) generator.CodeGenerator {
	// make a copy of the default generator
	var gen = *conf.generator.(*cgenerator.CGenerator)
	gen.OutPath = genDir
	return &gen
}

func (cTestHelper) prepareTempDir(t *testing.T, conf testSpec, srcDir, tempDir, tempRoot string) func(err error) error {
	return nil
}

func (h cTestHelper) build(t *testing.T, conf testSpec, dir string, expectedError error, errorTransformer func(err error) error) {
	includeDir, err := filepath.Abs(dir) // main.c/cpp will include generated headers from here
	assert.NoErr(t, err)

	tempRoot, err := ioutil.TempDir("", "objectbox-generator-test-build")
	assert.NoErr(t, err)
	defer os.RemoveAll(tempRoot)

	buildDir := path.Join(tempRoot, "build")
	cmakeConfDir := path.Join(tempRoot, "conf") // using "conf" dir to write CMakeLists.txt and main.c/cpp
	assert.NoErr(t, os.Mkdir(buildDir, 0700))
	assert.NoErr(t, os.Mkdir(cmakeConfDir, 0700))

	mainFile := path.Join(cmakeConfDir, "main.c")
	if h.cpp {
		mainFile = path.Join(cmakeConfDir, "main.cpp")
	}

	{ // write CMakeLists.txt to the conf dir
		var tplArguments = struct {
			Cpp     bool
			Ext     string
			Include string
			Main    string
		}{h.cpp, conf.generatedExt, includeDir, mainFile}

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)
		assert.NoErr(t, cmakeListsTpl.Execute(writer, tplArguments))
		assert.NoErr(t, writer.Flush())
		assert.NoErr(t, ioutil.WriteFile(path.Join(cmakeConfDir, "CMakeLists.txt"), b.Bytes(), 0600))
	}

	{ // write main.c/cpp to the conf dir - a simple one, just include all sources
		var mainSrc = ""

		files, err := ioutil.ReadDir(includeDir)
		assert.NoErr(t, err)
		for _, file := range files {
			if conf.generator.IsGeneratedFile(file.Name()) {
				mainSrc = mainSrc + "#include \"" + file.Name() + "\"\n"
			}
		}
		mainSrc = mainSrc + "int main(){ return 0; }\n\n"
		assert.NoErr(t, ioutil.WriteFile(mainFile, []byte(mainSrc), 0600))
	}

	// configure the cmake project
	stdOut, stdErr, err := cmake(buildDir, cmakeConfDir)
	if err != nil {
		assert.Failf(t, "cmake build configuration failed: \n%s\n%s\n%s", stdOut, stdErr, err)
	}
	if testing.Verbose() {
		t.Logf("configuration output:\n%s", string(stdOut))
	}

	// build the code
	stdOut, stdErr, err = cmake(includeDir, "--build", buildDir)

	checkBuildError(t, errorTransformer, stdOut, stdErr, err, expectedError)

	if testing.Verbose() {
		t.Logf("build output:\n%s", string(stdOut))
	}
}

func cmake(cwd string, args ...string) (stdOut []byte, stdErr []byte, err error) {
	var cmd = exec.Command("cmake", args...)
	cmd.Dir = cwd
	stdOut, err = cmd.Output()
	if ee, ok := err.(*exec.ExitError); ok {
		stdErr = ee.Stderr
	}
	return
}
