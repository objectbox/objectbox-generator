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

#define VERIFY_ARGUMENT_NOT_NULL(condition) \
    ((condition) ? (void) (0) : throw std::invalid_argument(std::string("Argument #condition must not be null")))

#define VERIFY_STATE(condition) \
    ((condition) ? (void) (0) : throw std::runtime_error(std::string("State condition failed: #condition")))

namespace {
static const char* ErrorAllocationError = "Out-of-memory when trying to allocate an error text.";
static const char* ErrorUnknown = "Unknown error";

const char* mallocedString(const char* str1, const char* str2 = nullptr) {
    size_t len1 = strlen(str1);
    size_t len2 = str2 ? strlen(str2) : 0;

    char* ptr = static_cast<char*>(malloc(len1 + len2 + 1));
    if (!ptr) return ErrorAllocationError;

    memcpy(ptr, str1, len1);
    if (str2) memcpy(ptr + len1, str2, len2);
    memset(ptr + len1 + len2, 0, 1);  // null terminate the resulting string

    return ptr;
}

inline size_t paddedSize(size_t len, size_t padding) { return (len + padding - 1) & (~(padding - 1)); }

FBS_bytes* mallocedBytesCopy(const char* name, void* src, size_t size) {
    size_t paddedStructSize = paddedSize(sizeof(FBS_bytes), sizeof(void*));
    void* memory = malloc(paddedStructSize + size);
    if (!memory) {
        throw std::runtime_error("Could not allocate " + std::to_string(size) + " bytes for the " + name);
    }
    FBS_bytes* result = static_cast<FBS_bytes*>(memory);
    result->size = size;
    if (size == 0) {
        result->data = nullptr;
    } else {
        result->data = static_cast<uint8_t*>(memory) + paddedStructSize;
        memcpy(result->data, src, size);
    }
    return result;
}

template <typename FUN, typename RET>
auto runCpp(const char** outError, RET&& resultOnError, FUN fn) -> decltype(fn()) {
    try {
        return fn();
    } catch (const std::exception& e) {
        // NOTE: ignore clion "condition is always true", it isn't...
        if (outError) *outError = mallocedString(e.what());
    } catch (...) {
        // NOTE: ignore clion "condition is always true", it isn't...
        if (outError) *outError = ErrorUnknown;
    }
    return std::forward<RET>(resultOnError);
}

template <typename FUN>
void runCpp(const char** outError, FUN fn) {
    runCpp(outError, 0, [&]() {
        fn();
        return 0;
    });
}
}  // namespace
