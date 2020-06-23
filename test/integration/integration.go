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
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"

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

func generateCCpp(t *testing.T, srcFile string, cpp bool, outDir string) {
	var args = []string{"run", path.Join(repoRoot(t), "cmd/objectbox-generator")}
	if cpp {
		args = append(args, "-cpp")
	} else {
		args = append(args, "-c")
	}
	args = append(args, "-out", outDir)
	args = append(args, "-persist", path.Join(outDir, "objectbox-model.json"))
	args = append(args, srcFile)

	t.Logf("executing generator %v", args)

	var cmd = exec.Command("go", args...)
	stdOut, err := cmd.Output()
	if ee, ok := err.(*exec.ExitError); ok {
		t.Fatalf("code generation failed: \n%s\n%s", string(stdOut), string(ee.Stderr))
	}
	t.Logf("code generation successful: %s", string(stdOut))
	assert.NoErr(t, err)
	return
}

func TestCCpp(t *testing.T, cpp bool) {
	var testSrcDir string
	var err error
	if cpp {
		testSrcDir, err = filepath.Abs("cpp")
	} else {
		testSrcDir, err = filepath.Abs("c")
	}
	assert.NoErr(t, err)

	cmak := cmake.Cmake{
		Name:        t.Name(),
		IsCpp:       true,
		Standard:    11,
		Files:       []string{path.Join(testSrcDir, "main.cpp")},
		IncludeDirs: append(build.IncludeDirs(repoRoot(t)), testSrcDir),
		LinkDirs:    build.LibDirs(repoRoot(t)),
		LinkLibs:    []string{"objectbox"},
	}

	assert.NoErr(t, cmak.CreateTempDirs())
	defer cmak.RemoveTempDirs()

	if *inSource {
		cmak.ConfDir = testSrcDir
	}
	cmak.IncludeDirs = append(cmak.IncludeDirs, cmak.ConfDir) // because of the generated files

	if !cpp {
		cmak.LinkLibs = append(cmak.LinkLibs, "flatccrt")
	}

	// Generate all FBS files, putting output into the tempoorary directory
	{
		inputFiles, err := filepath.Glob("*.fbs")
		assert.NoErr(t, err)
		for _, file := range inputFiles {
			generateCCpp(t, file, cpp, cmak.ConfDir)
		}
	}

	// Build
	if !testing.Short() {
		assert.NoErr(t, cmak.WriteCMakeListsTxt())
		if testing.Verbose() {
			cml, err := cmak.GetCMakeListsTxt()
			assert.NoErr(t, err)
			t.Logf("Using CMakeLists.txt: %s", cml)
		}

		if stdOut, stdErr, err := cmak.Configure(); err != nil {
			t.Fatalf("cmake configuration failed: \n%s\n%s\n%s", stdOut, stdErr, err)
		} else {
			t.Logf("configuration output:\n%s", string(stdOut))
		}

		if stdOut, stdErr, err := cmak.Build(); err != nil {
			t.Fatalf("cmake build failed: \n%s\n%s\n%s", stdOut, stdErr, err)
		} else {
			t.Logf("build output:\n%s", string(stdOut))
		}
	}

	// Execute
	if !testing.Short() {
		var testExecutable = path.Join(cmak.BuildDir, cmak.Name)
		if runtime.GOOS == "windows" {
			testExecutable = testExecutable + ".exe"
			assert.NoErr(t, comparison.CopyFile(path.Join(repoRoot(t), build.ObjectBoxCDir, "lib", "objectbox.dll"), path.Join(cmak.BuildDir, "objectbox.dll"), 0))
		}
		var cmd = exec.Command(testExecutable)
		cmd.Dir = cmak.BuildDir
		stdOut, err := cmd.Output()
		if ee, ok := err.(*exec.ExitError); ok {
			t.Fatalf("compiled test failed: %s\n%s\n%s", err, string(stdOut), string(ee.Stderr))
		}
		t.Logf("compiled test output: \n%s", string(stdOut))
		assert.NoErr(t, err)
	}
}
