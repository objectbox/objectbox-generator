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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-generator/test/assert"
	"github.com/objectbox/objectbox-generator/test/build"
)

type goTestHelper struct {
}

func (goTestHelper) prepareTempDir(t *testing.T, srcDir, tempDir, tempRoot string) func(err error) error {
	// When outside of the project's directory, we need to set up the whole temp dir as its own module, otherwise
	// it won't find this `objectbox-go`. Therefore, we create a go.mod file pointing it to the right path.
	cwd, err := os.Getwd()
	assert.NoErr(t, err)
	var modulePath = "example.com/virtual/objectbox-generator/test/comparison/" + srcDir
	var goMod = "module " + modulePath + "\n" +
		"replace " + moduleName + " => " + filepath.Join(cwd, "/../../") + "\n" +
		"require " + moduleName + " v0.0.0"
	assert.NoErr(t, ioutil.WriteFile(path.Join(tempDir, "go.mod"), []byte(goMod), 0600))

	// NOTE: we can't change directory using os.Chdir() because it applies to a process/thread, not a goroutine.
	// Therefore, we just map paths in received errors, so they match the expected ones.
	return func(err error) error {
		if err == nil {
			return nil
		}
		var str = strings.Replace(err.Error(), tempRoot+string(os.PathSeparator), "", -1)
		str = strings.Replace(str, modulePath, moduleName+"/test/comparison/"+srcDir, -1)
		return errors.New(str)
	}
}

func (goTestHelper) build(t *testing.T, dir string, errorTransformer func(err error) error) {
	var expectedError error
	if fileExists(path.Join(dir, "compile-error.expected")) {
		content, err := ioutil.ReadFile(path.Join(dir, "compile-error.expected"))
		assert.NoErr(t, err)
		expectedError = errors.New(string(content))
	}

	stdOut, stdErr, err := build.Package(dir)
	if err == nil && expectedError == nil {
		// successful
		return
	}

	if err == nil && expectedError != nil {
		assert.Failf(t, "Unexpected PASS during compilation")
	}

	// On Windows, we're getting a `go finding` message during the build - remove it to be consistent.
	var reg = regexp.MustCompile("go: finding " + moduleName + " v0.0.0[ \r\n]+")
	stdErr = reg.ReplaceAll(stdErr, nil)

	var receivedError = errorTransformer(fmt.Errorf("%s\n%s\n%s", stdOut, stdErr, err))

	// Fix paths in the error output on Windows so that it matches the expected error (which always uses '/').
	if os.PathSeparator != '/' {
		// Make sure the expected error doesn't contain the path separator already - to make it easier to debug.
		if strings.Contains(expectedError.Error(), string(os.PathSeparator)) {
			assert.Failf(t, "compile-error.expected contains this OS path separator '%v' so paths can't be normalized to '/'", string(os.PathSeparator))
		}
		receivedError = errors.New(strings.Replace(receivedError.Error(), string(os.PathSeparator), "/", -1))
	}

	assert.Eq(t, expectedError, receivedError)
}
