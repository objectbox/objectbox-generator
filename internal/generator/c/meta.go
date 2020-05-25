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
	"github.com/objectbox/objectbox-go/internal/generator/fbsparser/reflection"
	"github.com/objectbox/objectbox-go/internal/generator/model"
)

type fbsProperty struct {
	mProp    *model.Property
	fbsField *reflection.Field
}

// Merge implements model.PropertyMeta interface
func (mp *fbsProperty) Merge(property *model.Property) model.PropertyMeta {
	return &fbsProperty{property, mp.fbsField}
}

func (mp *fbsProperty) CppName() string {
	if reservedKeywords[mp.mProp.Name] {
		return mp.mProp.Name + "_"
	}
	return mp.mProp.Name
}

func (mp *fbsProperty) CppType() string {
	var fbsType = mp.fbsField.Type(nil)
	var baseType = fbsType.BaseType()
	var cppType = fbsTypeToCppType[baseType]
	if baseType == reflection.BaseTypeVector {
		cppType = cppType + "<" + fbsTypeToCppType[fbsType.Element()] + ">"
	}
	return cppType
}
