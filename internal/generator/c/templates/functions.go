/*
 * Copyright 2019 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package templates

import (
	"strings"
	"text/template"

	"github.com/objectbox/objectbox-generator/internal/generator/model"
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
		var result string
		var count = 0

		for flag, name := range model.PropertyFlagNames {
			if val&flag != 0 { // if this flag is set
				if count > 0 {
					result += " | "
				}
				result += "OBXPropertyFlags_" + cccToUc(name)
				count++
			}
		}

		if count > 1 {
			// if there's more than one, we need to cast the result of their combination back to the right type
			result = "OBXPropertyFlags(" + result + ")"
		}

		return result
	},
}
