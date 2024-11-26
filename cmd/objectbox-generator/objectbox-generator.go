/*
 * ObjectBox Generator - a build time tool for ObjectBox
 * Copyright (C) 2018-2024 ObjectBox Ltd. All rights reserved.
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

// Package main provides objectbox-generator executable.
// Generates objectbox related code by reading models (e.g. .fbs schemas, .go files).
// Currently support generation of C, C++ and Go code.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	generatorcmd "github.com/objectbox/objectbox-generator/v4/cmd"
	"github.com/objectbox/objectbox-generator/v4/internal/generator"
	cgenerator "github.com/objectbox/objectbox-generator/v4/internal/generator/c"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/flatbuffersc"
	gogenerator "github.com/objectbox/objectbox-generator/v4/internal/generator/go"
	jsgenerator "github.com/objectbox/objectbox-generator/v4/internal/generator/js"
)

func main() {
	if runFlatcIfRequested() {
		return
	}

	generatorcmd.Main(&command{})
}

// implements generatorcmd.generatorCommand
type command struct {
	langs                map[string]*bool
	optional             *string
	empty_string_as_null *bool // pointers due to flag API (https://pkg.go.dev/flag#Bool)
	nan_as_null          *bool
}

func (cmd command) ShowUsage() {
	fmt.Fprint(flag.CommandLine.Output(), `Usage:
  objectbox-generator [flags] {path}
      * to execute "clean" action (see below) on the path, removing previously generated code and missing entities,
      * and execute code generation on the path afterwards.

      The given {path} can be one of the following:
        * a directory - a non-recursive clean and generation is performed on the given directory,
        * a glob path pattern (e.g. contains a "*") - performs clean and generation on the matching paths,
        * a Go-style path pattern (e.g. "./..." - a recursive match of the current dir) - performs clean and generation on the matching paths,


or
  objectbox-generator [flags] {model/file/path.fbs}
      to generate the binding code for a single file


or
  objectbox-generator [flags] clean {path}
      to remove the generated files instead of creating them - this removes *.obx.* and objectbox-model.h but keeps objectbox-model.json

or
  objectbox-generator FLATC [flatc arguments]
      to execute FlatBuffers flatc command line tool Any arguments after the FLATC keyword are passed through.

path:
  * a source file path or a valid path pattern (e.g. ./...)
  
Available flags:
`)
	flag.PrintDefaults()
}

func (cmd *command) ConfigureFlags() {
	cmd.langs = make(map[string]*bool)
	cmd.langs["c"] = flag.Bool("c", false, "generate plain C code")
	cmd.langs["cpp"] = flag.Bool("cpp", false, "generate C++ code (at least C++14)")
	cmd.langs["cpp11"] = flag.Bool("cpp11", false, "generate C++11 code")
	cmd.langs["js"] = flag.Bool("js", false, "generate JS code")
	cmd.langs["go"] = flag.Bool("go", false, "generate Go code")

	// for c++ generator
	cmd.optional = flag.String("optional", "", "C++ wrapper type to use for fields annotated \"optional\"; one of: std::optional, std::unique_ptr, std::shared_ptr")
	cmd.empty_string_as_null = flag.Bool("empty-string-as-null", false, "C++: empty strings are treated as 0 (null)")
	cmd.nan_as_null = flag.Bool("nan-as-null", false, "C++: NaNs are treated as 0 (null)")
}

func (cmd *command) ParseFlags(remainingPosArgs *[]string, options *generator.Options) error {
	var selectedLang string
	for lang, val := range cmd.langs {
		if *val {
			if len(selectedLang) != 0 {
				return fmt.Errorf("only one output language can be specified at the moment, you've selected %s and %s", selectedLang, lang)
			}
			selectedLang = lang
		}
	}

	if len(*cmd.optional) != 0 && selectedLang != "cpp" {
		return errors.New("argument -optional is only allowed in combination with -cpp")
	}

	switch selectedLang {
	case "go":
		options.CodeGenerator = &gogenerator.GoGenerator{}
	case "c":
		options.CodeGenerator = &cgenerator.CGenerator{
			PlainC:      true,
			LangVersion: -1,    // unspecified, take the default
			Optional:    "ptr", // dummy value for checks to evaluate to true if "optional" annotation is used
		}
	case "cpp":
		options.CodeGenerator = &cgenerator.CGenerator{
			PlainC:            false,
			LangVersion:       14,
			Optional:          *cmd.optional,
			EmptyStringAsNull: *cmd.empty_string_as_null,
			NaNAsNull:         *cmd.nan_as_null,
		}
	case "cpp11":
		options.CodeGenerator = &cgenerator.CGenerator{
			PlainC:            false,
			LangVersion:       11,
			Optional:          *cmd.optional,
			EmptyStringAsNull: *cmd.empty_string_as_null,
			NaNAsNull:         *cmd.nan_as_null,
		}
	case "js":
		options.CodeGenerator = &jsgenerator.JSGenerator{
			Optional:          *cmd.optional,
			EmptyStringAsNull: *cmd.empty_string_as_null,
			NaNAsNull:         *cmd.nan_as_null,
		}
	default:
		return errors.New("you must specify an output language")
	}
	return nil
}

// runFlatcIfRequested checks command line arguments and if they start with FLATC, executes flatc compiler with the remainder of the arguments
func runFlatcIfRequested() bool {
	if len(os.Args) < 2 || strings.ToLower(os.Args[1]) != "flatc" {
		return false
	}

	code, err := flatbuffersc.ExecuteFlatc(os.Args[2:])
	if err != nil {
		fmt.Println(err)
		if code == 0 {
			code = 1
		}
		os.Exit(code)
	}
	return true
}
