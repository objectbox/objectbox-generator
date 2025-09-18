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

package gogenerator

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/objectbox/objectbox-generator/v4/internal/generator/binding"
	"github.com/objectbox/objectbox-generator/v4/internal/generator/model"
)

type uid = uint64
type id = uint32

var supportedEntityAnnotations = map[string]bool{
	"name":         false, // TODO
	"sync":         true,
	"transient":    true,
	"uid":          true,
	"external-name": true,
}

var supportedPropertyAnnotations = map[string]bool{
	"-":            true,
	"converter":    true,
	"date":         true,
	"date-nano":    true,
	"id":           true,
	"id-companion": true,
	"index":        true,
	"inline":       true,
	"lazy":         true,
	"link":         true,
	"name":         true,
	"type":         true,
	"uid":          true,
	"unique":       true,
	"external-name": true,
	"external-type": true,
}

// astReader contains information about the processed set of Entities
type astReader struct {
	Package *types.Package
	Imports map[string]string

	// model produced by reading the schema
	model *model.ModelInfo

	err    error
	source *file
}

// Entity holds the model information necessary to generate the binding code
type Entity struct {
	*binding.Object

	Fields []*Field // the tree of struct fields (necessary for embedded structs)

	binding *astReader // parent
}

// Merge implements model.EntityMeta interface
func (entity *Entity) Merge(mEntity *model.Entity) model.EntityMeta {
	entity.ModelEntity = mEntity
	return entity
}

// Property represents a mapping between a struct field and a DB field
type Property struct {
	*binding.Field

	IsBasicType bool
	GoType      string
	FbType      string
	Converter   *string

	// type casts for named types
	CastOnRead  string
	CastOnWrite string

	GoField *Field // actual code field this property represents
	Entity  *Entity

	annotations map[string]*binding.Annotation
}

// Merge implements model.PropertyMeta interface
func (property *Property) Merge(mProperty *model.Property) model.PropertyMeta {
	property.ModelProperty = mProperty
	return property
}

// Field is a field in an entity-struct. Not all fields become properties (e.g. to-many relations don't have a property)
type Field struct {
	Entity             *Entity // parent entity
	Name               string
	Type               string
	IsPointer          bool
	Property           *Property                 // nil if it's an embedded struct
	Fields             []*Field                  // inner fields, nil if it's a property
	StandaloneRelation *model.StandaloneRelation // to-many relation stored as a standalone relation in the model
	IsLazyLoaded       bool                      // only standalone (to-many) relations currently support lazy loading
	Meta               *Field                    // self reference for recursive ".Meta.Fields" access in the template

	path   string // relative addressing path for embedded structs
	parent *Field // when included in parent.Fields[], nil for top-level fields (directly in the entity)
}

func NewBinding() (*astReader, error) {
	return &astReader{model: &model.ModelInfo{}}, nil
}

func (r *astReader) CreateFromAst(f *file) (err error) {
	r.source = f
	r.Package = types.NewPackage(f.dir, f.pkgName)
	r.Imports = make(map[string]string)

	// this will hold the pointer to the latest GenDecl encountered (parent of the current struct)
	var prevDecl *ast.GenDecl

	// traverse the AST to process all structs
	f.walk(func(node ast.Node) bool {
		return r.entityLoader(node, &prevDecl)
	})

	if r.err != nil {
		return r.err
	}

	return nil
}

// this function only processes structs and cuts-off on types that can't contain a struct
func (r *astReader) entityLoader(node ast.Node, prevDecl **ast.GenDecl) bool {
	if r.err != nil {
		return false
	}

	switch v := node.(type) {
	case *ast.TypeSpec:
		if strct, isStruct := v.Type.(*ast.StructType); isStruct {
			var name = v.Name.Name

			if name == "" {
				// NOTE this should probably not happen
				r.err = fmt.Errorf("encountered a struct without a name")
				return false
			}

			var comments []*ast.Comment

			if v.Doc != nil && v.Doc.List != nil {
				// this will be defined in case the struct is inside a block of multiple types - `type (...)`
				comments = v.Doc.List

			} else if prevDecl != nil && *prevDecl != nil && (**prevDecl).Doc != nil && (**prevDecl).Doc.List != nil {
				// otherwise (`type A struct {`), use the docs from the parent GenDecl
				comments = (**prevDecl).Doc.List
			}

			r.err = r.createEntityFromAst(strct, name, comments)

			// no need to go any deeper in the AST
			return false
		}

		return true

	case *ast.GenDecl:
		// store the "parent" declaration - we need it to get the comments
		*prevDecl = v
		return true
	case *ast.File:
		return true
	}

	return false
}

func (r *astReader) createEntityFromAst(strct *ast.StructType, name string, comments []*ast.Comment) error {
	var modelEntity = model.CreateEntity(r.model, 0, 0)
	var entity = &Entity{Object: binding.CreateObject(modelEntity), binding: r}
	modelEntity.Meta = entity
	entity.SetName(name)

	if comments != nil {
		if err := entity.setAnnotations(comments); err != nil {
			return fmt.Errorf("%s on entity %s", err, entity.Name)
		}
	}

	{
		var fieldList = astStructFieldList{strct, r.source}
		var recursionStack = map[string]bool{}
		recursionStack[entity.Name] = true
		var err error
		entity.Fields, err = entity.addFields(nil, fieldList, entity.Name, "", &recursionStack)
		if err != nil {
			return err
		}
	}

	// TODO this is a new feature based on a transient/"-" annotation, previously not supported in Go
	// if entity.IsSkipped {
	// 	return nil
	// }

	if err := modelEntity.AutosetIdProperty([]model.PropertyType{model.PropertyTypeLong, model.PropertyTypeString}); err != nil {
		return fmt.Errorf("%s on entity %s", err, entity.Name)
	}

	// special handling for string IDs = they are transformed to uint64 in the binding
	if idProp, err := modelEntity.IdProperty(); err != nil {
		return fmt.Errorf("%s on entity %s", err, entity.Name)
	} else if idProp.Type == model.PropertyTypeString {
		var idPropMeta = idProp.Meta.(*Property)
		idProp.Type = model.PropertyTypeLong
		idPropMeta.FbType = "Uint64"
		idPropMeta.GoType = "uint64"

		if idPropMeta.annotations["converter"] == nil {
			var converter = "objectbox.StringIdConvert"
			idPropMeta.Converter = &converter
		}
	} else if !idProp.Meta.(*Property).hasValidTypeAsId() {
		return fmt.Errorf("id field '%s' has unsupported type '%s' on entity %s - must be one of [int64, uint64, string]",
			idProp.Meta.(*Property).Name, idProp.Meta.(*Property).GoType, entity.Name)
	} else {
		idProp.Meta.(*Property).FbType = "Uint64" // always stored as Uint64
	}

	r.model.Entities = append(r.model.Entities, modelEntity)

	return nil
}

func (entity *Entity) addFields(parent *Field, fields fieldList, fieldPath, prefix string, recursionStack *map[string]bool) ([]*Field, error) {
	var propertyLog = func(text string, property *Property) {
		log.Printf("%s property %s found in %s", text, property.Name, fieldPath)
	}
	var propertyError = func(err error, property *Property) error {
		return fmt.Errorf("%s on property %s found in %s", err, property.Name, fieldPath)
	}

	var children []*Field

	for i := 0; i < fields.Length(); i++ {
		f := fields.Field(i)

		var modelProperty = model.CreateProperty(entity.ModelEntity, 0, 0)
		var property = &Property{
			Field: binding.CreateField(modelProperty),

			Entity: entity, // TODO remove, there is Field.ModelProperty.Entity.Meta
		}
		modelProperty.Meta = property

		if name, err := f.Name(); err != nil {
			property.Name = strconv.FormatInt(int64(i), 10) // just for the error message
			return nil, propertyError(err, property)
		} else {
			property.SetName(name)
		}

		// this is used to correctly render embedded-structs initialization template
		var field = &Field{
			Entity:   entity,
			Name:     property.Name,
			Property: property,
			path:     fieldPath,
			parent:   parent,
		}
		field.Meta = field
		property.GoField = field

		if err := property.setAnnotations(f.Tag()); err != nil {
			return nil, propertyError(err, property)
		}

		if property.IsSkipped {
			continue
		}

		// if a field (type) is from a different package
		pkg, err := f.Package()
		if err != nil {
			return nil, propertyError(err, property)
		}
		var addImportPath = func() {} // called later when other info is available
		if pkg.Path() != entity.binding.Package.Path() {
			// check if it's available (starts with an uppercase letter)
			if len(field.Name) == 0 || field.Name[0] < 65 || field.Name[0] > 90 {
				propertyLog("Notice: skipping unavailable (private)", property)
				continue
			}

			// prepare the function for importing the package (if necessary, decided bellow)
			addImportPath = func() {
				if pkg.Name() == path.Base(pkg.Path()) {
					entity.binding.Imports[pkg.Path()] = pkg.Path()
				} else {
					entity.binding.Imports[pkg.Name()] = pkg.Path()
				}
			}
		}

		children = append(children, field)

		if property.annotations["type"] != nil {
			var annotatedType = property.annotations["type"].Value
			if len(annotatedType) > 1 && annotatedType[0] == '*' {
				field.IsPointer = true
				annotatedType = annotatedType[1:]
			}

			if err := property.setBasicType(annotatedType); err != nil {
				return nil, propertyError(err, property)
			}

		} else if innerStructFields, err := field.processType(f); err != nil {
			return nil, propertyError(err, property)

		} else if field.Type == "time.Time" {
			// first, try to handle time.Time struct - automatically set a converter if it's declared a date by the user
			if property.annotations["date"] == nil && property.annotations["date-nano"] == nil {
				property.annotations["date"] = &binding.Annotation{}
				propertyLog("Notice: time.Time is stored and read using millisecond precision in UTC by default on", property)
				log.Printf("To silence this notice either define your own converter using `converter` and " +
					"`type` annotations or add a `date` annotation explicitly")
			}

			// store the field as an int64
			// Note - property.annotations["type"] is not set or this code branch wouldn't be executed
			if err := property.setBasicType("int64"); err != nil {
				return nil, propertyError(err, property)
			}
			property.IsBasicType = false // override the value set by setBasicType

			if property.annotations["converter"] == nil {
				var converter = "objectbox.TimeInt64Convert"
				if property.annotations["date-nano"] != nil {
					converter = "objectbox.NanoTimeInt64Convert"
				}
				property.Converter = &converter
				property.annotations["type"] = &binding.Annotation{Value: "int64"}
			}

		} else if innerStructFields != nil {
			// if it was recognized as a struct that should be embedded, add all the fields

			var innerPrefix = prefix
			if property.annotations["inline"] == nil {
				// if NOT inline, use prefix based on the field name
				if len(innerPrefix) == 0 {
					innerPrefix = field.Name
				} else {
					innerPrefix = innerPrefix + "_" + field.Name
				}
			}

			// Structs may be chained in a cycle (using pointers), causing an infinite recursion.
			// Let's make sure this doesn't happen because it causes the generator (and a whole OS) to "freeze".
			if field.Type != "" {
				if (*recursionStack)[field.Type] {
					return nil, propertyError(fmt.Errorf("embedded struct cycle detected: %v", fieldPath), property)
				}
				(*recursionStack)[field.Type] = true
			}

			// apply some struct-related settings to the field
			field.Property = nil
			field.Fields, err = entity.addFields(field, innerStructFields, fieldPath+"."+property.Name, innerPrefix, recursionStack)
			if err != nil {
				return nil, err
			}

			addImportPath() // for structs, we're explicitly using the type so add the import

			delete(*recursionStack, field.Type)

			// this struct itself is not added, just the inner properties
			// so skip the the following steps of adding the property
			continue
		}

		// add import if necessary, i.e. we're explicitly using the type (but not for converters)
		if property.annotations["converter"] == nil && property.CastOnWrite != "" {
			addImportPath()
		}

		if property.annotations["converter"] != nil {
			if property.annotations["type"] == nil {
				return nil, propertyError(errors.New("type annotation has to be specified when using converters"), property)
			}
			property.Converter = &property.annotations["converter"].Value

			// converters use errors.New in the template
			entity.binding.Imports["errors"] = "errors"
		}

		if err := property.ProcessAnnotations(property.annotations); err != nil {
			return nil, propertyError(err, property)
		}

		if len(prefix) != 0 {
			property.ModelProperty.Name = prefix + "_" + property.ModelProperty.Name
			property.Name = prefix + "_" + property.Name
		}

		entity.ModelEntity.Properties = append(entity.ModelEntity.Properties, modelProperty)
	}

	return children, nil
}

// processType analyzes field type information and configures it.
// It might result in setting a field.Type (in case it's one of the basic types),
// field.StandaloneRelation (in case of many-to-many relations) or field.SimpleRelation (one-to-many relations).
// It also updates (fixes) the field.Name on embedded fields
func (field *Field) processType(f field) (fields fieldList, err error) {
	var typ = f.Type()
	var property = field.Property

	if err := property.setBasicType(typ.String()); err == nil {
		// if it's one of the basic supported types
		return nil, nil
	}

	// if not, get the underlying type and try again
	baseType, err := typ.UnderlyingOrError()
	if err != nil {
		return nil, err
	}

	// check if it needs a type cast (it is a named type, not an alias)
	var isNamed bool

	// in case it's a pointer, get it's underlying type
	if pointer, isPointer := baseType.(*types.Pointer); isPointer {
		baseType = pointer.Elem().Underlying()
		field.IsPointer = true
		isNamed = typesTypeErrorful{Type: baseType}.IsNamed()
	} else {
		isNamed = typ.IsNamed()
	}

	if err := property.setBasicType(baseType.String()); err == nil {
		// if the baseType is one of the basic supported types

		// check if it needs a type cast (it is a named type, not an alias)
		if isNamed {
			property.CastOnRead = baseType.String()
			property.CastOnWrite = path.Base(typ.String()) // sometimes, it may contain a full import path
		}

		return nil, nil
	}

	// try if it's a struct - it can be either embedded or a relation
	if strct, isStruct := baseType.(*types.Struct); isStruct {
		// fill in the field information
		field.fillInfo(f, typ)

		// if it's a one-to-many relation
		if property.annotations["link"] != nil {
			err := property.setRelationAnnotation(typeBaseName(typ.String()), false)
			property.IsBasicType = false // override the value set by setBasicType
			return nil, err
		}

		// otherwise inline all fields
		return structFieldList{strct}, nil
	}

	// check if it's a slice of a non-base type
	if slice, isSlice := baseType.(*types.Slice); isSlice {
		var elementType = slice.Elem()

		// it's a many-to-many relation
		if err := property.setRelationAnnotation(typeBaseName(elementType.String()), true); err != nil {
			return nil, err
		}

		var relDetails = make(map[string]*binding.Annotation)
		relDetails["name"] = &binding.Annotation{Value: field.Name}
		relDetails["to"] = property.annotations["link"]
		relDetails["uid"] = property.annotations["uid"]
		if property.annotations["external-name"] != nil {
			relDetails["external-name"] = property.annotations["external-name"]
		}
		if property.annotations["external-type"] != nil {
			relDetails["external-type"] = property.annotations["external-type"]
		}
		if rel, err := field.Entity.AddRelation(relDetails); err != nil {
			return nil, err
		} else {
			field.StandaloneRelation = rel
		}

		if field.Property.annotations["lazy"] != nil {
			// relations only
			field.IsLazyLoaded = true
		}

		// fill in the field information
		field.fillInfo(f, typesTypeErrorful{elementType})

		if _, isPointer := elementType.(*types.Pointer); isPointer {
			field.Type = "[]*" + field.Type
		} else {
			field.Type = "[]" + field.Type
		}

		// we need to skip adding this field (it's not persisted in DB) so we add an empty list of fields
		return structFieldList{}, nil
	}

	return nil, fmt.Errorf("unknown type %s", typ.String())
}

func (field *Field) fillInfo(f field, typ typeErrorful) {
	if namedType, isNamed := f.TypeInternal().(*types.Named); isNamed {
		field.Type = namedType.Obj().Name()
	} else {
		field.Type = typ.String()
	}

	// strip the '*' if it's a pointer type
	if len(field.Type) > 1 && field.Type[0] == '*' {
		field.Type = field.Type[1:]
	}

	// strip leading dots (happens sometimes, I think it's for local types from type-checked package)
	field.Type = strings.TrimLeft(field.Type, ".")

	// if the package path is specified (happens for embedded fields), check whether it's current package
	if strings.ContainsRune(strings.Replace(field.Type, "\\", "/", -1), '/') {
		// if the package is the current package, strip the path & name
		var parts = strings.Split(field.Type, ".")

		if len(parts) == 2 && parts[0] == field.Entity.binding.Package.Path() {
			field.Type = parts[len(parts)-1]
		}

	}
	// get just the last component from `packagename.typename` for the field name
	var parts = strings.Split(field.Name, ".")
	field.Name = parts[len(parts)-1]
}

func (entity *Entity) setAnnotations(comments []*ast.Comment) error {
	lines := parseCommentsLines(comments)
	var annotations = make(map[string]*binding.Annotation)

	for _, tags := range lines {
		// only handle comments in the form of:   // `tags`
		if len(tags) > 1 && tags[0] == tags[len(tags)-1] && tags[0] == '`' {
			if err := parseAnnotations(tags, &annotations, supportedEntityAnnotations); err != nil {
				return err
			}
		}
	}

	return entity.ProcessAnnotations(annotations)
}

func parseCommentsLines(comments []*ast.Comment) []string {
	var lines []string

	for _, comment := range comments {
		text := comment.Text
		text = strings.TrimSpace(text)

		// text is a single/multi line comment
		if strings.HasPrefix(text, "//") {
			text = strings.TrimPrefix(text, "//")
			lines = append(lines, strings.TrimSpace(text))

		} else if strings.HasPrefix(text, "/*") {
			text = strings.TrimPrefix(text, "/*")
			text = strings.TrimPrefix(text, "*")
			text = strings.TrimSuffix(text, "*/")
			text = strings.TrimSuffix(text, "*")
			text = strings.TrimSpace(text)
			for _, line := range strings.Split(text, "\n") {
				lines = append(lines, strings.TrimSpace(line))
			}
		} else {
			// unknown format, ignore
		}
	}

	return lines
}

func (property *Property) hasValidTypeAsId() bool {
	var goType = strings.ToLower(property.GoType)
	return goType == "int64" || goType == "uint64" || goType == "string"
}

func (property *Property) setAnnotations(tags string) error {
	var annotations = make(map[string]*binding.Annotation)
	if err := parseAnnotations(tags, &annotations, supportedPropertyAnnotations); err != nil {
		return err
	}

	if err := property.PreProcessAnnotations(annotations); err != nil {
		return err
	}

	property.annotations = annotations
	return nil
}

// setRelationAnnotation sets a relation on the property.
// If the user has previously defined a relation manually, it must match the arguments (relation target)
func (property *Property) setRelationAnnotation(target string, manyToMany bool) error {
	if property.annotations["link"] == nil {
		property.annotations["link"] = &binding.Annotation{}
	}

	if len(property.annotations["link"].Value) == 0 {
		// set the relation target to the type of the target entity
		// TODO this doesn't respect `objectbox:"name:entity"` on the entity (but we don't support that at the moment)
		property.annotations["link"].Value = target
	} else if property.annotations["link"].Value != target {
		return fmt.Errorf("relation target mismatch, expected %s, got %s", target, property.annotations["link"].Value)
	}

	if manyToMany {
		// nothing to do here, it's handled as a standalone relation so this "property" is skipped completely

	} else {
		// add this field as an ID field
		if err := property.setBasicType("uint64"); err != nil {
			return err
		}
	}

	return nil
}

func parseAnnotations(tags string, annotations *map[string]*binding.Annotation, supportedAnnotations map[string]bool) error {
	if len(tags) > 1 && tags[0] == tags[len(tags)-1] && (tags[0] == '`' || tags[0] == '"') {
		tags = tags[1 : len(tags)-1]
	}

	if tags == "" {
		return nil
	}

	// if it's a top-level call, i.e. tags is something like `objectbox:"tag1 tag2:value2" irrelevant:"value"`
	var tag = reflect.StructTag(tags)
	if contents, found := tag.Lookup("objectbox"); found {
		tags = contents
	} else if contents, found := tag.Lookup("ObjectBox"); found {
		tags = contents
	} else {
		return nil
	}

	return binding.ParseAnnotations(tags, annotations, supportedAnnotations)
}

func (property *Property) setBasicType(baseType string) error {
	property.GoType = baseType
	property.IsBasicType = true

	ts := property.GoType
	if property.GoType == "string" {
		property.ModelProperty.Type = model.PropertyTypeString
		property.FbType = "UOffsetT"
	} else if ts == "int" || ts == "int64" {
		property.ModelProperty.Type = model.PropertyTypeLong
		property.FbType = "Int64"
	} else if ts == "uint" || ts == "uint64" {
		property.ModelProperty.Type = model.PropertyTypeLong
		property.FbType = "Uint64"
		property.ModelProperty.AddFlag(model.PropertyFlagUnsigned)
	} else if ts == "int32" || ts == "rune" {
		property.ModelProperty.Type = model.PropertyTypeInt
		property.FbType = "Int32"
	} else if ts == "uint32" {
		property.ModelProperty.Type = model.PropertyTypeInt
		property.FbType = "Uint32"
		property.ModelProperty.AddFlag(model.PropertyFlagUnsigned)
	} else if ts == "int16" {
		property.ModelProperty.Type = model.PropertyTypeShort
		property.FbType = "Int16"
	} else if ts == "uint16" {
		property.ModelProperty.Type = model.PropertyTypeShort
		property.FbType = "Uint16"
		property.ModelProperty.AddFlag(model.PropertyFlagUnsigned)
	} else if ts == "int8" {
		property.ModelProperty.Type = model.PropertyTypeByte
		property.FbType = "Int8"
	} else if ts == "uint8" || ts == "byte" {
		property.ModelProperty.Type = model.PropertyTypeByte
		property.FbType = "Uint8"
		property.ModelProperty.AddFlag(model.PropertyFlagUnsigned)
	} else if ts == "[]byte" {
		property.ModelProperty.Type = model.PropertyTypeByteVector
		property.FbType = "UOffsetT"
	} else if ts == "[]float32" {
		property.ModelProperty.Type = model.PropertyTypeFloatVector
		property.FbType = "UOffsetT"
	} else if ts == "[]string" {
		property.ModelProperty.Type = model.PropertyTypeStringVector
		property.FbType = "UOffsetT"
	} else if ts == "float64" {
		property.ModelProperty.Type = model.PropertyTypeDouble
		property.FbType = "Float64"
	} else if ts == "float32" {
		property.ModelProperty.Type = model.PropertyTypeFloat
		property.FbType = "Float32"
	} else if ts == "bool" {
		property.ModelProperty.Type = model.PropertyTypeBool
		property.FbType = "Bool"
	} else {
		property.IsBasicType = false
		return fmt.Errorf("unknown type %s", ts)
	}

	return nil
}

// ObTypeString is called from the template
func (property *Property) ObTypeString() string {
	return model.PropertyTypeNames[property.ModelProperty.Type]
}

// HasNonIdProperty called from the template. The goal is to void GO error "variable declared and not used"
func (entity *Entity) HasNonIdProperty() bool {
	// since every entity MUST have an ID property, just check whether there's more than one property...
	return len(entity.ModelEntity.Properties) > 1
}

// HasRelations called from the template.
func (entity *Entity) HasRelations() bool {
	for _, field := range entity.Fields {
		if field.HasRelations() {
			return true
		}
	}

	return false
}

// HasLazyLoadedRelations called from the template.
func (entity *Entity) HasLazyLoadedRelations() bool {
	for _, field := range entity.Fields {
		if field.HasLazyLoadedRelations() {
			return true
		}
	}

	return false
}

// HasRelations called from the template.
func (field *Field) HasRelations() bool {
	if field.StandaloneRelation != nil || len(field.Property.ModelProperty.RelationTarget) > 0 {
		return true
	}

	for _, inner := range field.Fields {
		if inner.HasRelations() {
			return true
		}
	}

	return false
}

// HasLazyLoadedRelations called from the template.
func (field *Field) HasLazyLoadedRelations() bool {
	if field.StandaloneRelation != nil && field.IsLazyLoaded {
		return true
	}

	for _, inner := range field.Fields {
		if inner.HasLazyLoadedRelations() {
			return true
		}
	}

	return false
}

// Path returns full path to the field (in embedded struct)
// called from the template
func (field *Field) Path() string {
	var parts = strings.Split(field.path, ".")

	// strip the first component
	parts = parts[1:]

	parts = append(parts, field.Name)
	return strings.Join(parts, ".")
}

// HasPointersInPath checks whether there are any pointer-based fields in the path.
// Called from the template.
func (field *Field) HasPointersInPath() bool {
	if field.IsPointer {
		return true
	}

	if field.parent == nil {
		return false
	}

	return field.parent.HasPointersInPath()
}

// Path is called from the template. It returns full path to the property (in embedded struct).
func (property *Property) Path() string {
	return property.GoField.Path()
}

// AnnotatedType returns "type" annotation value
func (property *Property) AnnotatedType() string {
	return property.annotations["type"].Value
}

// TplReadValue returns a code to read the property value on a given object.
func (property *Property) TplReadValue(objVar, castType string) string {
	var valueAccessor = objVar

	if castType == "ptr-cast" {
		valueAccessor = valueAccessor + ".(*" + property.Entity.Name + ")"
	} else if castType == "val-cast" {
		valueAccessor = valueAccessor + ".(" + property.Entity.Name + ")"
	}

	valueAccessor = valueAccessor + "." + property.Path()

	if property.Converter != nil {
		return *property.Converter + "ToDatabaseValue(" + valueAccessor + ")" // returns value & error
	}

	// While not explicitly, this is currently only true if called from GetId() template part.
	// NOTE: currently we don't handle this for converters - they should work on uint64 for IDs
	if property.ModelProperty.IsIdProperty() && property.GoType != "uint64" {
		valueAccessor = "uint64(" + valueAccessor + ")"
	}

	return valueAccessor + ", nil" // return value & err=nil
}

// TplSetAndReturn returns a code to write the property value on a given object.
func (property *Property) TplSetAndReturn(objVar, castType, rhs string) string {
	var lhs = objVar

	if castType == "ptr-cast" {
		lhs = lhs + ".(*" + property.Entity.Name + ")"
	} else if castType == "val-cast" {
		lhs = lhs + ".(" + property.Entity.Name + ")"
	}

	lhs = lhs + "." + property.Path()

	var ret = "nil"

	if property.Converter != nil {
		lhs = `var err error
` + lhs + `, err`
		rhs = *property.Converter + "ToEntityProperty(" + rhs + ")"
		ret = "err"
	}

	// While not explicitly, this is currently only true if called from SetId() template part.
	if property.ModelProperty.IsIdProperty() && property.GoType != "uint64" {
		// TODO this won't compile because converters (i.e. `rhs`) now return two values
		rhs = property.GoType + "(" + rhs + ")"
	}

	return lhs + " = " + rhs + `
return ` + ret
}

func typeBaseName(name string) string {
	// strip the '*' if it's a pointer type
	name = strings.TrimPrefix(name, "*")

	// get just the last component from `packagename.typename` for the field name
	if strings.ContainsRune(name, '.') {
		name = strings.TrimPrefix(path.Ext(name), ".")
	}

	return name
}
