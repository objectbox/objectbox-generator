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

package comparison

import (
	"flag"
	"strings"
	"testing"

	"github.com/objectbox/objectbox-generator/v4/test/assert"
)

// used during development of generator to overwrite the "golden" files
var overwriteExpected = flag.Bool("update", false,
	"Update all '.expected' files with the generated content. "+
		"It's up to the developer to actually check before committing whether the newly generated files are correct")

// used during development of generator to test a single directory
var target = flag.String("target", "", "Specify target subdirectory to generate")

func TestCompare(t *testing.T) {
	if *target == "" {
		for key, _ := range confs {
			generateAllDirs(t, *overwriteExpected, key)
		}
	} else if parts := strings.Split(*target, "/"); len(parts) == 1 {
		generateAllDirs(t, *overwriteExpected, parts[0])
	} else if len(parts) == 2 {
		srcType, genType := typesFromConfKey(parts[0])
		conf, ok := confs[parts[0]]
		assert.True(t, ok)
		conf.helper.init(t, conf)
		generateOneDir(t, *overwriteExpected, conf, srcType, genType, parts[1])
	} else {
		t.Fatal("invalid target specification, expected 1 or two parts separated by '/'")
	}
}
