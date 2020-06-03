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

package cgenerator

import (
	"strings"

	"github.com/objectbox/objectbox-go/internal/generator/binding"
	"github.com/objectbox/objectbox-go/internal/generator/fbsparser/reflection"
	"github.com/objectbox/objectbox-go/internal/generator/model"
)

type fbsObject struct {
	*binding.Object
	mEntity   *model.Entity
	fbsObject *reflection.Object
}

// Merge implements model.PropertyMeta interface
func (mo *fbsObject) Merge(entity *model.Entity) model.EntityMeta {
	return &fbsObject{mo.Object, entity, mo.fbsObject}
}

// CppName returns C++ variable name with reserved keywords suffixed by an underscore
func (mo *fbsObject) CppName() string {
	return cppName(mo.Name)
}

// CppNamespaceStart returns c++ namespace opening declaration
func (mo *fbsObject) CppNamespaceStart() string {
	if len(mo.Namespace) == 0 {
		return ""
	}

	var nss = strings.Split(mo.Namespace, ".")
	for i, ns := range nss {
		nss[i] = "namespace " + ns + " {"
	}
	return strings.Join(nss, "\n")
}

// CppNamespaceEnd returns c++ namespace closing declaration
func (mo *fbsObject) CppNamespaceEnd() string {
	if len(mo.Namespace) == 0 {
		return ""
	}
	var result = ""
	var nss = strings.Split(mo.Namespace, ".")
	for _, ns := range nss {
		// print in reversed order
		result = "}  // namespace " + ns + "\n" + result
	}
	return result
}

type fbsField struct {
	*binding.Field
	mProp    *model.Property
	fbsField *reflection.Field
}

// Merge implements model.PropertyMeta interface
func (mp *fbsField) Merge(property *model.Property) model.PropertyMeta {
	return &fbsField{mp.Field, property, mp.fbsField}
}

// CppName returns C++ variable name with reserved keywords suffixed by an underscore
func (mp *fbsField) CppName() string {
	return cppName(mp.Name)
}

// CppType returns C++ type name
func (mp *fbsField) CppType() string {
	var fbsType = mp.fbsField.Type(nil)
	var baseType = fbsType.BaseType()
	var cppType = fbsTypeToCppType[baseType]
	if baseType == reflection.BaseTypeVector {
		cppType = cppType + "<" + fbsTypeToCppType[fbsType.Element()] + ">"
	}
	return cppType
}

// FbIsVector returns true if the property is considered a vector type.
func (mp *fbsField) FbIsVector() bool {
	switch mp.mProp.Type {
	case model.PropertyTypeString:
		return true
	case model.PropertyTypeByteVector:
		return true
	case model.PropertyTypeStringVector:
		return true
	}
	return false
}

// CElementType returns C vector element type name
func (mp *fbsField) CElementType() string {
	switch mp.mProp.Type {
	case model.PropertyTypeByteVector:
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

// FbOffsetFactory returns an offset factory used to build flatbuffers if this property is a complex type.
// See also FbOffsetType().
func (mp *fbsField) FbOffsetFactory() string {
	switch mp.mProp.Type {
	case model.PropertyTypeString:
		return "CreateString"
	case model.PropertyTypeByteVector:
		return "CreateVector"
	case model.PropertyTypeStringVector:
		return "CreateVectorOfStrings"
	}
	return ""
}

// FbOffsetType returns a type used to read flatbuffers if this property is a complex type.
// See also FbOffsetFactory().
func (mp *fbsField) FbOffsetType() string {
	switch mp.mProp.Type {
	case model.PropertyTypeString:
		return "flatbuffers::Vector<char>"
	case model.PropertyTypeByteVector:
		return "flatbuffers::Vector<" + fbsTypeToCppType[mp.fbsField.Type(nil).Element()] + ">"
	case model.PropertyTypeStringVector:
		return "" // NOTE custom handling in the template
	}
	return ""
}

// FbDefaultValue returns a default value for scalars
func (mp *fbsField) FbDefaultValue() string {
	switch mp.mProp.Type {
	case model.PropertyTypeFloat:
		return "0.0f"
	case model.PropertyTypeDouble:
		return "0.0"
	}
	return "0"
}
