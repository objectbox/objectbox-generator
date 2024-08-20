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

// Package gogen provides cmd line Main() function for objectbox-go project, actual command (executable) is part of objectbox-go project.
// We're keeping this code here instead of objectbox-go  to minimize the exposed API (avoid exposing Options and GoGenerator symbols).
package gogen

import (
	"flag"
	"fmt"
	"os"

	generatorcmd "github.com/objectbox/objectbox-generator/v4/cmd"
	"github.com/objectbox/objectbox-generator/v4/internal/generator"
	gogenerator "github.com/objectbox/objectbox-generator/v4/internal/generator/go"
)

const VersionId = generator.VersionId

func Main() {
	generatorcmd.Main(&command{})
}

// implements generatorcmd.generatorCommand
type command struct {
	byValue bool
}

func (cmd command) ShowUsage() {
	fmt.Fprint(flag.CommandLine.Output(), `Usage:
	objectbox-gogen [flags] {source-file}
		to generate the binding code

or

	objectbox-gogen clean {path}
		to remove the generated files instead of creating them - this removes *.obx.go and objectbox-model.go but keeps objectbox-model.json

path:
  * a source file path or a valid path pattern as accepted by the go tool (e.g. ./...)
  * if not given, the generator expects GOFILE environment variable to be set

Available flags:
`)
	flag.PrintDefaults()
}

func (cmd *command) ConfigureFlags() {
	flag.BoolVar(&cmd.byValue, "byValue", false, "getters should return a struct value (a copy) instead of a struct pointer")
}

func (cmd *command) ParseFlags(remainingPosArgs *[]string, options *generator.Options) error {
	options.CodeGenerator = &gogenerator.GoGenerator{
		ByValue: cmd.byValue,
	}

	if len(options.InPath) == 0 {
		// if the command is run by go:generate some environment variables are set
		// https://golang.org/pkg/cmd/go/internal/generate/
		if gofile, exists := os.LookupEnv("GOFILE"); exists {
			options.InPath = gofile
		}
	}

	return nil
}
