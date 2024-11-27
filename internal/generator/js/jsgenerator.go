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

package jsgenerator

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/objectbox/objectbox-generator/v4/internal/generator"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/flatbuffersc"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/js/templates"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/model"
)

// JS generator, given a .fbs and an optional *model.json file, is responsible for generating:
// - objectbox-model.js
// - sche
type JSGenerator struct {
	Optional          string // std::optional, std::unique_ptr, std::shared_ptr
	EmptyStringAsNull bool
	NaNAsNull         bool
}

// Return the names of the generated JS binding file (only one!) for the given entity file.
// For example: given a schema.fbs file, outputs schema.obx.fbs.
func (gen *JSGenerator) BindingFiles(forFile string, options generator.Options) []string {

	if len(options.OutPath) > 0 {
		forFile = filepath.Join(options.OutPath, filepath.Base(forFile))
	}
	var extension = filepath.Ext(forFile)
	var base = forFile[0 : len(forFile)-len(extension)]
	return []string{base + ".obx.js"}
}

// Return the model filename for the given model JSON file.
func (gen *JSGenerator) ModelFile(forFile string, options generator.Options) string {
	if len(options.OutPath) > 0 {
		forFile = filepath.Join(options.OutPath, filepath.Base(forFile))
	}
	var extension = filepath.Ext(forFile)
	var fileStem = forFile[0 : len(forFile)-len(extension)]
	return fileStem + ".js"
}

func (JSGenerator) IsGeneratedFile(file string) bool {
	var name = filepath.Base(file)
	return name == "objectbox-model.js" || name == "schema.obx.js"
}

func (JSGenerator) IsSourceFile(file string) bool {
	return strings.HasSuffix(file, ".fbs")
}

func (gen *JSGenerator) ParseSource(sourceFile string) (*model.ModelInfo, error) {
	schemaReflection, err := flatbuffersc.ParseSchemaFile(sourceFile)
	if err != nil {
		return nil, err // already includes file name so no more context should be necessary
	}

	reader := fbSchemaReader{model: &model.ModelInfo{}, optional: gen.Optional}
	if err = reader.read(schemaReflection); err != nil {
		return nil, fmt.Errorf("error generating model from schema %s: %s", sourceFile, err)
	}

	return reader.model, nil
}

// Generate the schema.obx.js file, given the merged model info
func (gen *JSGenerator) WriteBindingFiles(sourceFile string, options generator.Options, mergedModel *model.ModelInfo) error {
	var err, err2 error

	var bindingFile = gen.BindingFiles(sourceFile, options)[0]

	// First generate the binding source
	var bindingSource []byte
	if bindingSource, err = gen.generateBindingFile(bindingFile, mergedModel); err != nil {
		return fmt.Errorf("can't generate binding file %s: %s", sourceFile, err)
	}

	if formattedSource, err := format(bindingSource); err != nil {
		// We just store error but still write the file so that we can check it manually
		err2 = fmt.Errorf("failed to format generated binding file %s: %s", bindingFile, err)
	} else {
		bindingSource = formattedSource
	}

	if err = generator.WriteFile(bindingFile, bindingSource, sourceFile); err != nil {
		return fmt.Errorf("can't write binding file %s: %s", sourceFile, err)
	} else if err2 != nil {
		// Now when the binding has been written (for debugging purposes), we can return the error
		return err2
	}

	return nil
}

func (gen *JSGenerator) generateBindingFile(bindingFile string, modelInfo *model.ModelInfo) (data []byte, err error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	var replaceSpecialChars = strings.NewReplacer("-", "_", ".", "_")
	var fileIdentifier = strings.ToLower(filepath.Base(bindingFile))
	fileIdentifier = replaceSpecialChars.Replace(fileIdentifier)

	// Arguments for the template
	type TplArgs struct {
		Model             *model.ModelInfo
		GeneratorVersion  int
		FileIdentifier    string
		Optional          string
		LangVersion       int
		EmptyStringAsNull bool
		NaNAsNull         bool
	}
	var tplArgs TplArgs
	tplArgs.Model = modelInfo
	tplArgs.GeneratorVersion = generator.VersionId
	tplArgs.FileIdentifier = fileIdentifier
	tplArgs.Optional = gen.Optional
	tplArgs.EmptyStringAsNull = gen.EmptyStringAsNull
	tplArgs.NaNAsNull = gen.NaNAsNull

	var tpl = templates.JsBindingTemplate

	if err = tpl.Execute(writer, tplArgs); err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	return b.Bytes(), nil
}

// Generate the objectbox-model.js, given the merged model info
func (gen *JSGenerator) WriteModelBindingFile(options generator.Options, mergedModel *model.ModelInfo) error {
	var err, err2 error

	var modelFile = gen.ModelFile(options.ModelInfoFile, options)
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

	if err = templates.JsModelTemplate.Execute(writer, tplArguments); err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	return b.Bytes(), nil
}

func removeEmptyLines(source []byte) []byte {
	// Split the source into lines
	lines := bytes.Split(source, []byte("\n"))

	// Filter out empty or whitespace-only lines
	var filteredLines [][]byte
	for _, line := range lines {
		if strings.TrimSpace(string(line)) != "" {
			filteredLines = append(filteredLines, line)
		}
	}

	// Join the filtered lines back together
	return bytes.Join(filteredLines, []byte("\n"))
}

func format(source []byte) ([]byte, error) {
	// NOTE we could do JS source formatting here

	// Replace tabs with spaces
	formatted := bytes.ReplaceAll(source, []byte("\t"), []byte("    "))
	//formatted = removeEmptyLines(formatted)

	return formatted, nil
}
