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

package gogenerator

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"path/filepath"
	"strings"

	"github.com/objectbox/objectbox-generator/internal/generator"
	"github.com/objectbox/objectbox-generator/internal/generator/go/templates"
	"github.com/objectbox/objectbox-generator/internal/generator/model"
)

type GoGenerator struct {
	binding *astReader
}

// BindingFile returns a name of the binding file for the given entity file.
func (GoGenerator) BindingFile(forFile string) string {
	var extension = filepath.Ext(forFile)
	return forFile[0:len(forFile)-len(extension)] + ".obx" + extension
}

// ModelFile returns the model GO file for the given JSON info file path
func (GoGenerator) ModelFile(forFile string) string {
	var extension = filepath.Ext(forFile)
	return forFile[0:len(forFile)-len(extension)] + ".go"
}

func (GoGenerator) IsGeneratedFile(file string) bool {
	var name = filepath.Base(file)
	return name == "objectbox-model.go" || strings.HasSuffix(name, ".obx.go")
}

func (goGen *GoGenerator) ParseSource(sourceFile string) (*model.ModelInfo, error) {
	var f *file
	var err error

	if f, err = parseFile(sourceFile); err != nil {
		return nil, fmt.Errorf("can't parse file %s: %s", sourceFile, err)
	}

	if goGen.binding, err = NewBinding(); err != nil {
		return nil, fmt.Errorf("can't init Go AST reader: %s", err)
	}

	if err = goGen.binding.CreateFromAst(f); err != nil {
		return nil, fmt.Errorf("can't prepare bindings for %s: %s", sourceFile, err)
	}

	return goGen.binding.model, nil
}

func (goGen *GoGenerator) WriteBindingFiles(sourceFile string, options generator.Options, mergedModel *model.ModelInfo) error {
	var err, err2 error

	var bindingSource []byte
	if bindingSource, err = goGen.generateBindingFile(options, mergedModel); err != nil {
		return fmt.Errorf("can't generate binding file %s: %s", sourceFile, err)
	}

	var bindingFile = BindingFile(sourceFile)
	if formattedSource, err := format.Source(bindingSource); err != nil {
		// we just store error but still write the file so that we can check it manually
		err2 = fmt.Errorf("failed to format generated binding file %s: %s", bindingFile, err)
	} else {
		bindingSource = formattedSource
	}

	if err = generator.WriteFile(bindingFile, bindingSource, sourceFile); err != nil {
		return fmt.Errorf("can't write binding file %s: %s", sourceFile, err)
	} else if err2 != nil {
		// now when the binding has been written (for debugging purposes), we can return the error
		return err2
	}

	return nil
}

func (goGen *GoGenerator) generateBindingFile(options generator.Options, m *model.ModelInfo) (data []byte, err error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	var tplArguments = struct {
		Model            *model.ModelInfo
		Binding          *astReader
		GeneratorVersion int
		Options          generator.Options
	}{m, goGen.binding, generator.VersionId, options}

	if err = templates.BindingTemplate.Execute(writer, tplArguments); err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	return b.Bytes(), nil
}

func (goGen *GoGenerator) WriteModelBindingFile(options generator.Options, modelInfo *model.ModelInfo) error {
	var err, err2 error

	var modelFile = ModelFile(options.ModelInfoFile)
	var modelSource []byte

	if modelSource, err = goGen.generateModelFile(modelInfo); err != nil {
		return fmt.Errorf("can't generate model file %s: %s", modelFile, err)
	}

	if formattedSource, err := format.Source(modelSource); err != nil {
		// we just store error but still writ the file so that we can check it manually
		err2 = fmt.Errorf("failed to format generated model file %s: %s", modelFile, err)
	} else {
		modelSource = formattedSource
	}

	if err = generator.WriteFile(modelFile, modelSource, options.ModelInfoFile); err != nil {
		return fmt.Errorf("can't write model file %s: %s", modelFile, err)
	} else if err2 != nil {
		// now when the model has been written (for debugging purposes), we can return the error
		return err2
	}

	return nil
}

func (goGen *GoGenerator) generateModelFile(m *model.ModelInfo) (data []byte, err error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	var tplArguments = struct {
		Package          string
		Model            *model.ModelInfo
		GeneratorVersion int
	}{goGen.binding.Package.Name(), m, generator.VersionId}

	if err = templates.ModelTemplate.Execute(writer, tplArguments); err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	return b.Bytes(), nil
}
