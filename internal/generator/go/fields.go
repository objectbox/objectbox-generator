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

package gogenerator

import (
	"fmt"
	"go/ast"
	"go/types"
	"path"
	"strings"
)

// these interfaces are used in the binding to iterate over fields coming from multiple sources (AST & type checker)
type fieldList interface {
	Length() int
	Field(int) field
}

type field interface {
	Name() (string, error)
	Tag() string
	Type() typeErrorful
	TypeInternal() types.Type
	Package() (*types.Package, error)
}

type typeErrorful interface {
	String() string
	UnderlyingOrError() (types.Type, error)

	// whether it's an alias of a basic type or rather a named type
	IsNamed() bool
}

//region ast.StructType wrappers

type astStructFieldList struct {
	*ast.StructType
	source *file
}

func (fields astStructFieldList) Length() int {
	return len(fields.Fields.List)
}

func (fields astStructFieldList) Field(i int) field {
	return &astStructField{fields.Fields.List[i], fields.source}
}

type astStructField struct {
	*ast.Field
	source *file
}

func (field astStructField) Name() (string, error) {
	if len(field.Names) == 0 {
		// in case of an unnamed field, use the type name
		var typ = types.ExprString(field.Field.Type)

		// strip the '*' if it's a pointer type
		typ = strings.TrimPrefix(typ, "*")

		// remove the package from the name
		if strings.ContainsRune(typ, '.') {
			typ = strings.TrimPrefix(path.Ext(typ), ".")
		}

		return typ, nil
	} else if len(field.Names) == 1 {
		return field.Names[0].Name, nil
	} else {
		return "", fmt.Errorf("the field has too many names: %v", len(field.Names))
	}
}

func (field astStructField) Tag() string {
	if field.Field.Tag != nil {
		return field.Field.Tag.Value
	}
	return ""
}

func (field astStructField) Type() typeErrorful {
	return astTypeExpr{Expr: field.Field.Type, source: field.source}
}

func (field astStructField) TypeInternal() types.Type {
	return astTypeExpr{Expr: field.Field.Type, source: field.source}
}

func (field astStructField) Package() (*types.Package, error) {
	// handle fields referring to an imported type
	if selector, ok := field.Field.Type.(*ast.SelectorExpr); ok {
		var importAlias = selector.X.(*ast.Ident).Name
		pkg, err := field.source.importedPackage(importAlias)
		if err != nil {
			return nil, err
		}
		return pkg, nil
	}

	// by default, use current package
	return types.NewPackage(field.source.dir, field.source.pkgName), nil
}

type astTypeExpr struct {
	ast.Expr
	source *file
}

func (expr astTypeExpr) String() string {
	return types.ExprString(expr.Expr)
}

func (expr astTypeExpr) IsNamed() bool {
	if t, err := expr.source.getType(expr.Expr); err != nil {
		panic(err)
	} else {
		return typesTypeErrorful{Type: t}.IsNamed()
	}
}

func (expr astTypeExpr) Underlying() types.Type {
	if t, err := expr.UnderlyingOrError(); err != nil {
		panic(err)
	} else {
		return t
	}
}

func (expr astTypeExpr) UnderlyingOrError() (types.Type, error) {
	t, err := expr.source.getType(expr.Expr)
	if err != nil {
		return nil, err
	}
	return t.Underlying(), nil
}

//endregion

//region types.Struct wrappers

type structFieldList struct {
	*types.Struct
}

func (fields structFieldList) Length() int {
	if fields.Struct == nil {
		return 0
	}

	return fields.Struct.NumFields()
}

func (fields structFieldList) Field(i int) field {
	return structField{fields.Struct.Field(i), fields.Tag(i)}
}

type structField struct {
	*types.Var
	tag string
}

func (field structField) Name() (string, error) {
	return field.Var.Name(), nil
}

func (field structField) Tag() string {
	return field.tag
}

func (field structField) Type() typeErrorful {
	return typesTypeErrorful{field.Var.Type()}
}

func (field structField) TypeInternal() types.Type {
	return field.Var.Type()
}

func (field structField) Package() (*types.Package, error) {
	var pkg = field.Var.Pkg()
	if pkg == nil {
		return nil, fmt.Errorf("could not get package for %v", *field.Var)
	}
	return pkg, nil
}

type typesTypeErrorful struct {
	types.Type
}

func (typ typesTypeErrorful) String() string {
	return typ.Type.String()
}

func (typ typesTypeErrorful) IsNamed() bool {
	_, isNamed := typ.Type.(*types.Named)
	return isNamed
}

func (typ typesTypeErrorful) UnderlyingOrError() (types.Type, error) {
	return typ.Type.Underlying(), nil
}

//endregion
