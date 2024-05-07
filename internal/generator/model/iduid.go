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

package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// IdUid represents a "ID:UID" string as used in the model jSON
type IdUid string

// CreateIdUid creates a string representation of ID and UID
func CreateIdUid(id Id, uid Uid) IdUid {
	return IdUid(strconv.FormatInt(int64(id), 10) + ":" + strconv.FormatUint(uid, 10))
}

var componentNamesErr = [2]string{"id", "uid"}

// Validate performs initial validation of loaded data so that it doesn't have to be checked in each function
func (str *IdUid) Validate() error {
	if _, err := str.GetUid(); err != nil {
		return err
	}

	if _, err := str.GetId(); err != nil {
		return err
	}

	if len(strings.Split(string(*str), ":")) != 2 {
		return errors.New("invalid id format - too many colons")
	}

	return nil
}

// GetId returns the ID part
func (str *IdUid) GetId() (Id, error) {
	id, err := str.getComponent(0, 32, false)
	if err != nil {
		return 0, err
	}
	return Id(id), nil
}

// GetId returns the ID part, not returning an error in case of a zero value
func (str *IdUid) GetIdAllowZero() (Id, error) {
	id, err := str.getComponent(0, 32, true)
	if err != nil {
		return 0, err
	}
	return Id(id), nil
}

// GetUid returns the UID part
func (str *IdUid) GetUid() (Uid, error) {
	return str.getComponent(1, 64, false)
}

// GetUid returns the UID part, not returning an error in case of a zero value
func (str *IdUid) GetUidAllowZero() (Uid, error) {
	return str.getComponent(1, 64, true)
}

// Get returns a pair of ID and UID
func (str *IdUid) Get() (Id, Uid, error) {
	if id, err := str.GetId(); err != nil {
		return 0, 0, err
	} else if uid, err := str.GetUid(); err != nil {
		return 0, 0, err
	} else {
		return id, uid, nil
	}
}

func (str IdUid) getComponent(n, bitsize int, allowZero bool) (uint64, error) {
	if len(str) == 0 {
		return 0, errors.New(componentNamesErr[n] + " is undefined")
	}

	idStr := strings.Split(string(str), ":")[n]
	if component, err := strconv.ParseUint(idStr, 10, bitsize); err != nil {
		return 0, fmt.Errorf("can't parse '%s' as unsigned int: %s", idStr, err)
	} else if component == 0 && !allowZero {
		return 0, errors.New(componentNamesErr[n] + " is zero")
	} else {
		return component, nil
	}
}

func (str IdUid) getIdSafe() Id {
	i, _ := str.GetId()
	return i
}

func (str IdUid) getUidSafe() Uid {
	i, _ := str.GetUid()
	return i
}
