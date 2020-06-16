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

package cgenerator

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/objectbox/objectbox-generator/internal/generator"
	"github.com/objectbox/objectbox-generator/internal/generator/c/templates"
	"github.com/objectbox/objectbox-generator/internal/generator/flatbuffersc"
	"github.com/objectbox/objectbox-generator/internal/generator/model"
)

type CGenerator struct {
	OutPath string
	PlainC  bool
}

// BindingFile returns a name of the binding file for the given entity source file.
func (gen *CGenerator) BindingFile(forFile string) string {
	if len(gen.OutPath) > 0 {
		forFile = filepath.Join(gen.OutPath, filepath.Base(forFile))
	}
	var extension = filepath.Ext(forFile)
	var suffix string
	if !gen.PlainC {
		suffix = "-cpp"
	}
	return forFile[0:len(forFile)-len(extension)] + suffix + ".obx.h"
}

// ModelFile returns the model GO file for the given JSON info file path
func (gen *CGenerator) ModelFile(forFile string) string {
	if len(gen.OutPath) > 0 {
		forFile = filepath.Join(gen.OutPath, filepath.Base(forFile))
	}
	var extension = filepath.Ext(forFile)
	return forFile[0:len(forFile)-len(extension)] + ".h"
}

func (CGenerator) IsGeneratedFile(file string) bool {
	var name = filepath.Base(file)
	return name == "objectbox-model.h" || strings.HasSuffix(name, ".obx.h")
}

func (gen *CGenerator) ParseSource(sourceFile string) (*model.ModelInfo, error) {
	schemaReflection, err := flatbuffersc.ParseSchemaFile(sourceFile)
	if err != nil {
		return nil, err // already includes file name so no more context should be necessary
	}

	reader := fbSchemaReader{model: &model.ModelInfo{}}
	if err = reader.read(schemaReflection); err != nil {
		return nil, fmt.Errorf("error generating model from schema %s: %s", sourceFile, err)
	}

	return reader.model, nil
}

func (gen *CGenerator) WriteBindingFiles(sourceFile string, _ generator.Options, mergedModel *model.ModelInfo) error {
	var err, err2 error

	var bindingFile = gen.BindingFile(sourceFile)

	var bindingSource []byte
	if bindingSource, err = gen.generateBindingFile(bindingFile, mergedModel); err != nil {
		return fmt.Errorf("can't generate binding file %s: %s", sourceFile, err)
	}

	if formattedSource, err := format(bindingSource); err != nil {
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

func (gen *CGenerator) generateBindingFile(bindingFile string, m *model.ModelInfo) (data []byte, err error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	var replaceSpecialChars = strings.NewReplacer("-", "_", ".", "_")
	var ifdefGuard = strings.ToUpper(filepath.Base(bindingFile))
	ifdefGuard = replaceSpecialChars.Replace(ifdefGuard)

	var tplArguments = struct {
		Model            *model.ModelInfo
		GeneratorVersion int
		IfdefGuard       string
	}{m, generator.VersionId, ifdefGuard}

	var tpl = templates.CppBindingTemplate
	if gen.PlainC {
		tpl = templates.CBindingTemplate
	}

	if err = tpl.Execute(writer, tplArguments); err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	return b.Bytes(), nil
}

func (gen *CGenerator) WriteModelBindingFile(options generator.Options, mergedModel *model.ModelInfo) error {
	var err, err2 error

	var modelFile = gen.ModelFile(options.ModelInfoFile)
	var modelSource []byte

	if modelSource, err = generateModelFile(mergedModel); err != nil {
		return fmt.Errorf("can't generate model file %s: %s", modelFile, err)
	}

	if formattedSource, err := format(modelSource); err != nil {
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

func generateModelFile(m *model.ModelInfo) (data []byte, err error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	var tplArguments = struct {
		Model            *model.ModelInfo
		GeneratorVersion int
	}{m, generator.VersionId}

	if err = templates.ModelTemplate.Execute(writer, tplArguments); err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	return b.Bytes(), nil
}

func format(source []byte) ([]byte, error) {
	// replace tabs with spaces
	source = bytes.ReplaceAll(source, []byte("\t"), []byte("    "))

	// TODO c source formatting
	return source, nil
}
