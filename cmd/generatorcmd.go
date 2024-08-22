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

// Package generatorcmd provides common functionality for code-generator executables.
// Generates objectbox related code by reading models (e.g. .fbs schemas, .go files).
// Currently support generation of C, C++ and Go code.
package generatorcmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/objectbox/objectbox-generator/v4/internal/generator"
)

const defaultErrorCode = 2

// / generatorCommand defines an interface for command-line applications to implement
type generatorCommand interface {
	ShowUsage()
	ConfigureFlags()
	ParseFlags(remainingPosArgs *[]string, options *generator.Options) error
}

func Main(impl generatorCommand) {
	clean, options := getArgs(impl)

	var err error
	if clean {
		fmt.Printf("Removing ObjectBox bindings for %s\n", options.InPath)
		err = generator.Clean(options.CodeGenerator, options.InPath)
	} else {
		fmt.Printf("Generating ObjectBox bindings for %s\n", options.InPath)
		err = generator.Process(options)
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

func showUsageAndExit(impl generatorCommand, a ...interface{}) {
	if len(a) > 0 {
		a = append(a, "\n\n")
		fmt.Fprint(flag.CommandLine.Output(), a...)
	}
	impl.ShowUsage()
	os.Exit(1)
}

func getArgs(impl generatorCommand) (clean bool, options generator.Options) {
	var printVersion bool
	var printHelp bool
	flag.Usage = impl.ShowUsage
	impl.ConfigureFlags()
	flag.StringVar(&options.OutPath, "out", "", "output path for generated source files")
	flag.StringVar(&options.OutHeadersPath, "out-headers", "", "optional: output path for generated header files") // opt-in: C and C++
	flag.StringVar(&options.ModelInfoFile, "model", "", "path to the model information persistence file (JSON)")
	// TODO remove in v0.15.0 or later
	flag.StringVar(&options.ModelInfoFile, "persist", "", "[DEPRECATED, use 'model'] path to the model information persistence file (JSON)")
	flag.BoolVar(&printVersion, "version", false, "print the generator version info")
	flag.BoolVar(&printHelp, "help", false, "print this help")
	flag.Parse()

	if printHelp {
		impl.ShowUsage()
		os.Exit(0)
	}

	if printVersion {
		fmt.Println(fmt.Sprintf("ObjectBox Generator v%s #%d", generator.Version, generator.VersionId))
		os.Exit(0)
	}

	// process positional args
	var args = flag.Args()

	if len(args) > 0 && args[0] == "clean" {
		clean = true
		args = args[1:]
	}

	if len(args) > 0 {
		options.InPath = args[0]
		args = args[1:]
	}

	if err := impl.ParseFlags(&args, &options); err != nil {
		showUsageAndExit(impl, err)
	}

	if len(options.InPath) == 0 {
		showUsageAndExit(impl, "path not specified")
	}

	if len(args) > 0 {
		showUsageAndExit(impl, "unknown arguments", args)
	}

	return
}
