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

package model

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// LoadOrCreateModel reads a model file or creates a new one if it doesn't exist
func LoadOrCreateModel(path string) (model *ModelInfo, err error) {
	if fileExists(path) {
		return LoadModelFromJSONFile(path)
	}
	return createModelJSONFile(path)
}

// Close and unlock model
func (model *ModelInfo) Close() error {
	return model.file.Close()
}

// Write current model data to file
func (model *ModelInfo) Write() error {
	data, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		return err
	}

	if err = model.file.Truncate(0); err != nil {
		return err
	}

	if _, err := model.file.WriteAt(data, 0); err != nil {
		return err
	}

	if err = model.file.Sync(); err != nil {
		return err
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func LoadModelFromJSONFile(path string) (model *ModelInfo, err error) {
	model = &ModelInfo{}

	if model.file, err = os.OpenFile(path, os.O_RDWR, 0); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(io.Reader(model.file))

	if err == nil {
		err = json.Unmarshal(data, model)
	}

	if err != nil {
		defer model.Close()
		return nil, fmt.Errorf("can't read file %s: %s", path, err)
	}

	// until objectbox-go 0.9 we didn't have model version in the file but it was basically version 4; recognize this
	if model.ModelVersion == 0 && model.MinimumParserVersion == 0 && len(model.Note1) == 0 {
		model.ModelVersion = 4
		model.MinimumParserVersion = 4
	}

	model.fillMissing()

	return model, nil
}

func createModelJSONFile(path string) (model *ModelInfo, err error) {
	model = createModelInfo()

	// create a file handle so to have an exclusive access
	if model.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600); err != nil {
		return nil, err
	}

	// write it with initial content (so that we know it's writable & it would have correct contents on next tool run)
	if err = model.Write(); err != nil {
		defer model.Close()
		return nil, fmt.Errorf("can't write file %s: %s", path, err)
	}

	return model, nil
}
