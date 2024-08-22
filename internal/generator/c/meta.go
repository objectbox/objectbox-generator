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

package cgenerator

import (
	"fmt"
	"sort"
	"strings"

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

// CppName returns C++ symbol/variable name with reserved keywords suffixed by an underscore
func (mo *fbsObject) CppName() string {
	return cppName(mo.Name)
}

// CName returns CppName() prefixed by a namespace (with underscores)
func (mo *fbsObject) CName() string {
	var prefix string
	if len(mo.Namespace) != 0 {
		prefix = strings.Replace(mo.Namespace, ".", "_", -1) + "_"
	}

	return prefix + mo.CppName()
}

func cppNamespacePrefix(ns string) string {
	if len(ns) == 0 {
		return ""
	}
	return strings.Join(strings.Split(ns, "."), "::") + "::"
}

// CppNamespacePrefix returns c++ namespace prefix for symbol definition
func (mo *fbsObject) CppNamespacePrefix() string {
	return cppNamespacePrefix(mo.Namespace)
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

// PreDeclareCppRelTargets returns C++ struct pre-declarations for related entities.
func (mo *fbsObject) PreDeclareCppRelTargets() (string, error) {
	// first create a map `(ns.entity) => bool`, then sort it to keep the code from changing, and generate the C++ decl.
	var m = make(map[string]bool)

	for _, rel := range mo.ModelEntity.Relations {
		m[rel.Target.Meta.(*fbsObject).Namespace+"."+rel.Target.Name] = true
	}
	for _, prop := range mo.ModelEntity.Properties {
		if len(prop.RelationTarget) > 0 {
			m[prop.Meta.(*fbsField).relTargetNamespace()+"."+prop.RelationTarget] = true
		}
	}

	// sort
	var sortedUniqueTargets []string
	for k := range m {
		sortedUniqueTargets = append(sortedUniqueTargets, k)
	}
	sort.Strings(sortedUniqueTargets)

	// generated the code
	var code string
	for _, name := range sortedUniqueTargets {
		var nss []string // namespaces
		if strings.HasPrefix(name, ".") {
			name = name[1:] // no NS
		} else {
			nss = strings.Split(name, ".")
			name = nss[len(nss)-1]
			nss = nss[0 : len(nss)-1]
		}

		var line string
		for _, ns := range nss {
			line = line + "namespace " + ns + " { "
		}
		line = line + "struct " + cppName(name) + "; "
		line = line + strings.Repeat("}", len(nss))
		code = code + line + "\n"
	}
	return code, nil
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

// CppName returns C++ variable name with reserved keywords suffixed by an underscore
func (mp *fbsField) CppName() string {
	return cppName(mp.Name)
}

// CppNameRelationTarget returns C++ target class name with reserved keywords suffixed by an underscore
func (mp *fbsField) CppNameRelationTarget() string {
	return cppNamespacePrefix(mp.relTargetNamespace()) + cppName(mp.ModelProperty.RelationTarget)
}

// CppType returns C++ type name
func (mp *fbsField) CppType() string {
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
func (mp *fbsField) CppFbType() string {
	var cppType = mp.CppType()
	if cppType == "bool" {
		cppType = "uint8_t"
	}
	return cppType
}

// CppTypeWithOptional returns full C++ type name, including wrapper if the value is not defined
func (mp *fbsField) CppTypeWithOptional() (string, error) {
	var cppType = mp.CppType()
	if len(mp.Optional) != 0 {
		if mp.ModelProperty.IsIdProperty() {
			return "", fmt.Errorf("ID property must not be optional: %s.%s", mp.ModelProperty.Entity.Name, mp.ModelProperty.Name)
		}
		cppType = mp.Optional + "<" + cppType + ">"
	}
	return cppType, nil
}

// CppValOp returns field value access operator
func (mp *fbsField) CppValOp() string {
	if len(mp.Optional) != 0 {
		return "->"
	}
	return "."
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

// FbOffsetFactory returns an offset factory used to build flatbuffers if this property is a complex type.
// See also FbOffsetType().
func (mp *fbsField) FbOffsetFactory() string {
	switch mp.ModelProperty.Type {
	case model.PropertyTypeString:
		return "CreateString"
	case model.PropertyTypeByteVector:
		return "CreateVector"
	case model.PropertyTypeFloatVector:
		return "CreateVector"
	case model.PropertyTypeStringVector:
		return "CreateVectorOfStrings"
	}
	return ""
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

// CppName returns C++ variable name with reserved keywords suffixed by an underscore
func (mr *standaloneRel) CppName() string {
	return cppName(mr.ModelRelation.Name)
}
