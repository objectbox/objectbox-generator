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

#include <stdio.h>
#include <catch.hpp>
#include "flatbuffersc.h"

namespace {
/// Provides a memory-managed wrapper to be used with the C-API argument `char** out_error`
struct Error {
    const char* text = nullptr;
    ~Error() { fbs_error_free(text); }

    /// @returns a valid pointer that can be used as a `char** out_error` argument for C calls
    /// @note in case the internal error text was previously set (an error occurred before), it's freed first before
    ///       returning a clean pointer
    const char** ptr() {
        if (text) {
            fbs_error_free(text);
            text = nullptr;
        }
        return &text;
    }
};

bool fileExists(const char* path) {
    if (FILE* file = fopen(path, "r")) {
        fclose(file);
        return true;
    }
    return false;
}

}  // namespace

TEST_CASE("schema-parser-errors", "") {
    // must not crash even when an error occurs and out_error is nullptr
    FBS_bytes* schema = fbs_schema_parse_file(nullptr, nullptr);
    REQUIRE(schema == nullptr);

    // missing filename must produce an error
    Error error;
    schema = fbs_schema_parse_file(nullptr, error.ptr());
    REQUIRE(schema == nullptr);
    REQUIRE(error.text != nullptr);
    REQUIRE_THAT(error.text, Catch::Contains("must not be null"));
}

TEST_CASE("schema-parser", "") {
    Error error;
    FBS_bytes* schema = fbs_schema_parse_file(TEST_SRC_DIRECTORY "schema.fbs", error.ptr());
    REQUIRE(error.text == nullptr);
    REQUIRE(schema != nullptr);
    REQUIRE(schema->size > 0);
    REQUIRE(schema->data != nullptr);

    // A rough check the comments are included correctly
    std::string str(static_cast<char*>(schema->data), schema->size);
    REQUIRE(schema->size == str.size());

    // `///` three slashes on an otherwise blank line a comment
    REQUIRE_THAT(str, Catch::Contains("A real or imaginary living creature or entity"));
    REQUIRE_THAT(str, Catch::Contains("Note: name may be nil"));
    REQUIRE_THAT(str, Catch::Contains("All worldly belongings of this being"));

    // `//` Not a doc comment because it doesn't have three slashes
    REQUIRE_THAT(str, !Catch::Contains("An individual article"));

    // `//<` doesn't make a comment at the end of the line
    REQUIRE_THAT(str, !Catch::Contains("Current health points"));

    // `/**` multiline comments don't work
    REQUIRE_THAT(str, !Catch::Contains("A celestial body"));
    REQUIRE_THAT(str, !Catch::Contains("Distinguished"));

    fbs_schema_free(schema);
}

TEST_CASE("flatc-main", "") {
    Error error;

    const char* generatedFile = "schema_generated.h";
    remove(generatedFile);

    const size_t argc = 2;
    const char* args[] = {"--cpp", TEST_SRC_DIRECTORY "schema.fbs"};
    int code = fbs_flatc(args, argc, error.ptr());
    CAPTURE(error.text);
    REQUIRE(code == 0);

    REQUIRE(fileExists(generatedFile));
}