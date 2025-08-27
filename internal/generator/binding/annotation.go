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

package binding

import (
	"fmt"
	"strings"
)

// Annotation is a tag on a struct-field
type Annotation struct {
	Value string
	// Details is used to map complex annotations, e.g. many-to-many with brackets syntax, sync annotation.
	// e.g. relation(name=manyToManyRelName,to=TargetEntity)
	Details map[string]*Annotation
}

// HasDetail checks if the annotation has a "Detail" with the given name.
func (a Annotation) HasDetail(name string) bool {
	return a.Details != nil && a.Details[name] != nil
}

// HasBooleanDetail checks an annotation in the given map if it has a "Detail" with an empty value.
func HasBooleanDetail(annotations map[string]*Annotation, parentAnnotation, name string) (bool, error) {
	if annotations[parentAnnotation] != nil {
		a := annotations[parentAnnotation]
		if a.HasDetail(name) {
			if len(a.Details[name].Value) != 0 {
				return true, fmt.Errorf("'%s' annotation's '%s' attribute value must be empty", parentAnnotation, name)
			}
			return true, nil
		}
	}
	return false, nil
}

// ParseAnnotations parses annotations in any of the following formats.
// name="name",index - creates two annotations, name and index, the former having a non-empty value
// relation(name=manyToManyRelName,to=TargetEntity) - creates a single annotation relation with two items as details
// id - creates a single annotation
// NOTE: this started as a very simple parser but it seems like the requirements are ever-increasing... maybe some form
//
//	of recursive tokenization would be better in case we decided to rework.
func ParseAnnotations(str string, annotations *map[string]*Annotation, supportedAnnotations map[string]bool) error {
	var s annotationInProgress
	for i := 0; i < len(str); i++ {
		var char = str[i]

		if !s.valueQuoted && (char == '=' || char == ':' || char == '(') { // start a value
			if len(s.name) == 0 {
				return fmt.Errorf("invalid annotation format: name expected before '%s' at position %d in `%s` ", string(char), i, str)
			}
			s.value = &Annotation{}

			// special handling for "recursive" details (many-to-many relation)
			if char == '(' {
				// find the closing bracket
				var detailsStr string
				for j := i + 1; j < len(str); j++ {
					if str[j] == ')' { // NOTE we're ignoring potential closing brackets in quotes
						detailsStr += str[i+1 : j]
						i = j // skip up to this character in the parent loop
						break
					}
				}
				if len(detailsStr) == 0 {
					return fmt.Errorf("invalid annotation details format, closing bracket ')' not found in `%s`", str[i+1:])
				}
				s.name = strings.ToLower(strings.TrimSpace(s.name))
				s.value.Details = make(map[string]*Annotation)
				var supportedDetails map[string]bool
				if s.name == "relation" {
					supportedDetails = map[string]bool{"to": true, "name": true, "uid": true, "external-name": true, "external-type": true}
				} else if s.name == "sync" {
					supportedDetails = map[string]bool{"sharedglobalids": true}
				} else if s.name == "id" {
					supportedDetails = map[string]bool{"assignable": true}
				} else {
					return fmt.Errorf("invalid annotation format: details only supported for `relation` & `sync` annotations, found `%s`", s.name)
				}
				if err := ParseAnnotations(detailsStr, &s.value.Details, supportedDetails); err != nil {
					return err
				}
				if s.name == "relation" {
					if s.value.Details["name"] == nil {
						return fmt.Errorf("invalid annotation format: relation name missing in `%s`", str)
					}
					s.key = fmt.Sprintf("relation-%10d-%s", relationsCount(*annotations), s.value.Details["name"].Value)
				}
				if err := s.finishAnnotation(annotations, supportedAnnotations); err != nil {
					return err
				}
				s = annotationInProgress{} // reset
			} else if j := skipSpacesUntil(str, i+1, func(c uint8) bool {
				return c != ' '
			}); j != i {
				i = j - 1 // continue processing on the next non-space character
			}

		} else if !s.valueQuoted && (char == ',' || char == ' ') { // finish an annotation on a separator
			// A space may also be used before an equal sign, which means this isn't an annotation separator after all.
			if char == ' ' {
				// look ahead for a next character; if it's an equal sign, continue reading the same annotation
				if j := skipSpacesUntil(str, i, func(c uint8) bool {
					return c == '='
				}); j != i {
					i = j - 1 // continue processing on the equal sign
					continue
				}
			}
			if err := s.finishAnnotation(annotations, supportedAnnotations); err != nil {
				return err
			}
			s = annotationInProgress{} // reset
		} else if s.value != nil { // continue a value (set contents)
			if char == '"' {
				if len(s.value.Value) == 0 {
					s.valueQuoted = true
				} else {
					s.valueQuoted = false
					s.valueFinished = true
				}
			} else if s.valueFinished {
				return fmt.Errorf("invalid annotation format: no more characters may follow after a quoted value at position %d in `%s`", i, str)
			} else {
				s.value.Value += string(char)
			}
		} else { // continue a name
			s.name += string(char)
		}
	}

	return s.finishAnnotation(annotations, supportedAnnotations)
}

type annotationInProgress struct {
	name          string
	key           string
	value         *Annotation
	valueQuoted   bool
	valueFinished bool
}

// Skips spaces until encountering an equal sign, returning the position of the equal sign if found, or the input `i` otherwise.
func skipSpacesUntil(str string, i int, fn func(uint8) bool) int {
	for j := i; j < len(str); j++ {
		if str[j] == ' ' {
			continue
		} else if fn(str[j]) {
			return j
		} else {
			break
		}
	}
	return i
}

func (s *annotationInProgress) finishAnnotation(annotations *map[string]*Annotation, supportedAnnotations map[string]bool) error {
	s.name = strings.ToLower(strings.TrimSpace(s.name))
	if len(s.name) == 0 {
		return nil
	}
	if s.value == nil {
		s.value = &Annotation{} // empty value
	} else {
		s.value.Value = strings.TrimSpace(s.value.Value)
	}
	var key = s.key
	if len(key) == 0 {
		key = s.name
	}
	if (*annotations)[key] != nil {
		return fmt.Errorf("duplicate annotation %s", key)
	} else if !supportedAnnotations[s.name] {
		return fmt.Errorf("unknown annotation '%s'", s.name)
	} else {
		if strings.HasPrefix(key, "hnsw-") {
			var indexAnnotation = (*annotations)["index"]
			if indexAnnotation == nil || indexAnnotation.Value != "hnsw" {
				return fmt.Errorf("The HNSW annotation '%s' is only allowed after an 'index' annotation set to 'hnsw'.", key)
			}
		}
		(*annotations)[key] = s.value
	}
	return nil
}

// counts all "relation-" prefixed annotations (standalone relations) - used to ensure consistent processing order
func relationsCount(annotations map[string]*Annotation) uint {
	var count uint
	for key := range annotations {
		if strings.HasPrefix(key, "relation-") {
			count++
		}
	}
	return count
}
