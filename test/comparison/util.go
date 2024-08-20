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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-generator/v4/test/assert"
)

func checkBuildError(t *testing.T, errorTransformer func(err error) error, stdOut []byte, stdErr []byte, err error, expectedError error) {
	if err == nil {
		if expectedError != nil {
			assert.Failf(t, "Unexpected PASS during compilation, expected error: %s", expectedError)
		}
		return
	}

	var receivedError = errorTransformer(fmt.Errorf("%s\n%s\n%s", stdOut, stdErr, err))

	// Fix paths in the error output on Windows so that it matches the expected error (which always uses '/').
	if os.PathSeparator != '/' && expectedError != nil {
		// Make sure the expected error doesn't contain the path separator already - to make it easier to debug.
		if strings.Contains(expectedError.Error(), string(os.PathSeparator)) {
			assert.Failf(t, "compile-error.expected contains this OS path separator '%v' so paths can't be normalized to '/'", string(os.PathSeparator))
		}
		receivedError = errors.New(strings.Replace(receivedError.Error(), string(os.PathSeparator), "/", -1))
	}

	assert.Eq(t, expectedError, receivedError)
}

func repoRoot(t *testing.T) string {
	cwd, err := os.Getwd()
	assert.NoErr(t, err)
	return filepath.ToSlash(filepath.Join(cwd, "..", ".."))
}
