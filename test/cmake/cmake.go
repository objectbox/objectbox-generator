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

// package cmake provides tools to create, configure and build C & C++ projects using CMake
package cmake

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Cmake struct {
	// Configs required for CMakeLists.txt
	Name        string // Executable name
	IsCpp       bool
	Standard    int // CMAKE_C/CXX_STANDARD
	Files       []string
	IncludeDirs []string
	LinkLibs    []string // Library names or full paths
	LinkDirs    []string // Where should the linker look for libraries

	// Build configuration
	SourceDir string
	ConfDir   string
	BuildDir  string

	tempRoot string
}

// CreateCmake creates new cmake configuration.
// If useTempDirs is true, ConfDir and BuildDir are created as temporary directory
func CreateCmake(name, sourceDir string, useTempDirs bool) (Cmake, error) {
	var cmake = Cmake{
		Name:      name,
		SourceDir: sourceDir,
	}

	if useTempDirs {
		var err error
		if cmake.tempRoot, err = ioutil.TempDir("", name+"cmake"); err != nil {
			return cmake, err
		}

		cmake.BuildDir = filepath.Join(cmake.tempRoot, "build")
		cmake.ConfDir = filepath.Join(cmake.tempRoot, "conf")
		if err = os.Mkdir(cmake.BuildDir, 0700); err != nil {
			os.RemoveAll(cmake.tempRoot)
			return cmake, err
		}
		if err = os.Mkdir(cmake.ConfDir, 0700); err != nil {
			os.RemoveAll(cmake.tempRoot)
			return cmake, err
		}
	}

	return cmake, nil
}

func (cmake *Cmake) RemoveTempDirs() error {
	if len(cmake.tempRoot) == 0 {
		return errors.New("temp dirs were not used")
	}
	return os.RemoveAll(cmake.tempRoot)
}

// WriteCMakeListsTxt writes cmake specification to file and returns the file name.
func (cmake *Cmake) WriteCMakeListsTxt() error {
	// open the file, overwrite if it exists
	file, err := os.OpenFile(filepath.Join(cmake.ConfDir, "CMakeLists.txt"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if err := cmakeListsTpl.Execute(writer, cmake); err != nil {
		return err
	}

	return writer.Flush()
}

var cmakeListsTpl = template.Must(template.New("CMakeLists.txt").
	Funcs(template.FuncMap{
		"Join": func(ss []string) string {
			return strings.Join(ss, " ")
		},
	}).
	Parse(`
cmake_minimum_required(VERSION 3.0)
set(CMAKE_C{{if .IsCpp}}XX{{end}}_STANDARD {{.Standard}})
project({{.Name}} C{{if .IsCpp}}XX{{end}})

add_executable(${PROJECT_NAME} {{Join .Files}})
{{if .IncludeDirs}}target_include_directories(${PROJECT_NAME} PRIVATE {{Join .IncludeDirs}}){{end}}
{{if .LinkLibs}}target_link_libraries(${PROJECT_NAME} {{Join .LinkLibs}}){{end}}
{{if .LinkDirs}}target_link_directories(${PROJECT_NAME} {{Join .LinkDirs}}){{end}}
`))

// Configure runs cmake configuration step.
func (cmake *Cmake) Configure() ([]byte, []byte, error) {
	return cmakeExec(cmake.BuildDir, cmake.ConfDir)
}

// Configure runs cmake build step.
func (cmake *Cmake) Build() ([]byte, []byte, error) {
	return cmakeExec(cmake.SourceDir, "--build", cmake.BuildDir)
}

func cmakeExec(cwd string, args ...string) (stdOut []byte, stdErr []byte, err error) {
	var cmd = exec.Command("cmake", args...)
	cmd.Dir = cwd
	stdOut, err = cmd.Output()
	if ee, ok := err.(*exec.ExitError); ok {
		stdErr = ee.Stderr
	}
	return
}
