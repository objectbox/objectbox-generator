/*
 * Copyright 2020 ObjectBox Ltd. All rights reserved.
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
            // TODO flatc.cpp issues a warning in this case...
            // Warn(parser.error_, false);
            throw std::runtime_error(parser.error_);
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
