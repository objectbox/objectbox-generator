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

// package cmake provides tools to create, configure and build C & C++ projects using CMake
package cmake

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"text/template"
)

// Cmake contains all configuration necessary to configure and build a CMake project
type Cmake struct {
	// Configs required for CMakeLists.txt
	Name           string // Executable name
	IsCpp          bool
	Standard       int // CMAKE_C/CXX_STANDARD
	Files          []string
	IncludeDirs    []string
	LinkLibs       []string // Library names or full paths
	LinkDirs       []string // Where should the linker look for libraries
	Generator      string
	ConfigureFlags []string

	// Build configuration
	ConfDir  string
	BuildDir string

	tempRoot string
}

// CreateTempDirs creates temporary directories for conf and build dir
func (cmake *Cmake) CreateTempDirs() error {
	if len(cmake.tempRoot) != 0 {
		return errors.New("temp root is already set")
	}

	tempRoot, err := ioutil.TempDir("", cmake.Name+"cmake")
	if err != nil {
		return err
	}

	buildDir, err := createTempDir(tempRoot, "build")
	confDir, err2 := createTempDir(tempRoot, "conf")

	if err != nil || err2 != nil {
		os.RemoveAll(cmake.tempRoot)
		if err != nil {
			return err
		}
		return err2
	}

	cmake.tempRoot = tempRoot
	cmake.BuildDir = buildDir
	cmake.ConfDir = confDir
	return nil
}

func createTempDir(parent, name string) (string, error) {
	var path = filepath.Join(parent, name)
	if err := os.Mkdir(path, 0700); err != nil {
		return "", err
	}
	return path, nil
}

func (cmake *Cmake) RemoveTempDirs() error {
	if len(cmake.tempRoot) == 0 {
		return nil
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

func (cmake *Cmake) GetCMakeListsTxt() (string, error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	if err := cmakeListsTpl.Execute(writer, cmake); err != nil {
		return "", err
	} else if err = writer.Flush(); err != nil {
		return "", err
	}
	return string(b.Bytes()), nil
}

var cmakeListsTpl = template.Must(template.New("CMakeLists.txt").
	Funcs(template.FuncMap{
		"Join": func(ss []string) string {
			switch len(ss) {
			case 0:
				return ""
			case 1:
				return filepath.ToSlash(ss[0])
			default:
				var result string
				for _, s := range ss {
					result = result + "\n\t" + filepath.ToSlash(s)
				}
				return result
			}
		},
	}).
	Parse(`
cmake_minimum_required(VERSION 3.5)
{{if .Standard}}set(CMAKE_C{{if .IsCpp}}XX{{end}}_STANDARD {{.Standard}}){{end}}
project({{.Name}} C{{if .IsCpp}}XX{{end}})

add_executable(${PROJECT_NAME} {{Join .Files}})
{{if .IncludeDirs}}target_include_directories(${PROJECT_NAME} PRIVATE {{Join .IncludeDirs}}){{end}}
{{if .LinkLibs}}target_link_libraries(${PROJECT_NAME} PRIVATE {{Join .LinkLibs}}){{end}}
{{if .LinkDirs}}target_link_directories(${PROJECT_NAME} PRIVATE {{Join .LinkDirs}}){{end}}
`))

// Configure runs cmake configuration step.
func (cmake *Cmake) Configure() ([]byte, []byte, error) {
	if len(cmake.Generator) == 0 && runtime.GOOS == "windows" {
		// Using MinGW because MSVC doesn't support linking to .dll - app is supposed to load them on runtime
		cmake.Generator = "MinGW Makefiles"
	}

	var args = []string{cmake.ConfDir}
	args = append(args, cmake.ConfigureFlags[:]...)

	if len(cmake.Generator) > 0 && cmake.Generator != "default" {
		args = append(args, "-G", cmake.Generator)
	}

	return cmakeExec(cmake.BuildDir, args[:]...)
}

// Configure runs cmake configuration step without any CI quirks as is.
func (cmake *Cmake) ConfigureRaw() ([]byte, []byte, error) {
	var args = []string{cmake.ConfDir}
	args = append(args, cmake.ConfigureFlags[:]...)

	if len(cmake.Generator) > 0 && cmake.Generator != "default" {
		args = append(args, "-G", cmake.Generator)
	}

	return cmakeExec(cmake.BuildDir, args[:]...)
}

// Build runs cmake build "Name" target step.
func (cmake *Cmake) BuildTarget() ([]byte, []byte, error) {
	return cmakeExec(cmake.ConfDir,
		"--build", cmake.BuildDir,
		"--target", cmake.Name,
		"--parallel "+strconv.FormatInt(int64(runtime.NumCPU()/2), 10),
	)
}

// Build runs cmake build "Name" target step.
func (cmake *Cmake) BuildWithTarget(target string) ([]byte, []byte, error) {
	return cmakeExec(cmake.ConfDir,
		"--build", cmake.BuildDir,
		"--target", target,
		"--parallel "+strconv.FormatInt(int64(runtime.NumCPU()/2), 10),
	)
}

// Build runs cmake build step.
func (cmake *Cmake) BuildDefaults() ([]byte, []byte, error) {
	return cmakeExec(cmake.ConfDir,
		"--build", cmake.BuildDir,
		"--parallel "+strconv.FormatInt(int64(runtime.NumCPU()/2), 10),
	)
}

// Build runs cmake default build step with config
func (cmake *Cmake) BuildDefaultsWithConfig(config string) ([]byte, []byte, error) {
	return cmakeExec(cmake.ConfDir,
		"--build", cmake.BuildDir,
		"--config", config,
		"--parallel "+strconv.FormatInt(int64(runtime.NumCPU()/2), 10),
	)
}

// Build runs cmake default build step with config
func (cmake *Cmake) BuildTargetWithConfig(config, target string) ([]byte, []byte, error) {
	return cmakeExec(cmake.ConfDir,
		"--build", cmake.BuildDir,
		"--config", config,
		"--target", target,
		"--parallel "+strconv.FormatInt(int64(runtime.NumCPU()/2), 10),
	)
}

func CopyDir(cwd string, src string, dst string) ([]byte, []byte, error) {
	return cmakeExec(cwd, "-E", "copy_directory", src, dst)
}

func CopyFile(cwd string, file string, dst string) ([]byte, []byte, error) {
	return cmakeExec(cwd, "-E", "copy", file, dst)
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
