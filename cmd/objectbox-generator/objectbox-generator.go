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

	"github.com/objectbox/objectbox-generator/cmd"
	"github.com/objectbox/objectbox-generator/internal/generator"
	"github.com/objectbox/objectbox-generator/internal/generator/c"
	"github.com/objectbox/objectbox-generator/internal/generator/flatbuffersc"
	"github.com/objectbox/objectbox-generator/internal/generator/go"
)

func main() {
	if runFlatcIfRequested() {
		return
	}

	generatorcmd.Main(&command{})
}

// implements generatorcmd.generatorCommand
type command struct {
	langs map[string]*bool
}

func (cmd command) ShowUsage() {
	fmt.Fprint(flag.CommandLine.Output(), `Usage:
  objectbox-generator [flags] {source-file}
      to generate the binding code

or
  objectbox-generator [flags] clean {path}
      to remove the generated files instead of creating them - this removes *.obx.h and objectbox-model.h but keeps objectbox-model.json

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
	cmd.langs["cpp"] = flag.Bool("cpp", false, "generate C++ code")
	cmd.langs["go"] = flag.Bool("go", false, "generate Go code")
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

	switch selectedLang {
	case "go":
		options.CodeGenerator = &gogenerator.GoGenerator{}
	case "c":
		options.CodeGenerator = &cgenerator.CGenerator{
			PlainC: true,
		}
	case "cpp":
		options.CodeGenerator = &cgenerator.CGenerator{
			PlainC: false,
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
