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
	"sort"
	"strings"
	"text/template"

	"github.com/objectbox/objectbox-generator/v4/internal/generator/model"
)

// cccToUc converts CapitalCamelCase to UPPER_CASE - only used converty Go PropertyFlags names to C/Core names.
// Note: this isn't library quality, e.g. only handles ascii letters.
func cccToUc(str string) string {
	var result string
	for _, char := range str {
		// if it's an uppercase character and not the first one, prepend an underscore ("space")
		if char >= 65 && char <= 90 && len(result) > 0 {
			result += "_"
		}
		result += strings.ToUpper(string(char))
	}
	return result
}

var funcMap = template.FuncMap{
	"PropTypeName": func(val model.PropertyType) string {
		return model.PropertyTypeNames[val]
	},
	"CorePropFlags": func(val model.PropertyFlags) string {
		var result []string

		// Get sorted flag names to avoid changes in the generated code. Go map iteration order is not guaranteed.
		for flag, name := range model.PropertyFlagNames {
			if val&flag != 0 { // if this flag is set
				result = append(result, "OBXPropertyFlags_"+cccToUc(name))
			}
		}

		if len(result) > 1 {
			sort.Strings(result)
			// if there's more than one, we need to cast the result of their combination back to the right type
			return strings.Join(result, " | ")
		} else if len(result) > 0 {
			return result[0]
		}
		return ""
	},
	"CoreEntityFlags": func(val model.EntityFlags) string {
		var result []string

		// Get sorted flag names to avoid changes in the generated code. Go map iteration order is not guaranteed.
		for flag, name := range model.EntityFlagNames {
			if val&flag != 0 { // if this flag is set
				result = append(result, "OBXEntityFlags_"+cccToUc(name))
			}
		}

		if len(result) > 1 {
			sort.Strings(result)
			// if there's more than one, we need to cast the result of their combination back to the right type
			return strings.Join(result, " | ")
		} else if len(result) > 0 {
			return result[0]
		}
		return ""
	},
	"CoreHnswFlags": func(val model.HnswFlags) string {
		var result []string

		// Get sorted flag names to avoid changes in the generated code. Go map iteration order is not guaranteed.
		for flag, name := range model.HnswFlagNames {
			if val&flag != 0 { // if this flag is set
				result = append(result, "OBXHnswFlags_"+name)
			}
		}

		if len(result) > 1 {
			sort.Strings(result)
			// if there's more than one, we need to cast the result of their combination back to the right type
			return "(" + strings.Join(result, " | ") + ")"
		} else if len(result) > 0 {
			return result[0]
		}
		return "OBXHnswFlags_NONE"
	},
	"CoreExternalTypes": func(val model.ExternalType) string {
		return "OBXExternalPropertyType_" + model.ExternalTypeNames[val]
	},
	"PrintComments": func(tabs int, comments []string) string {
		var result string
		for _, comment := range comments {
			result += "/// " + comment + "\n" + strings.Repeat("\t", tabs)
		}
		return result
	},
	"IsOptionalPtr": func(optional string) bool {
		return optional == "std::unique_ptr" || optional == "std::shared_ptr"
	},
	"ToUpper": strings.ToUpper,
}
