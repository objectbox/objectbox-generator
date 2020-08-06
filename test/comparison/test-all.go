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

package comparison

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-generator/internal/generator"
	"github.com/objectbox/objectbox-generator/test/assert"
)

func typesFromConfKey(confKey string) (srcType, genType string) {
	types := strings.Split(confKey, "-")
	srcType = types[0]
	genType = types[len(types)-1]
	return
}

// generateAllDirs walks through the "data" and generates bindings for each subdirectory of langDir
// set overwriteExpected to TRUE to update all ".expected" files with the generated content
func generateAllDirs(t *testing.T, overwriteExpected bool, confKey string) {
	t.Logf("Testing %s code generator", confKey)

	srcType, genType := typesFromConfKey(confKey)
	testCases, err := ioutil.ReadDir(path.Join("testdata", srcType))
	assert.NoErr(t, err)

	conf, ok := confs[confKey]
	assert.True(t, ok)
	conf.helper.init(t, conf)

	for _, testCase := range testCases {
		if !testCase.IsDir() {
			continue
		}

		var tc = testCase.Name() // need to create a variable in order to be captured properly by the lambda below
		t.Run(confKey+"/"+testCase.Name(), func(t *testing.T) {
			t.Parallel()
			generateOneDir(t, overwriteExpected, conf, srcType, genType, tc)
		})
	}
}

func generateOneDir(t *testing.T, overwriteExpected bool, conf testSpec, srcType, genType, testCase string) {
	var srcDir = filepath.Join("testdata", srcType, testCase) // where input files, e.g. schema.fbs, are
	var expDir = srcDir                                       // where expected files for the current type are
	var genDir = srcDir                                       // where should the generator output the files
	if srcType != genType {
		expDir = filepath.Join(srcDir, genType)
		genDir = filepath.Join(srcDir, genType)
	}

	var errorTransformer = func(err error) error {
		return err
	}

	var cleanup = func() {}
	defer func() {
		cleanup()
	}()

	// Test in a temporary directory - if tested by an end user, the repo is read-only.
	tempRoot, err := ioutil.TempDir("", "objectbox-generator-test")
	assert.NoErr(t, err)

	// we can't defer directly because compilation step is run in a separate goroutine after this function exits
	cleanup = func() {
		assert.NoErr(t, os.RemoveAll(tempRoot))
	}

	genDir = filepath.Join(tempRoot, testCase)
	t.Logf("Testing in a temporary directory %s", genDir)
	assert.NoErr(t, os.MkdirAll(genDir, 0700))

	if conf.helper != nil {
		if errTrans := conf.helper.prepareTempDir(t, conf, srcDir, genDir, tempRoot); errTrans != nil {
			errorTransformer = errTrans
		}
	}

	// Go generator updates generator go.mod when loading files (adds the missing objectbox-go import).
	// Therefore, we'll load files from the temp dir instead
	srcDir = genDir

	modelInfoFile := generator.ModelInfoFile(genDir)
	modelInfoExpectedFile := generator.ModelInfoFile(srcDir) + ".expected"

	modelCodeFile := conf.generator.ModelFile(modelInfoFile)
	modelCodeExpectedFile := conf.generator.ModelFile(generator.ModelInfoFile(expDir)) + ".expected"

	// run the generation twice, first time with deleting old modelInfo
	for i := 0; i <= 1; i++ {
		if i == 0 {
			t.Logf("Testing %s->%s '%s' without model info JSON", srcType, genType, testCase)
			os.Remove(modelInfoFile)
		} else if testing.Short() {
			continue // don't test twice in "short" tests
		} else {
			t.Logf("Testing %s->%s '%s' with previous model info JSON", srcType, genType, testCase)
		}

		// setup the desired directory contents by copying "*.initial" files to their name without the extension
		setupInitialFiles(t, srcDir, genDir)
		setupInitialFiles(t, expDir, genDir)

		if 0 != generateAllFiles(t, overwriteExpected, conf, srcDir, expDir, genDir, modelInfoFile, errorTransformer) {
			assertSameFile(t, modelInfoFile, modelInfoExpectedFile, overwriteExpected)
			assertSameFile(t, modelCodeFile, modelCodeExpectedFile, overwriteExpected)
		}
	}

	// verify the result can be built
	if !testing.Short() && conf.helper != nil {
		// override the defer to prevent cleanup before compilation is actually run
		var cleanupAfterCompile = cleanup
		cleanup = func() {}

		t.Run("compile", func(t *testing.T) {
			defer cleanupAfterCompile()
			t.Parallel()
			var expectedError error
			if fileExists(path.Join(expDir, "compile-error.expected")) {
				content, err := ioutil.ReadFile(path.Join(expDir, "compile-error.expected"))
				assert.NoErr(t, err)
				expectedError = errors.New(string(content))
			}
			conf.helper.build(t, conf, genDir, expectedError, errorTransformer)
		})
	}
}

func setupInitialFiles(t *testing.T, srcDir, targetDir string) {
	srcFiles, err := filepath.Glob(filepath.Join(srcDir, "*.initial"))
	assert.NoErr(t, err)
	for _, srcFile := range srcFiles {
		targetFile := filepath.Base(srcFile)
		targetFile = targetFile[0 : len(targetFile)-len(".initial")]
		targetFile = filepath.Join(targetDir, targetFile)
		assert.NoErr(t, CopyFile(srcFile, targetFile, 0))
	}
}

func assertSameFile(t *testing.T, file string, expectedFile string, overwriteExpected bool) {
	if overwriteExpected && fileExists(file) {
		assert.NoErr(t, CopyFile(file, expectedFile, 0))
	}

	// if no file is expected
	if !fileExists(expectedFile) {
		// there can be no source file either
		if fileExists(file) {
			assert.Failf(t, "%s is missing but %s exists", expectedFile, file)
		}
		return
	}

	content, err := ioutil.ReadFile(file)
	assert.NoErr(t, err)

	contentExpected, err := ioutil.ReadFile(expectedFile)
	assert.NoErr(t, err)

	if 0 != bytes.Compare(content, contentExpected) {
		assert.Failf(t, "generated file %s is not the same as %s", file, expectedFile)
	}
}

func generateAllFiles(t *testing.T, overwriteExpected bool, conf testSpec, srcDir, expDir, genDir string, modelInfoFile string, errorTransformer func(error) error) int {
	// remove generated files during development (they might be syntactically wrong)
	if overwriteExpected {
		files, err := filepath.Glob(filepath.Join(genDir, "*"+conf.generatedExt))
		assert.NoErr(t, err)

		for _, file := range files {
			assert.NoErr(t, os.Remove(file))
		}
	}

	var positiveTestsCount = 0

	// process all source files in the directory
	inputFiles, err := filepath.Glob(filepath.Join(srcDir, "*"+conf.sourceExt))
	assert.NoErr(t, err)
	for _, sourceFile := range inputFiles {
		// skip generated files & "expected results" files
		if conf.generator.IsGeneratedFile(sourceFile) ||
			strings.HasSuffix(sourceFile, ".skip"+conf.sourceExt) ||
			strings.HasSuffix(sourceFile, "expected") ||
			strings.HasSuffix(sourceFile, "initial") {
			continue
		}

		t.Logf("  %s", filepath.Base(sourceFile))

		var options = generator.Options{
			ModelInfoFile: modelInfoFile,
			// NOTE zero seed for test-only - avoid changes caused by random numbers by fixing them to the same seed
			Rand:          rand.New(rand.NewSource(0)),
			CodeGenerator: conf.helper.generatorFor(t, conf, sourceFile, genDir),
		}
		err = errorTransformer(generator.Process(sourceFile, options))

		// handle negative test
		var shouldFail = strings.HasSuffix(filepath.Base(sourceFile), ".fail"+conf.sourceExt)
		if shouldFail {
			if err == nil {
				assert.Failf(t, "Unexpected PASS on a negative test %s", sourceFile)
			} else {
				var errPlatformIndependent = strings.Replace(err.Error(), "\\", "/", -1)
				assert.Eq(t, getExpectedError(t, sourceFile).Error(), errPlatformIndependent)
				continue
			}
		} else {
			positiveTestsCount++
		}

		assert.NoErr(t, err)

		var bindingFiles = options.CodeGenerator.BindingFiles(sourceFile)
		for _, bindingFile := range bindingFiles {
			var expectedFile = strings.Replace(bindingFile, genDir, expDir, 1) + ".expected"
			assertSameFile(t, bindingFile, expectedFile, overwriteExpected)
		}
	}
	return positiveTestsCount
}

var expectedErrorRegexp = regexp.MustCompile(`// *ERROR *=(.+)[\n|\r]`)
var expectedErrorRegexpMulti = regexp.MustCompile(`(?sU)/\* *ERROR.*[\n|\r](.+)\*/`)

func getExpectedError(t *testing.T, sourceFile string) error {
	source, err := ioutil.ReadFile(sourceFile)
	assert.NoErr(t, err)

	if match := expectedErrorRegexp.FindSubmatch(source); len(match) > 1 {
		return errors.New(strings.TrimSpace(string(match[1]))) // this is a "positive" return
	}

	if match := expectedErrorRegexpMulti.FindSubmatch(source); len(match) > 1 {
		return errors.New(strings.TrimSpace(string(match[1]))) // this is a "positive" return
	}

	assert.Failf(t, "missing error declaration in %s - add comment to the file // ERROR = expected error text", sourceFile)
	return nil
}
