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
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/objectbox/objectbox-generator/v4/internal/generator"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/model"
	"github.com/objectbox/objectbox-generator/v4/test/assert"
	"github.com/objectbox/objectbox-generator/v4/test/integration"
)

func TestCpp(t *testing.T) {
	dbDir, err := ioutil.TempDir("", "generator-test-db")
	assert.NoErr(t, err)
	defer os.RemoveAll(dbDir)
	var envVars = []string{"dbDir=" + dbDir}

	conf := &integration.CCppTestConf{}
	defer conf.Cleanup()

	// BEFORE start
	conf.CreateCMake(t, integration.CppDefault, "step-1.cpp")
	conf.Generate(t, map[string]string{"schema.fbs": `table EntityName {
	id:uint64;
	value:int;
}`})
	modelJSONFile := generator.ModelInfoFile(conf.Cmake.ConfDir)
	modelInfo, err := model.LoadModelFromJSONFile(modelJSONFile)
	assert.NoErr(t, err)
	assert.Eq(t, 1, len(modelInfo.Entities))
	assert.Eq(t, 2, len(modelInfo.Entities[0].Properties))
	{
		id, err := modelInfo.Entities[0].Properties[1].Id.GetId()
		assert.NoErr(t, err)
		assert.Eq(t, 2, int(id))
	}
	conf.Build(t)
	conf.Run(t, envVars)
	// BEFORE end

	// AFTER start
	modelInfo.Rand = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	newUid, err := modelInfo.GenerateUid()
	t.Logf("Changing property '%s' %s UID to %d",
		modelInfo.Entities[0].Properties[1].Name, modelInfo.Entities[0].Properties[1].Id, newUid)
	assert.NoErr(t, err)
	conf.CreateCMake(t, integration.CppDefault, "step-2.cpp")
	conf.Generate(t, map[string]string{"schema.fbs": `table EntityName {
	id:uint64;
` + "/// objectbox: uid=" + strconv.FormatInt(int64(newUid), 10) + `
	value:int;
}`})
	modelInfo, err = model.LoadModelFromJSONFile(modelJSONFile)
	assert.NoErr(t, err)
	assert.Eq(t, 1, len(modelInfo.Entities))
	assert.Eq(t, 2, len(modelInfo.Entities[0].Properties))
	{
		id, uid, err := modelInfo.Entities[0].Properties[1].Id.Get()
		assert.NoErr(t, err)
		assert.Eq(t, 3, int(id))
		assert.Eq(t, newUid, uid)
	}
	conf.Build(t)
	conf.Run(t, envVars)
	// AFTER end
}
