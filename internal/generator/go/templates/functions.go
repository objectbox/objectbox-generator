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

package templates

import (
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"StringTitle": strings.Title,
	"StringCamel": func(s string) string {
		result := strings.Title(s)
		return strings.ToLower(result[0:1]) + result[1:]
	},
	"TypeIdentifier": func(s string) string {
		if strings.HasPrefix(s, "[]") {
			return strings.Title(s[2:]) + "Vector"
		}
		return strings.Title(s)
	},
}
