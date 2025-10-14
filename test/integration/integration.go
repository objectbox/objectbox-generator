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

// package integration provides common tools for all integration test executors
package integration

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-generator/v4/internal/generator"
	cgenerator "github.com/objectbox/objectbox-generator/v4/internal/generator/c"
	"github.com/objectbox/objectbox-generator/v4/test/assert"
	"github.com/objectbox/objectbox-generator/v4/test/build"
	"github.com/objectbox/objectbox-generator/v4/test/cmake"
	"github.com/objectbox/objectbox-generator/v4/test/comparison"
)

type cCppStandard string

func (std cCppStandard) isCpp() bool {
	return strings.HasPrefix(string(std), "C++")
}

func (std cCppStandard) year() (int, error) {
	var prefixLen = 1
	if std.isCpp() {
		prefixLen = 3
	}
	num, err := strconv.ParseInt(string(std)[prefixLen:], 10, 64)
	return int(num), err
}

const Cpp11 = cCppStandard("C++11")
const Cpp14 = cCppStandard("C++14")
const Cpp17 = cCppStandard("C++17")

const CppDefault = Cpp14

// used during development to generate code into the source directory instead of temp
var inSource = flag.Bool("insource", false, "Output generated code to the source dir for development")

func repoRoot(t *testing.T) string {
	cwd, err := os.Getwd()
	assert.NoErr(t, err)
	return filepath.ToSlash(filepath.Join(cwd, "..", "..", ".."))
}

// Calls the generator command-line executable (using `go run`).
// Doesn't work well with go test cache - not invalidated when the generator code changes.
// Leaving the code here as a backup for now - if we wanted to do "real" integration test after all,
// e.g. in addition to generating with `generator.Process(srcFile, options)` as below.
//
// func generateCCpp(t *testing.T, srcFile string, cpp bool, outDir string) {
// 	var args = []string{"run", path.Join(repoRoot(t), "cmd/objectbox-generator")}
// 	if cpp {
// 		args = append(args, "-cpp")
// 	} else {
// 		args = append(args, "-c")
// 	}
// 	args = append(args, "-out", outDir)
// 	args = append(args, "-persist", path.Join(outDir, "objectbox-model.json"))
// 	args = append(args, srcFile)
//
// 	t.Logf("executing generator %v", args)
//
// 	var cmd = exec.Command("go", args...)
// 	stdOut, err := cmd.Output()
// 	if ee, ok := err.(*exec.ExitError); ok {
// 		t.Fatalf("code generation failed: \n%s\n%s", string(stdOut), string(ee.Stderr))
// 	}
// 	t.Logf("code generation successful: %s", string(stdOut))
// 	assert.NoErr(t, err)
// 	return
// }

func generateCCpp(t *testing.T, srcPath string, outDir string, cGenerator *cgenerator.CGenerator) {
	t.Logf("generating code for %s into %s", srcPath, outDir)
	var options = generator.Options{
		ModelInfoFile: path.Join(outDir, "objectbox-model.json"),
		CodeGenerator: cGenerator,
		InPath:        srcPath,
		OutPath:       outDir,
	}
	assert.NoErr(t, generator.Process(options))
}

type CCppTestConf struct {
	Cmake     *cmake.Cmake
	Generator *cgenerator.CGenerator
}

func sourceExt(cpp bool) string {
	if cpp {
		return "cpp"
	} else {
		return "c"
	}
}

// CommonExecute executes the integration with the simple/common setup
func (conf *CCppTestConf) CommonExecute(t *testing.T, lang cCppStandard) {
	conf.CreateCMake(t, lang, "main."+sourceExt(lang.isCpp()))
	conf.Generate(t, nil)
	conf.Build(t)
	conf.Run(t, nil)
}

func (conf *CCppTestConf) Cleanup() {
	if conf.Cmake != nil {
		conf.Cmake.RemoveTempDirs()
	}
}

// CreateCMake creates temporary directories and configures the CMake project
func (conf *CCppTestConf) CreateCMake(t *testing.T, lang cCppStandard, mainFile string) {
	var testSrcDir string
	var err error
	if lang.isCpp() {
		testSrcDir, err = filepath.Abs("cpp")
	} else {
		testSrcDir, err = filepath.Abs("c")
	}
	assert.NoErr(t, err)

	langYear, err := lang.year()
	assert.NoErr(t, err)

	if conf.Cmake != nil {
		t.Logf("Reusing the previous CMake configuration - just changing binary to %s", mainFile)
		assert.Eq(t, lang.isCpp(), conf.Cmake.IsCpp)
	} else {
		conf.Cmake = &cmake.Cmake{
			Name:        t.Name(),
			IsCpp:       true,
			Standard:    langYear,
			IncludeDirs: append(build.IncludeDirs(repoRoot(t)), testSrcDir, filepath.Join(repoRoot(t), "test", "integration")),
			LinkDirs:    build.LibDirs(repoRoot(t)),
			LinkLibs:    []string{"objectbox", "flatccrt"},
		}
		assert.NoErr(t, conf.Cmake.CreateTempDirs())
	}
	conf.Cmake.Files = []string{path.Join(testSrcDir, mainFile)}

	if *inSource {
		conf.Cmake.ConfDir = testSrcDir
	}
	conf.Cmake.IncludeDirs = append(conf.Cmake.IncludeDirs, conf.Cmake.ConfDir) // because of the generated files

	// Link the test executable statically on Windows or it won't execute in the temp dir (missing DLL)
	if runtime.GOOS == "windows" {
		conf.Cmake.LinkLibs = append(conf.Cmake.LinkLibs, "-static-libgcc")
		if lang.isCpp() {
			conf.Cmake.LinkLibs = append(conf.Cmake.LinkLibs, "-static-libstdc++")
		}
	}
}

// Generate loads *.fbs files in the current dir (or the given schema file) and generates the code
func (conf *CCppTestConf) Generate(t *testing.T, schemas map[string]string) {
	var srcPath string

	if len(schemas) != 0 {
		for name, content := range schemas {
			// passing an empty name and content is a trick to having multiple schames to enable wildcard generation.
			if len(name) == 0 && len(content) == 0 {
				continue
			}

			srcPath = filepath.Join(conf.Cmake.ConfDir, name)
			assert.NoErr(t, ioutil.WriteFile(srcPath, []byte(content), 0600))
		}
		if len(schemas) != 1 {
			srcPath = filepath.Join(conf.Cmake.ConfDir, "*.fbs")
		}
	} else {
		srcPath = "*.fbs"
	}

	var cGenerator = conf.Generator
	if cGenerator == nil {
		cGenerator = &cgenerator.CGenerator{
			PlainC: !conf.Cmake.IsCpp,
		}
	}

	generateCCpp(t, srcPath, conf.Cmake.ConfDir, cGenerator)
}

// Build compiles the test sources producing an executable
func (conf *CCppTestConf) Build(t *testing.T) {
	generatedSources, err := filepath.Glob(filepath.Join(conf.Cmake.ConfDir, "*obx."+sourceExt(conf.Cmake.IsCpp)))
	assert.NoErr(t, err)
	conf.Cmake.Files = append(conf.Cmake.Files, generatedSources...)

	assert.NoErr(t, conf.Cmake.WriteCMakeListsTxt())

	if !testing.Short() {
		if testing.Verbose() {
			cml, err := conf.Cmake.GetCMakeListsTxt()
			assert.NoErr(t, err)
			t.Logf("Using CMakeLists.txt: %s", cml)
		}

		if stdOut, stdErr, err := conf.Cmake.Configure(); err != nil {
			t.Fatalf("cmake configuration failed: \n%s\n%s\n%s", stdOut, stdErr, err)
		} else {
			t.Logf("configuration output:\n%s", string(stdOut))
		}

		if stdOut, stdErr, err := conf.Cmake.BuildTarget(); err != nil {
			t.Fatalf("cmake build failed: \n%s\n%s\n%s", stdOut, stdErr, err)
		} else {
			t.Logf("build output:\n%s", string(stdOut))
		}
	}
}

// Run executes the built test binary
func (conf *CCppTestConf) Run(t *testing.T, envVars []string) {
	if !testing.Short() {
		var testExecutable = path.Join(conf.Cmake.BuildDir, conf.Cmake.Name)
		if runtime.GOOS == "windows" {
			testExecutable = testExecutable + ".exe"
			conf.copyObjectBoxDll(t)
			conf.copyMinGwDlls(t)
		}
		var cmd = exec.Command(testExecutable)
		cmd.Dir = conf.Cmake.BuildDir
		cmd.Env = append(os.Environ(), envVars...)
		stdOut, err := cmd.Output()
		if ee, ok := err.(*exec.ExitError); ok {
			t.Fatalf("compiled test failed: %s\n%s\n%s", err, string(stdOut), string(ee.Stderr))
		}
		t.Logf("compiled test output: \n%s", string(stdOut))
		assert.NoErr(t, err)
	}
}

func (conf *CCppTestConf) copyObjectBoxDll(t *testing.T) {
	libDir := filepath.Join(repoRoot(t), build.ObjectBoxCDir, "lib")
	dllFiles, err := filepath.Glob(filepath.Join(libDir, "*.dll"))
	assert.NoErr(t, err)
	for _, dllFile := range dllFiles {
		dllName := filepath.Base(dllFile)
		t.Logf("Copying DLL from lib: %s", dllName)
		assert.NoErr(t, comparison.CopyFile(
			dllFile,
			filepath.Join(conf.Cmake.BuildDir, dllName),
			0))
	}
}

// Required on **some** CI runners only: copy MinGW runtime DLLs that objectbox.dll may depend on.
// Finds the C++ compiler path to locate the MinGW bin directory.
func (conf *CCppTestConf) copyMinGwDlls(t *testing.T) {
	cppPath, err := exec.LookPath("c++")
	if err == nil {
		mingwBinDir := filepath.Dir(cppPath)
		t.Logf("MinGW bin directory: %s", mingwBinDir)

		// Common MinGW runtime DLLs that objectbox.dll likely depends on
		runtimeDlls := []string{
			"libwinpthread-1.dll",
			"libgcc_s_seh-1.dll",
			"libgcc_s_dw2-1.dll", // Alternative GCC DLL
			"libstdc++-6.dll",
		}

		for _, dllName := range runtimeDlls {
			srcPath := filepath.Join(mingwBinDir, dllName)
			if _, err := os.Stat(srcPath); err == nil {
				t.Logf("Copying MinGW runtime DLL: %s", dllName)
				assert.NoErr(t, comparison.CopyFile(
					srcPath,
					filepath.Join(conf.Cmake.BuildDir, dllName),
					0))
			}
		}
	} else {
		t.Logf("Warning: Could not locate c++ compiler to find MinGW runtime DLLs: %v", err)
	}
}
