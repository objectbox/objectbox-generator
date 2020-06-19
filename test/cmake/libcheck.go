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

package cmake

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// LibraryExists tries to compile a simple program linking to the given library
func LibraryExists(name, includeFile string) (bool, error) {
	build := Cmake{
		Name:  "check-" + name,
		IsCpp: true,
		Files: []string{"main.cpp"},
	}
	if err := build.CreateTempDirs(); err != nil {
		return false, err
	}
	defer build.RemoveTempDirs()

	if len(name) > 0 {
		build.LinkLibs = []string{name}
	}

	if err := build.WriteCMakeListsTxt(); err != nil {
		return false, err
	}

	{ // write main.cpp
		mainPath := filepath.Join(build.SourceDir, build.Files[0])
		mainSrc := "int main(){ return 0; }\n\n"
		if len(includeFile) > 0 {
			mainSrc = "#include <" + includeFile + ">\n" + mainSrc
		}
		if err := ioutil.WriteFile(mainPath, []byte(mainSrc), 0600); err != nil {
			return false, err
		}
	}

	if stdOut, stdErr, err := build.Configure(); err != nil {
		return false, fmt.Errorf("cmake configuration failed: \n%s\n%s\n%s", stdOut, stdErr, err)
	}

	if stdOut, stdErr, err := build.Build(); err != nil {
		return false, fmt.Errorf("cmake build failed: \n%s\n%s\n%s", stdOut, stdErr, err)
	}

	return true, nil
}
