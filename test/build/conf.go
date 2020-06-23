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

package build

import (
	"path"
)

const flatbuffersDir = "third_party/flatbuffers-c-bridge/third_party/flatbuffers"
const flatccDir = "third_party/flatcc"
const ObjectBoxCDir = "third_party/objectbox-c"

func IncludeDirs(repoRoot string) []string {
	var result []string
	result = append(result, path.Join(repoRoot, flatbuffersDir, "include"))
	result = append(result, path.Join(repoRoot, flatccDir, "include"))
	result = append(result, path.Join(repoRoot, ObjectBoxCDir, "include"))
	result = append(result, path.Join(repoRoot, "third_party"))
	return result
}

func LibDirs(repoRoot string) []string {
	var result []string
	result = append(result, path.Join(repoRoot, ObjectBoxCDir, "lib"))
	result = append(result, path.Join(repoRoot, flatccDir, "lib"))
	return result
}
