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
	"fmt"
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
				result = append(result, "OBXPropertyFlags."+cccToUc(name))
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
				result = append(result, "OBXEntityFlags."+cccToUc(name))
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
				result = append(result, "OBXHnswFlags."+name)
			}
		}

		if len(result) > 1 {
			sort.Strings(result)
			// if there's more than one, we need to cast the result of their combination back to the right type
			return "(" + strings.Join(result, " | ") + ")"
		} else if len(result) > 0 {
			return result[0]
		}
		return "OBXHnswFlags.NONE"
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

	"AddField": func(property model.Property) string {
		varName := "object." + property.Name
		notSupportedComment := fmt.Sprint("// Not supported: ", model.PropertyTypeNames[property.Type])
		isNullable := (property.Flags & model.PropertyFlagNotNull) == 0

		var str string

		if isNullable {
			str += "if (" + varName + " != null) {\n"
		}

		switch property.Type {
		case model.PropertyTypeBool:
			varVal := fmt.Sprint(varName, " ? 1 : 0")
			str += fmt.Sprintln("fbb.addFieldInt8(", property.FbSlot(), ", ", varVal, ");")
		case model.PropertyTypeByte:
			str += fmt.Sprintln("fbb.addFieldInt8(", property.FbSlot(), ", ", varName, ");")
		case model.PropertyTypeShort:
			str += fmt.Sprintln("fbb.addFieldInt16(", property.FbSlot(), ", ", varName, ");")
		case model.PropertyTypeChar:
			str += fmt.Sprintln("fbb.addFieldInt16(", property.FbSlot(), ", ", varName, ");")
		case model.PropertyTypeInt:
			str += fmt.Sprintln("fbb.addFieldInt32(", property.FbSlot(), ", ", varName, ");")
		case model.PropertyTypeLong:
			str += fmt.Sprintln("fbb.addFieldInt64(", property.FbSlot(), ", ", varName, ");")
		case model.PropertyTypeFloat:
			str += fmt.Sprintln("fbb.addFieldFloat32(", property.FbSlot(), ", ", varName, ");")
		case model.PropertyTypeDouble:
			str += fmt.Sprintln("fbb.addFieldFloat64(", property.FbSlot(), ", ", varName, ");")
		case model.PropertyTypeString:
			return "" // Not an inline field
		case model.PropertyTypeDate:
			str += fmt.Sprintln("fbb.addFieldInt64(", property.FbSlot(), ", ", varName, ");")
		case model.PropertyTypeRelation:
			return "" // Not an inline field
		case model.PropertyTypeDateNano:
			str += notSupportedComment
		case model.PropertyTypeByteVector:
			return "" // Not an inline field
		case model.PropertyTypeFloatVector:
			return "" // Not an inline field
		case model.PropertyTypeStringVector:
			return "" // Not an inline field
		default:
			panic("Unknown property type")
		}

		if isNullable {
			str += "}"
		}
		return str
	},

	"CreateOffsetProperty": func(property model.Property) string {
		offsetVar := property.Name + "_offset"
		fieldVar := "object." + property.Name

		switch property.Type {
		case model.PropertyTypeString:
			return fmt.Sprint("const ", offsetVar, " = fbb.createString(", fieldVar, ");")
		case model.PropertyTypeByteVector:
			return fmt.Sprint("const ", offsetVar, " = fbb.createByteVector(Uint8Array.from(", fieldVar, "));")
		case model.PropertyTypeFloatVector:
			return fmt.Sprint("const ", offsetVar, " = fbb.createByteVector(new Uint8Array(Float32Array.from(", fieldVar, ")));")
		case model.PropertyTypeStringVector:
			return "" // TODO: string vectors not supported right now
		default:
			return ""
		}
	},

	"AddFieldOffset": func(property model.Property) string {
		offsetVarName := property.Name + "_offset"
		return fmt.Sprint("fbb.addFieldOffset(", property.FbSlot(), ",", offsetVarName, ");")
	},

	"WriteGetAssignOffset": func(property model.Property) string {
		value, err := property.FbvTableOffset()
		if err != nil {
			panic(err)
		}
		offsetVarName := property.Name + "_offset"
		return fmt.Sprint("const ", offsetVarName, " = bb.__offset(bbPos, ", value, ");")
	},

	"ReadProperty": func(property model.Property) string {
		offsetVarName := property.Name + "_offset"
		assignLhs := "outObject." + property.Name + " = "
		switch property.Type {
		case model.PropertyTypeBool:
			return fmt.Sprint(assignLhs, "bb.readInt8(bbPos + ", offsetVarName, ") ? true : false;")
		case model.PropertyTypeByte:
			return fmt.Sprint(assignLhs, "bb.readInt8(bbPos + ", offsetVarName, ");")
		case model.PropertyTypeShort:
			return fmt.Sprint(assignLhs, "bb.readInt16(bbPos + ", offsetVarName, ");")
		case model.PropertyTypeChar:
			return fmt.Sprint(assignLhs, "bb.readInt16(bbPos + ", offsetVarName, ");")
		case model.PropertyTypeInt:
			return fmt.Sprint(assignLhs, "bb.readInt32(bbPos + ", offsetVarName, ");")
		case model.PropertyTypeLong:
			return fmt.Sprint(assignLhs, "bb.readInt64(bbPos + ", offsetVarName, ");")
		case model.PropertyTypeFloat:
			return fmt.Sprint(assignLhs, "bb.readFloat32(bbPos + ", offsetVarName, ");")
		case model.PropertyTypeDouble:
			return fmt.Sprint(assignLhs, "bb.readFloat64(bbPos + ", offsetVarName, ");")
		case model.PropertyTypeString:
			return fmt.Sprint(assignLhs, "bb.__string(bbPos + ", offsetVarName, ");")
		case model.PropertyTypeDate:
			return fmt.Sprint(assignLhs, "bb.readInt64(bbPos + ", offsetVarName, ");")
		case model.PropertyTypeRelation:
			return fmt.Sprint("// ", assignLhs, "PropertyTypeRelation") // TODO
		case model.PropertyTypeDateNano:
			return fmt.Sprint("// ", assignLhs, "PropertyTypeDateNano") // TODO
		case model.PropertyTypeByteVector:
			return fmt.Sprint("// ", assignLhs, "PropertyTypeByteVector") // TODO
		case model.PropertyTypeFloatVector:
			return fmt.Sprint("// ", assignLhs, "PropertyTypeFloatVector") // TODO
		case model.PropertyTypeStringVector:
			return fmt.Sprint("// ", assignLhs, "PropertyTypeStringVector") // TODO
		default:
			return ""
		}
	},

	"OBXTypeToJSPropertyType": func(propertyType model.PropertyType) string {
		switch propertyType {
		case model.PropertyTypeBool:
			return "BoolProperty"
		case model.PropertyTypeByte:
			return "ByteProperty"
		case model.PropertyTypeShort:
			return "ShortProperty"
		case model.PropertyTypeInt:
			return "IntProperty"
		case model.PropertyTypeLong:
			return "LongProperty"
		case model.PropertyTypeFloat:
			return "FloatProperty"
		case model.PropertyTypeDouble:
			return "DoubleProperty"
		case model.PropertyTypeString:
			return "StringProperty"
		case model.PropertyTypeDate:
			return "DateProperty"
		case model.PropertyTypeFloatVector:
			return "Float32VectorProperty"
		// case model.PropertyTypeRelation:
		// 	return "number" // or Relation type?
		// case model.PropertyTypeDateNano:
		// 	return "bigint" // or Date?
		// case model.PropertyTypeByteVector, model.PropertyTypeFloatVector, model.PropertyTypeStringVector:
		// 	return "Array" // or specific type?
		default:
			panic("Unknown property type")
		}
	},

	"IsIdPropertyFlagPresent": func(flags model.PropertyFlags) bool {
		return flags&model.PropertyFlagId != 0
	},
}
