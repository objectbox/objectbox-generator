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

#pragma once

#ifdef __cplusplus
#include <cstddef>
#include <cstdint>
extern "C" {
#else
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#endif

/// Utility container for any data
typedef struct FBS_bytes {
    uint64_t size;
    void* data;
} FBS_bytes;

void fbs_error_free(const char* error);

/// Parses a FlatBuffers schema file.
/// @param out_error - error if any occurred in which case you must free it using fbs_error_free() after reading.
/// @param filename absolute path to a schema file file.
/// @return a pointer to the loaded FB of the schema. Must be freed after use by calling fbs_schema_free()
FBS_bytes* fbs_schema_parse_file(const char* filename, const char** out_error);

/// Frees memory of both FBS_bytes as well as the inner schema->data
void fbs_schema_free(FBS_bytes* schema);

/// Executes FlatBuffers compiler's (flatc executable) main() function with the given arguments.
/// @warning you must link to flatbuffersc-flatc in addition to flatbuffersc if you intend to use this function
/// @note may print warnings ond errors on standard output
/// @param args should contain arguments WITHOUT the current program name (as opposed to the usual main() signature that
/// has the program name at the first index)
/// @returns an exit code as returned by the FB compiler, zero on success
int fbs_flatc(const char** args, size_t count, const char** out_error);

#ifdef __cplusplus
}
#endif
