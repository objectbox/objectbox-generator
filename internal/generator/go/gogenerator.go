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
	binding *Binding
}

// BindingFile returns a name of the binding file for the given entity file.
func BindingFile(sourceFile string) string {
	var extension = filepath.Ext(sourceFile)
	return sourceFile[0:len(sourceFile)-len(extension)] + ".obx" + extension
}

// ModelFile returns the model GO file for the given JSON info file path
func ModelFile(modelInfoFile string) string {
	var extension = filepath.Ext(modelInfoFile)
	return modelInfoFile[0:len(modelInfoFile)-len(extension)] + ".go"
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
		return nil, fmt.Errorf("can't init Binding: %s", err)
	}

	if err = goGen.binding.CreateFromAst(f); err != nil {
		return nil, fmt.Errorf("can't prepare bindings for %s: %s", sourceFile, err)
	}

	// TODO convert binding to model
	return nil, nil
}

func (goGen *GoGenerator) WriteBindingFiles(sourceFile string, options generator.Options, mergedModel *model.ModelInfo) error {
	var err, err2 error

	var bindingSource []byte
	if bindingSource, err = goGen.generateBindingFile(options); err != nil {
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

func (goGen *GoGenerator) generateBindingFile(options generator.Options) (data []byte, err error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	var tplArguments = struct {
		Binding          *Binding
		GeneratorVersion int
		Options          generator.Options
	}{goGen.binding, generator.Version, options}

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

	if modelSource, err = generateModelFile(modelInfo); err != nil {
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
