/*
 * Copyright (C) 2020 ObjectBox Ltd. All rights reserved.
 * https://objectbox.io
 *
 * This file is part of ObjectBox Generator.
 *
 * ObjectBox Generator is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * ObjectBox Generator is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with ObjectBox Generator.  If not, see <http://www.gnu.org/licenses/>.
 */

package typeful

import (
	"testing"

	"github.com/objectbox/objectbox-generator/test/integration"
)

func TestCpp(t *testing.T) {
	conf := &integration.CCppTestConf{
		SchemaPath: "schema/schema.fbs",
	}
	defer conf.Cleanup()
	conf.CommonExecute(t, integration.Cpp11)
}
