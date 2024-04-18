/*
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
	"text/template"
)

// ModelTemplate is used to generate the model initialization code
var ModelTemplate = template.Must(template.New("model").Parse(
	`// Code generated by ObjectBox; DO NOT EDIT.

package {{.Package}}

import (
	"github.com/objectbox/objectbox-go/objectbox"
)

// ObjectBoxModel declares and builds the model from all the entities in the package. 
// It is usually used when setting-up ObjectBox as an argument to the Builder.Model() function.
func ObjectBoxModel() *objectbox.Model {
	model := objectbox.NewModel()
	model.GeneratorVersion({{.GeneratorVersion}})

	{{range $entity := .Model.Entities -}}
	model.RegisterBinding({{$entity.Name}}Binding)
	{{end -}}
	model.LastEntityId({{.Model.LastEntityId.GetId}}, {{.Model.LastEntityId.GetUid}})
	{{if .Model.LastIndexId}}model.LastIndexId({{.Model.LastIndexId.GetId}}, {{.Model.LastIndexId.GetUid}}){{end}}
	{{if .Model.LastRelationId}}model.LastRelationId({{.Model.LastRelationId.GetId}}, {{.Model.LastRelationId.GetUid}}){{end}}

	return model
}`))
