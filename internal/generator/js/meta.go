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

package jsgenerator

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/objectbox/objectbox-generator/v4/internal/generator/binding"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/flatbuffersc/reflection"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/model"
)

type fbsObject struct {
	*binding.Object
	fbsObject *reflection.Object
}

// Merge implements model.EntityMeta interface
func (mo *fbsObject) Merge(entity *model.Entity) model.EntityMeta {
	mo.ModelEntity = entity
	return mo
}

// JsName returns JS symbol/variable name with reserved keywords suffixed by an underscore
func (mo *fbsObject) JsName() string {
	return mo.Name
}

type fbsField struct {
	*binding.Field
	fbsField *reflection.Field
}

// Merge implements model.PropertyMeta interface
func (mp *fbsField) Merge(property *model.Property) model.PropertyMeta {
	mp.ModelProperty = property
	return mp
}

// JsName returns C++ variable name with reserved keywords suffixed by an underscore
func (mp *fbsField) JsName() string {
	return mp.Name
}

// JsType returns C++ type name
func (mp *fbsField) JsType() string {
	var fbsType = mp.fbsField.Type(nil)
	var baseType = fbsType.BaseType()
	var cppType = fbsTypeToCppType[baseType]
	if baseType == reflection.BaseTypeVector {
		cppType = cppType + "<" + fbsTypeToCppType[fbsType.Element()] + ">"
	} else if (mp.ModelProperty.IsIdProperty() || mp.ModelProperty.Type == model.PropertyTypeRelation) && cppType == "uint64_t" {
		cppType = "obx_id" // defined in objectbox.h
	}
	return cppType
}

// CppFbType returns C++ type name used in flatbuffers templated functions
func (mp *fbsField) JsFbType() string {
	var cppType = mp.JsType()
	if cppType == "bool" {
		cppType = "uint8_t"
	}
	return cppType
}

// CppTypeWithOptional returns full C++ type name, including wrapper if the value is not defined
func (mp *fbsField) JsTypeWithOptional() (string, error) {
	var cppType = mp.JsType()
	if len(mp.Optional) != 0 {
		if mp.ModelProperty.IsIdProperty() {
			return "", fmt.Errorf("ID property must not be optional: %s.%s", mp.ModelProperty.Entity.Name, mp.ModelProperty.Name)
		}
		cppType = mp.Optional + "<" + cppType + ">"
	}
	return cppType, nil
}

// FbIsVector returns true if the property is considered a vector type.
func (mp *fbsField) FbIsVector() bool {
	switch mp.ModelProperty.Type {
	case model.PropertyTypeString:
		return true
	case model.PropertyTypeByteVector:
		return true
	case model.PropertyTypeFloatVector:
		return true
	case model.PropertyTypeStringVector:
		return true
	}
	return false
}

// CElementType returns C vector element type name
func (mp *fbsField) CElementType() string {
	switch mp.ModelProperty.Type {
	case model.PropertyTypeByteVector:
		return fbsTypeToCppType[mp.fbsField.Type(nil).Element()]
	case model.PropertyTypeFloatVector:
		return fbsTypeToCppType[mp.fbsField.Type(nil).Element()]
	case model.PropertyTypeString:
		return "char"
	case model.PropertyTypeStringVector:
		return "char*"
	}
	return ""
}

// FlatccFnPrefix returns the field's type as used in Flatcc.
func (mp *fbsField) FlatccFnPrefix() string {
	return fbsTypeToFlatccFnPrefix[mp.fbsField.Type(nil).BaseType()]
}

// FbTypeSize returns the field's type flatbuffers size.
func (mp *fbsField) FbTypeSize() uint8 {
	return fbsTypeSize[mp.fbsField.Type(nil).BaseType()]
}

// Reference:
// https://stackoverflow.com/a/40811635

type P map[string]interface{}

func renderStr(format string, params P) string {
	buffer := &bytes.Buffer{}
	template.Must(template.New("").Parse(format)).Execute(buffer, params)
	return buffer.String()
}

// FbOffsetType returns a type used to read flatbuffers if this property is a complex type.
// See also FbOffsetFactory().
func (mp *fbsField) FbOffsetType() string {
	switch mp.ModelProperty.Type {
	case model.PropertyTypeString:
		return "flatbuffers::Vector<char>"
	case model.PropertyTypeByteVector:
		return "flatbuffers::Vector<" + fbsTypeToCppType[mp.fbsField.Type(nil).Element()] + ">"
	case model.PropertyTypeFloatVector:
		return "flatbuffers::Vector<" + fbsTypeToCppType[mp.fbsField.Type(nil).Element()] + ">"
	case model.PropertyTypeStringVector:
		return "" // NOTE custom handling in the template
	}
	return ""
}

// FbDefaultValue returns a default value for scalars
func (mp *fbsField) FbDefaultValue() string {
	switch mp.ModelProperty.Type {
	case model.PropertyTypeFloat:
		return "0.0f"
	case model.PropertyTypeDouble:
		return "0.0"
	}
	return "0"
}

// FbIsFloatingPoint returns true if type is float or double
func (mp *fbsField) FbIsFloatingPoint() bool {
	switch mp.ModelProperty.Type {
	case model.PropertyTypeFloat:
		return true
	case model.PropertyTypeDouble:
		return true
	}
	return false
}

// Try to determine the namespace of the target entity but don't fail if we can't because it's declared in a different
// file. Assume no namespace in that case and hope for the best.
func (mp *fbsField) relTargetNamespace() string {
	if targetEntity, err := mp.ModelProperty.Entity.Model.FindEntityByName(mp.ModelProperty.RelationTarget); err == nil {
		if targetEntity.Meta != nil {
			return targetEntity.Meta.(*fbsObject).Namespace
		}
	}
	return ""
}

type standaloneRel struct {
	ModelRelation *model.StandaloneRelation
}

// Merge implements model.PropertyMeta interface
func (mr *standaloneRel) Merge(rel *model.StandaloneRelation) model.StandaloneRelationMeta {
	mr.ModelRelation = rel
	return mr
}

// JsName returns JS variable name with reserved keywords suffixed by an underscore
func (mr *standaloneRel) JsName() string {
	return mr.ModelRelation.Name
}
