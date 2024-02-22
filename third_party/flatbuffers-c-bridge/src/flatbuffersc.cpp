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

#include "flatbuffersc.h"

#include <flatbuffers/idl.h>
#include <flatbuffers/util.h>

#include "utils.h"

void fbs_error_free(const char* error) {
    if (error == nullptr) return;
    if (error == ErrorAllocationError) return;
    if (error == ErrorUnknown) return;
    free((void*) error);
}

FBS_bytes* fbs_schema_parse_file(const char* filename, const char** out_error) {
    return runCpp(out_error, nullptr, [&]() -> FBS_bytes* {
        VERIFY_ARGUMENT_NOT_NULL(filename);

        std::string contents;
        if (!flatbuffers::LoadFile(filename, true, &contents)) {
            throw std::invalid_argument(std::string("unable to load file: ") + filename);
        }

        auto options = flatbuffers::IDLOptions();
        options.binary_schema_comments = true;  // include doc comments in the binary schema

        flatbuffers::Parser parser(options);
        if (!parser.Parse(contents.c_str(), nullptr, filename)) {
            throw std::runtime_error(parser.error_);
        }

        if (!parser.error_.empty()) {
            bool ignore = false;
            if (parser.has_warning_) {
                if (parser.error_.find("warning: field names should be lowercase snake_case, got:") != std::string::npos) {
                    ignore = true;
                }
            }
            if (!ignore) {
                throw std::runtime_error(parser.error_);
            }
        }

        parser.Serialize();

        size_t size = parser.builder_.GetSize();
        VERIFY_STATE(size > 0);
        return mallocedBytesCopy("schema", parser.builder_.GetBufferPointer(), size);
    });
}

void fbs_schema_free(FBS_bytes* schema) {
    if (schema == nullptr) return;
    free(schema);
}
