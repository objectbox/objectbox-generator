/*
 * ObjectBox Generator - a build time tool for ObjectBox
 * Copyright (C) 2020-2024 ObjectBox Ltd. All rights reserved.
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

package rename

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/objectbox/objectbox-generator/v4/internal/generator"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/model"
	"github.com/objectbox/objectbox-generator/v4/test/assert"
	"github.com/objectbox/objectbox-generator/v4/test/integration"
)

func getUid(t *testing.T, id model.IdUid) string {
	uid, err := id.GetUid()
	assert.NoErr(t, err)
	return strconv.FormatInt(int64(uid), 10)
}

func TestCpp(t *testing.T) {
	dbDir, err := ioutil.TempDir("", "generator-test-db")
	assert.NoErr(t, err)
	defer os.RemoveAll(dbDir)
	var envVars = []string{"dbDir=" + dbDir}

	conf := &integration.CCppTestConf{}
	defer conf.Cleanup()

	// BEFORE RENAME start
	conf.CreateCMake(t, integration.CppDefault, "step-1.cpp")
	conf.Generate(t, map[string]string{"schema.fbs": `table OldEntityName {
	id:uint64;
	oldPropertyName:int;
}`})
	modelJSONFile := generator.ModelInfoFile(conf.Cmake.ConfDir)
	modelInfo, err := model.LoadModelFromJSONFile(modelJSONFile)
	assert.NoErr(t, err)
	assert.Eq(t, 1, len(modelInfo.Entities))
	assert.Eq(t, 2, len(modelInfo.Entities[0].Properties))
	assert.Eq(t, "oldPropertyName", modelInfo.Entities[0].Properties[1].Name)
	entityUid := getUid(t, modelInfo.Entities[0].Id)
	propertyUid := getUid(t, modelInfo.Entities[0].Properties[1].Id)
	conf.Build(t)
	conf.Run(t, envVars)
	// BEFORE RENAME end

	// AFTER RENAME start
	conf.CreateCMake(t, integration.CppDefault, "step-2.cpp")
	conf.Generate(t, map[string]string{"schema.fbs": "/// objectbox: uid=" + entityUid + `
table NewEntityName {
	id:uint64;
` + "/// objectbox: uid=" + propertyUid + `
	newPropertyName:int;
}`})
	modelInfo, err = model.LoadModelFromJSONFile(modelJSONFile)
	assert.NoErr(t, err)
	assert.Eq(t, 1, len(modelInfo.Entities))
	assert.Eq(t, "NewEntityName", modelInfo.Entities[0].Name)
	assert.Eq(t, 2, len(modelInfo.Entities[0].Properties))
	assert.Eq(t, "newPropertyName", modelInfo.Entities[0].Properties[1].Name)
	assert.Eq(t, entityUid, getUid(t, modelInfo.Entities[0].Id))
	assert.Eq(t, propertyUid, getUid(t, modelInfo.Entities[0].Properties[1].Id))
	conf.Build(t)
	conf.Run(t, envVars)
	// AFTER RENAME end
}
