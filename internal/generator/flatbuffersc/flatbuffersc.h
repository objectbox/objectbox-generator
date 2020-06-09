#ifndef FLATBUFFERSC_H
#define FLATBUFFERSC_H

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
#endif  // FLATBUFFERSC_H
