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
