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

package cgenerator

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/objectbox/objectbox-generator/v4/internal/generator"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/c/templates"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/flatbuffersc"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/model"
)

type CGenerator struct {
	PlainC            bool
	LangVersion       int    // -1: unset, cpp: 11, 14, 17
	Optional          string // std::optional, std::unique_ptr, std::shared_ptr
	EmptyStringAsNull bool
	NaNAsNull         bool
}

// BindingFiles returns the names of the generated C or C++ language binding files for the given entity file.
func (gen *CGenerator) BindingFiles(forFile string, options generator.Options) []string {

	if len(options.OutPath) > 0 {
		forFile = filepath.Join(options.OutPath, filepath.Base(forFile))
	}
	var extension = filepath.Ext(forFile)
	var base = forFile[0 : len(forFile)-len(extension)]

	if gen.PlainC {
		return []string{base + ".obx.h"}
	}
	var headerBase = base
	if len(options.OutHeadersPath) > 0 {
		headerBase = filepath.Join(options.OutHeadersPath, filepath.Base(forFile))
		headerBase = headerBase[0 : len(headerBase)-len(extension)]
	}

	return []string{headerBase + ".obx.hpp", base + ".obx.cpp"}
}

// ModelFile returns the generated model C header file for the given JSON info file path
func (gen *CGenerator) ModelFile(forFile string, options generator.Options) string {

	if len(options.OutHeadersPath) > 0 {
		forFile = filepath.Join(options.OutHeadersPath, filepath.Base(forFile))
	} else if len(options.OutPath) > 0 {
		forFile = filepath.Join(options.OutPath, filepath.Base(forFile))
	}
	var extension = filepath.Ext(forFile)
	return forFile[0:len(forFile)-len(extension)] + ".h"
}

func (CGenerator) IsGeneratedFile(file string) bool {
	var name = filepath.Base(file)
	return name == "objectbox-model.h" ||
		strings.HasSuffix(name, ".obx.h") ||
		strings.HasSuffix(name, ".obx.hpp") ||
		strings.HasSuffix(name, ".obx.cpp")
}

func (CGenerator) IsSourceFile(file string) bool {
	return strings.HasSuffix(file, ".fbs")
}

func (gen *CGenerator) ParseSource(sourceFile string) (*model.ModelInfo, error) {
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

func (gen *CGenerator) WriteBindingFiles(sourceFile string, options generator.Options, mergedModel *model.ModelInfo) error {
	var err, err2 error

	var bindingFiles = gen.BindingFiles(sourceFile, options)

	for _, bindingFile := range bindingFiles {
		var bindingSource []byte
		if bindingSource, err = gen.generateBindingFile(bindingFile, bindingFiles[0], mergedModel); err != nil {
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
	}

	return nil
}

func (gen *CGenerator) generateBindingFile(bindingFile, headerFile string, m *model.ModelInfo) (data []byte, err error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	var replaceSpecialChars = strings.NewReplacer("-", "_", ".", "_")
	var fileIdentifier = strings.ToLower(filepath.Base(bindingFile))
	fileIdentifier = replaceSpecialChars.Replace(fileIdentifier)

	var tplArguments = struct {
		Model             *model.ModelInfo
		GeneratorVersion  int
		FileIdentifier    string
		HeaderFile        string
		Optional          string
		LangVersion       int
		EmptyStringAsNull bool
		NaNAsNull         bool
	}{m, generator.VersionId, fileIdentifier, filepath.Base(headerFile), gen.Optional, gen.LangVersion, gen.EmptyStringAsNull, gen.NaNAsNull}

	var tpl *template.Template

	if gen.PlainC {
		tpl = templates.CBindingTemplate
	} else if bindingFile == headerFile {
		tpl = templates.CppBindingTemplateHeader
	} else {
		tpl = templates.CppBindingTemplate
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

	if err = templates.ModelTemplate.Execute(writer, tplArguments); err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	return b.Bytes(), nil
}

func format(source []byte) ([]byte, error) {
	// NOTE we could do C/C++ source formatting here if there was an easy to integrate go module.
	// For now, we just try to do our best within the templates themselves.

	// replace tabs with spaces
	return bytes.ReplaceAll(source, []byte("\t"), []byte("    ")), nil
}
