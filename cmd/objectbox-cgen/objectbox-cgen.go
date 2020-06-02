/*
 * Copyright 2019 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/objectbox/objectbox-go/internal/generator"
	"github.com/objectbox/objectbox-go/internal/generator/c"
)

func main() {
	path, clean, options := getArgs()

	var err error
	if clean {
		fmt.Printf("Removing ObjectBox bindings for %s\n", path)
		err = generator.Clean(options.CodeGenerator, path)
	} else {
		fmt.Printf("Generating ObjectBox bindings for %s\n", path)
		err = generator.Process(path, options)
	}

	stopOnError(err)
}

func stopOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func showUsage() {
	fmt.Fprint(flag.CommandLine.Output(), `Usage:
	objectbox-cgen [flags] [path-pattern]
		to generate the binding code

or

	objectbox-cgen clean [path-pattern]
		to remove the generated files instead of creating them - this removes *.obx.h and objectbox-model.h but keeps objectbox-model.json

path-pattern:
  * a path or a valid path pattern (e.g. ./...)
  
Available flags:`)
	flag.PrintDefaults()
}

func showUsageAndExit() {
	showUsage()
	os.Exit(1)
}

func getArgs() (path string, clean bool, options generator.Options) {
	var gen = &cgenerator.CGenerator{}
	options.CodeGenerator = gen

	var printVersion bool
	var printHelp bool
	flag.Usage = showUsage
	flag.StringVar(&gen.OutPath, "out", "", "output path for generated source files")
	flag.StringVar(&options.ModelInfoFile, "persist", "", "path to the model information persistence file (JSON)")
	flag.BoolVar(&printVersion, "version", false, "print the generator version info")
	flag.BoolVar(&gen.PlainC, "c", false, "generate plain C code instead of default C++")
	flag.BoolVar(&printHelp, "help", false, "print this help")
	flag.Parse()

	if printHelp {
		showUsage()
		os.Exit(0)
	}

	if printVersion {
		fmt.Println(fmt.Sprintf("ObjectBox C/C++ binding code generator version: %d", generator.Version))
		os.Exit(0)
	}

	if flag.NArg() == 2 {
		clean = true
		if flag.Arg(0) != "clean" {
			fmt.Printf("Unknown argument %s", flag.Arg(0))
			showUsageAndExit()
		}

		path = flag.Arg(1)

	} else if flag.NArg() == 1 {
		path = flag.Arg(0)
	} else if flag.NArg() != 0 {
		showUsageAndExit()
	}

	if len(path) == 0 {
		showUsageAndExit()
	}

	return
}
