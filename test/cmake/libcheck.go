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

package cmake

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// LibraryExists tries to compile a simple program linking to the given library
func LibraryExists(name string, includeFiles, includeDirs, linkDirs, predefines []string) error {
	build := Cmake{
		Name:        "check-" + name,
		IsCpp:       true,
		Standard:    11,
		Files:       []string{"main.cpp"},
		IncludeDirs: includeDirs,
		LinkDirs:    linkDirs,
	}
	if err := build.CreateTempDirs(); err != nil {
		return err
	}
	defer build.RemoveTempDirs()

	if len(name) > 0 {
		build.LinkLibs = []string{name}
	}

	if err := build.WriteCMakeListsTxt(); err != nil {
		return err
	}

	{ // write main.cpp
		mainPath := filepath.Join(build.ConfDir, build.Files[0])
		var mainSrc string
		if len(predefines) > 0 {
			for _, predefine := range predefines {
				mainSrc = mainSrc + "#define " + predefine + "\n"
			}
		}
		if len(includeFiles) > 0 {
			for _, inc := range includeFiles {
				mainSrc = mainSrc + "#include <" + inc + ">\n"
			}
		}
		mainSrc = mainSrc + "\nint main(){ return 0; }\n\n"
		if err := ioutil.WriteFile(mainPath, []byte(mainSrc), 0600); err != nil {
			return err
		}
	}

	if stdOut, stdErr, err := build.Configure(); err != nil {
		return fmt.Errorf("cmake configuration failed: \n%s\n%s\n%s", stdOut, stdErr, err)
	}

	if stdOut, stdErr, err := build.BuildTarget(); err != nil {
		return fmt.Errorf("cmake build failed: \n%s\n%s\n%s", stdOut, stdErr, err)
	}

	return nil
}
