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
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-generator/internal/generator"
	gogenerator "github.com/objectbox/objectbox-generator/internal/generator/go"
	"github.com/objectbox/objectbox-generator/test/assert"
)

// this containing module name - used for test case modules
const goModuleName = "github.com/objectbox/objectbox-go"

var goGeneratorArgsRegexp = regexp.MustCompile("//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen (.+)[\n|\r]")

type goTestHelper struct{}

func (h *goTestHelper) init(t *testing.T, conf testSpec) {}

func (h goTestHelper) generatorFor(t *testing.T, conf testSpec, sourceFile string, genDir string) generator.CodeGenerator {
	source, err := ioutil.ReadFile(sourceFile)
	assert.NoErr(t, err)

	// make a copy of the default generator
	var gen = *conf.generator.(*gogenerator.GoGenerator)

	if match := goGeneratorArgsRegexp.FindSubmatch(source); len(match) > 1 {
		var args = argsToMap(string(match[1]))
		for name, value := range args {
			_ = value // get rid of the testHelper warning until we start using some options with values

			switch name {
			case "byValue":
				// TODO gen.ByValue = true
			default:
				t.Fatalf("unknown option '%s'", name)
			}
		}
	}
	return &gen
}

func argsToMap(args string) map[string]string {
	var result = map[string]string{}

	for _, arg := range strings.Split(strings.TrimSpace(args), "-") {
		arg = strings.TrimSpace(arg)

		if len(arg) == 0 {
			continue
		}

		var pair = strings.Split(arg, " ")
		if len(pair) == 1 {
			result[pair[0]] = ""
		} else {
			result[pair[0]] = pair[1]
		}
	}

	return result
}

func (goTestHelper) prepareTempDir(t *testing.T, conf testSpec, srcDir, tempDir, tempRoot string) func(err error) error {
	// copy the source dir, including the relative paths (to make sure expected errors contain same paths)
	assert.NoErr(t, copyDirectory(srcDir, tempDir, 0700, 0600))

	// When outside of the project's directory, we need to set up the whole temp dir as its own module, otherwise
	// it won't find this `objectbox-go`. Therefore, we create a go.mod file pointing it to the right path.
	cwd, err := os.Getwd()
	assert.NoErr(t, err)
	var modulePath = "example.com/virtual/objectbox-generator/test/comparison/" + srcDir
	var goMod = "module " + modulePath + "\n" +
		"replace " + goModuleName + " => " + filepath.Join(cwd, "/../../") + "\n" +
		"require " + goModuleName + " v0.0.0"
	assert.NoErr(t, ioutil.WriteFile(path.Join(tempDir, "go.mod"), []byte(goMod), 0600))

	// NOTE: we can't change directory using os.Chdir() because it applies to a process/thread, not a goroutine.
	// Therefore, we just map paths in received errors, so they match the expected ones.
	return func(err error) error {
		if err == nil {
			return nil
		}
		var str = strings.Replace(err.Error(), tempRoot+string(os.PathSeparator), "", -1)
		str = strings.Replace(str, modulePath, goModuleName+"/test/comparison/"+srcDir, -1)
		return errors.New(str)
	}
}

func (goTestHelper) build(t *testing.T, conf testSpec, dir string, expectedError error, errorTransformer func(err error) error) {
	stdOut, stdErr, err := gobuild(dir)
	if err == nil && expectedError == nil {
		// successful
		return
	}

	// On Windows, we're getting a `go finding` message during the build - remove it to be consistent.
	var reg = regexp.MustCompile("go: finding " + goModuleName + " v0.0.0[ \r\n]+")
	stdErr = reg.ReplaceAll(stdErr, nil)

	checkBuildError(t, errorTransformer, stdOut, stdErr, err, expectedError)
}

func gobuild(path string) (stdOut []byte, stdErr []byte, err error) {
	var cmd = exec.Command("go", "build")
	cmd.Dir = path
	stdOut, err = cmd.Output()
	if ee, ok := err.(*exec.ExitError); ok {
		stdErr = ee.Stderr
	}
	return
}
