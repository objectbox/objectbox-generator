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

// package main provides cmd line Main() function for objectbox-go project, actual command (executable) is part of objectbox-go project
package gogen

import (
	"flag"
	"fmt"
	"os"

	"github.com/objectbox/objectbox-generator/internal/generator"
	"github.com/objectbox/objectbox-generator/internal/generator/go"
)

const defaultErrorCode = 2

func Main() {
	path, clean, options := getArgs()

	var err error
	if clean {
		fmt.Printf("Removing ObjectBox bindings for %s\n", path)
		err = generator.Clean(options.CodeGenerator, path)
	} else {
		fmt.Printf("Generating ObjectBox bindings for %s\n", path)
		err = generator.Process(path, options)
	}

	stopOnError(0, err)
}

func stopOnError(code int, err error) {
	if err != nil {
		fmt.Println(err)
		if code == 0 {
			code = defaultErrorCode
		}
		os.Exit(code)
	}
}

func showUsage() {
	fmt.Fprint(flag.CommandLine.Output(), `Usage:
	objectbox-gogen [flags] [path-pattern]
		to generate the binding code

or

	objectbox-gogen clean [path-pattern]
		to remove the generated files instead of creating them - this removes *.obx.go and objectbox-model.go but keeps objectbox-model.json

path-pattern:
  * a path or a valid path pattern as accepted by the go tool (e.g. ./...)
  * if not given, the generator expects GOFILE environment variable to be set

Available flags:
`)
	flag.PrintDefaults()
}

func showUsageAndExit(a ...interface{}) {
	if len(a) > 0 {
		a = append(a, "\n\n")
		fmt.Fprint(flag.CommandLine.Output(), a...)
	}
	showUsage()
	os.Exit(1)
}

func getArgs() (path string, clean bool, options generator.Options) {
	var gen = &gogenerator.GoGenerator{}
	options.CodeGenerator = gen

	var printVersion bool
	var printHelp bool
	flag.Usage = showUsage
	flag.StringVar(&path, "source", "", "@deprecated, equivalent to passing the given source file path as as the path-pattern argument")
	flag.StringVar(&options.ModelInfoFile, "persist", "", "path to the model information persistence file")
	flag.BoolVar(&gen.ByValue, "byValue", false, "getters should return a struct value (a copy) instead of a struct pointer")
	flag.BoolVar(&printVersion, "version", false, "print the generator version info")
	flag.BoolVar(&printHelp, "help", false, "print this help")
	flag.Parse()

	if printHelp {
		showUsage()
		os.Exit(0)
	}

	if printVersion {
		fmt.Println(fmt.Sprintf("ObjectBox Go Generator v%s #%d", generator.Version, generator.VersionId))
		os.Exit(0)
	}

	if len(path) != 0 {
		fmt.Println("'source' flag is deprecated and will be removed in the future - use a standard positional argument instead. See command help for more information.")
	}

	var argPath string

	if flag.NArg() == 2 {
		clean = true
		if flag.Arg(0) != "clean" {
			showUsageAndExit("Unknown argument %s\n", flag.Arg(0))
		}

		argPath = flag.Arg(1)

	} else if flag.NArg() == 1 {
		argPath = flag.Arg(0)
	} else if flag.NArg() != 0 {
		showUsageAndExit()
	}

	// if the path-pattern positional argument was given
	if len(argPath) > 0 {
		if len(path) == 0 {
			path = argPath
		} else if argPath != path {
			fmt.Printf("Path argument mismatch - given 'source' flag '%s' and the positional path argument '%s'\n", path, argPath)
			showUsageAndExit()
		}
	}

	if len(path) == 0 {
		// if the command is run by go:generate some environment variables are set
		// https://golang.org/pkg/cmd/go/internal/generate/
		if gofile, exists := os.LookupEnv("GOFILE"); exists {
			path = gofile
		}

		if len(path) == 0 {
			showUsageAndExit()
		}
	}

	return
}
