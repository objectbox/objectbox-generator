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

package cgenerator

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/objectbox/objectbox-go/internal/generator"
	"github.com/objectbox/objectbox-go/internal/generator/c/templates"
	"github.com/objectbox/objectbox-go/internal/generator/fbsparser"
	"github.com/objectbox/objectbox-go/internal/generator/model"
)

type CGenerator struct {
}

// BindingFile returns a name of the binding file for the given entity file.
func BindingFile(sourceFile string) string {
	var extension = filepath.Ext(sourceFile)
	return sourceFile[0:len(sourceFile)-len(extension)] + ".obx.h"
}

// ModelFile returns the model GO file for the given JSON info file path
func ModelFile(modelInfoFile string) string {
	var extension = filepath.Ext(modelInfoFile)
	return modelInfoFile[0:len(modelInfoFile)-len(extension)] + ".h"
}

func (CGenerator) IsGeneratedFile(file string) bool {
	var name = filepath.Base(file)
	return name == "objectbox-model.h" || name == "objectbox-model.c" ||
		strings.HasSuffix(name, ".obx.c") || strings.HasSuffix(name, ".obx.h")
}

func (gen *CGenerator) ParseSource(sourceFile string) (*model.ModelInfo, error) {
	schemaReflection, err := fbsparser.ParseSchemaFile(sourceFile)
	if err != nil {
		return nil, err // already includes file name so no more context should be necessary
	}

	reader := fbSchemaReader{model: &model.ModelInfo{}}
	if err = reader.read(schemaReflection); err != nil {
		return nil, fmt.Errorf("error generating model from schema %s: %s", sourceFile, err)
	}

	return reader.model, nil
}

func (gen *CGenerator) WriteBindingFiles(sourceFile string, options generator.Options, mergedModel *model.ModelInfo) error {
	// TODO
	return nil
}

func (gen *CGenerator) WriteModelBindingFile(options generator.Options, mergedModel *model.ModelInfo) error {
	var err, err2 error

	var modelFile = ModelFile(options.ModelInfoFile)
	var modelSource []byte

	if modelSource, err = generateModelFile(mergedModel); err != nil {
		return fmt.Errorf("can't generate model file %s: %s", modelFile, err)
	}

	// TODO c source formatting
	// if formattedSource, err := format.Source(modelSource); err != nil {
	// 	// we just store error but still writ the file so that we can check it manually
	// 	err2 = fmt.Errorf("failed to format generated model file %s: %s", modelFile, err)
	// } else {
	// 	modelSource = formattedSource
	// }

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
	}{m, generator.Version}

	if err = templates.ModelTemplate.Execute(writer, tplArguments); err != nil {
		return nil, fmt.Errorf("template execution failed: %s", err)
	}

	if err = writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer: %s", err)
	}

	return b.Bytes(), nil
}
