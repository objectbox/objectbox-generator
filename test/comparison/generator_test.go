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
	"flag"
	"strings"
	"testing"
)

// used during development of generator to overwrite the "golden" files
var overwriteExpected = flag.Bool("update", false,
	"Update all '.expected' files with the generated content. "+
		"It's up to the developer to actually check before committing whether the newly generated files are correct")

// used during development of generator to test a single directory
var target = flag.String("target", "", "Specify target subdirectory to generate")

func TestCompare(t *testing.T) {
	if *target == "" {
		for langDir, _ := range confs {
			generateAllDirs(t, *overwriteExpected, langDir)
		}
	} else if parts := strings.Split(*target, "/"); len(parts) == 1 {
		generateAllDirs(t, *overwriteExpected, parts[0])
	} else {
		generateOneDir(t, *overwriteExpected, confs[parts[0]], "testdata/"+*target)
	}
}
