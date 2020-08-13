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
	"testing"

	"github.com/objectbox/objectbox-generator/internal/generator"
	"github.com/objectbox/objectbox-generator/internal/generator/c"
	"github.com/objectbox/objectbox-generator/test/assert"
	"github.com/objectbox/objectbox-generator/test/build"
	"github.com/objectbox/objectbox-generator/test/cmake"
	"github.com/objectbox/objectbox-generator/test/comparison"
)

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

func generateCCpp(t *testing.T, srcFile string, cpp bool, outDir string) {
	t.Logf("generating code for %s into %s", srcFile, outDir)
	var options = generator.Options{
		ModelInfoFile: path.Join(outDir, "objectbox-model.json"),
		CodeGenerator: &cgenerator.CGenerator{
			PlainC: !cpp,
		},
		InPath:  srcFile,
		OutPath: outDir,
	}
	assert.NoErr(t, generator.Process(options))
}

type CCppTestConf struct {
	Cmake *cmake.Cmake
}

func sourceExt(cpp bool) string {
	if cpp {
		return "cpp"
	} else {
		return "c"
	}
}

// CommonExecute executes the integration with the simple/common setup
func (conf *CCppTestConf) CommonExecute(t *testing.T, cpp bool) {
	conf.CreateCMake(t, cpp, "main."+sourceExt(cpp))
	conf.Generate(t, "")
	conf.Build(t)
	conf.Run(t, nil)
}

func (conf *CCppTestConf) Cleanup() {
	if conf.Cmake != nil {
		conf.Cmake.RemoveTempDirs()
	}
}

// CreateCMake creates temporary directories and configures the CMake project
func (conf *CCppTestConf) CreateCMake(t *testing.T, cpp bool, mainFile string) {
	var testSrcDir string
	var err error
	if cpp {
		testSrcDir, err = filepath.Abs("cpp")
	} else {
		testSrcDir, err = filepath.Abs("c")
	}
	assert.NoErr(t, err)

	if conf.Cmake != nil {
		t.Logf("Reusing the previous CMake configuration - just changing binary to %s", mainFile)
		assert.Eq(t, cpp, conf.Cmake.IsCpp)
	} else {
		conf.Cmake = &cmake.Cmake{
			Name:        t.Name(),
			IsCpp:       true,
			Standard:    11,
			IncludeDirs: append(build.IncludeDirs(repoRoot(t)), testSrcDir, filepath.Join(repoRoot(t), "test", "integration")),
			LinkDirs:    build.LibDirs(repoRoot(t)),
			LinkLibs:    []string{"objectbox"},
		}
		assert.NoErr(t, conf.Cmake.CreateTempDirs())
	}
	conf.Cmake.Files = []string{path.Join(testSrcDir, mainFile)}

	if *inSource {
		conf.Cmake.ConfDir = testSrcDir
	}
	conf.Cmake.IncludeDirs = append(conf.Cmake.IncludeDirs, conf.Cmake.ConfDir) // because of the generated files

	if !cpp {
		conf.Cmake.LinkLibs = append(conf.Cmake.LinkLibs, "flatccrt")
	}

	// Link the test executable statically on Windows or it won't execute in the temp dir (missing DLL)
	if runtime.GOOS == "windows" {
		conf.Cmake.LinkLibs = append(conf.Cmake.LinkLibs, "-static-libgcc")
		if cpp {
			conf.Cmake.LinkLibs = append(conf.Cmake.LinkLibs, "-static-libstdc++")
		}
	}
}

// Generate loads *.fbs files in the current dir (or the given schema file) and generates the code
func (conf *CCppTestConf) Generate(t *testing.T, schema string) {
	var inputFiles []string

	if len(schema) != 0 {
		tmpFile := filepath.Join(conf.Cmake.ConfDir, "schema.fbs")
		defer os.Remove(tmpFile)
		inputFiles = []string{tmpFile}

		assert.NoErr(t, ioutil.WriteFile(tmpFile, []byte(schema), 0600))
	} else {
		var err error
		inputFiles, err = filepath.Glob("*.fbs")
		assert.NoErr(t, err)
	}

	for _, file := range inputFiles {
		generateCCpp(t, file, conf.Cmake.IsCpp, conf.Cmake.ConfDir)
	}
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

		if stdOut, stdErr, err := conf.Cmake.Build(); err != nil {
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
			assert.NoErr(t, comparison.CopyFile(
				path.Join(repoRoot(t), build.ObjectBoxCDir, "lib", "objectbox.dll"),
				path.Join(conf.Cmake.BuildDir, "objectbox.dll"),
				0))
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
