#ifndef FLATBUFFERSC_H
#define FLATBUFFERSC_H

#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/// Utility container for any data
typedef struct FBS_bytes {
    uint64_t size;
    void* data;
} FBS_bytes;

void fbs_error_free(const char* error);

/// Parses a FlatBuffers schema file.
/// @param filename absolute path to a schema file file.
/// @param out_error - error if any occurred in which case you must free it using fbs_error_free() after reading.
/// @return a pointer to the loaded FB of the schema. Must be freed after use by calling fbs_schema_free()
FBS_bytes* fbs_schema_parse_file(const char** out_error, const char* filename);

/// Frees memory of both FBS_bytes as well as the inner schema->data
void fbs_schema_free(FBS_bytes* schema);

#ifdef __cplusplus
}
#endif
#endif  // FLATBUFFERSC_H
