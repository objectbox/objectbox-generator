/*
 * ObjectBox Generator - a build time tool for ObjectBox
 * Copyright (C) 2020-2024 ObjectBox Ltd. All rights reserved.
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

package comparison

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CopyFile(sourceFile, targetFile string, permsOverride os.FileMode) error {
	data, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	// copy permissions either from the existing target file or from the source file
	var perm os.FileMode = permsOverride
	if perm == 0 {
		if info, _ := os.Stat(targetFile); info != nil {
			perm = info.Mode()
		} else if info, err := os.Stat(sourceFile); info != nil {
			perm = info.Mode()
		} else {
			return err
		}
	}

	err = ioutil.WriteFile(targetFile, data, perm)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func copyDirectory(sourceDir, targetDir string, dirPerms, filePerms os.FileMode) error {
	if err := os.MkdirAll(targetDir, dirPerms); err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDir, entry.Name())
		targetPath := filepath.Join(targetDir, entry.Name())

		info, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			if err := copyDirectory(sourcePath, targetPath, dirPerms, filePerms); err != nil {
				return err
			}
		} else if info.Mode().IsRegular() {
			if err := CopyFile(sourcePath, targetPath, filePerms); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("not a regular file or directory: %s", sourcePath)
		}
	}
	return nil
}
