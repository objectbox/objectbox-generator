/*
 * Copyright 2017 Google Inc. All rights reserved.
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
// This file is heavily based on flatbuffers src/flatc_main.cpp - it exposes the main() function

#include "bfbs_gen_lua.h"
#include "bfbs_gen_nim.h"
#include "flatbuffers/base.h"
#include "flatbuffers/code_generator.h"
#include "flatbuffers/flatc.h"
#include "flatbuffers/util.h"
#include "idl_gen_binary.h"
#include "idl_gen_cpp.h"
#include "idl_gen_csharp.h"
#include "idl_gen_dart.h"
#include "idl_gen_fbs.h"
#include "idl_gen_go.h"
#include "idl_gen_java.h"
#include "idl_gen_json_schema.h"
#include "idl_gen_kotlin.h"
#include "idl_gen_lobster.h"
#include "idl_gen_php.h"
#include "idl_gen_python.h"
#include "idl_gen_rust.h"
#include "idl_gen_swift.h"
#include "idl_gen_text.h"
#include "idl_gen_ts.h"

#include "flatbuffersc.h"
#include "utils.h"

namespace {
static void Warn(const flatbuffers::FlatCompiler* flatc, const std::string& warn, bool show_exe_name) {
    printf("flatc warning: %s\n", warn.c_str());
}

static void Error(const flatbuffers::FlatCompiler* flatc, const std::string& err, bool usage, bool show_exe_name) {
    printf("flatc error: %s\n", err.c_str());
    if (usage && flatc) {
        printf("%s", flatc->GetUsageString("").c_str());
    }
    throw std::runtime_error(err);
}
}  // namespace

namespace flatbuffers {
void LogCompilerWarn(const std::string &warn) {
    Warn(static_cast<const flatbuffers::FlatCompiler *>(nullptr), warn, true);
}
void LogCompilerError(const std::string &err) {
    Error(static_cast<const flatbuffers::FlatCompiler *>(nullptr), err, false,
          true);
}
}  // namespace flatbuffers

int fbs_flatc(const char** args, size_t count, const char** out_error) {
    int code = 1;
    runCpp(out_error, [&]() {

        const std::string flatbuffers_version(flatbuffers::FLATBUFFERS_VERSION());

        flatbuffers::FlatCompiler::InitParams params;
        params.warn_fn = Warn;
        params.error_fn = Error;

        flatbuffers::FlatCompiler flatc(params);

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{
                "b", "binary", "",
                "Generate wire format binaries for any data definitions" },
            flatbuffers::NewBinaryCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "c", "cpp", "",
                                        "Generate C++ headers for tables/structs" },
            flatbuffers::NewCppCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "n", "csharp", "",
                                        "Generate C# classes for tables/structs" },
            flatbuffers::NewCSharpCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "d", "dart", "",
                                        "Generate Dart classes for tables/structs" },
            flatbuffers::NewDartCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "", "proto", "",
                                        "Input is a .proto, translate to .fbs" },
            flatbuffers::NewFBSCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "g", "go", "",
                                        "Generate Go files for tables/structs" },
            flatbuffers::NewGoCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "j", "java", "",
                                        "Generate Java classes for tables/structs" },
            flatbuffers::NewJavaCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "", "jsonschema", "", "Generate Json schema" },
            flatbuffers::NewJsonSchemaCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "", "kotlin", "",
                                        "Generate Kotlin classes for tables/structs" },
            flatbuffers::NewKotlinCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "", "lobster", "",
                                        "Generate Lobster files for tables/structs" },
            flatbuffers::NewLobsterCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "l", "lua", "",
                                        "Generate Lua files for tables/structs" },
            flatbuffers::NewLuaBfbsGenerator(flatbuffers_version));

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "", "nim", "",
                                        "Generate Nim files for tables/structs" },
            flatbuffers::NewNimBfbsGenerator(flatbuffers_version));

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "p", "python", "",
                                        "Generate Python files for tables/structs" },
            flatbuffers::NewPythonCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "", "php", "",
                                        "Generate PHP files for tables/structs" },
            flatbuffers::NewPhpCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "r", "rust", "",
                                        "Generate Rust files for tables/structs" },
            flatbuffers::NewRustCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{
                "t", "json", "", "Generate text output for any data definitions" },
            flatbuffers::NewTextCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "", "swift", "",
                                        "Generate Swift files for tables/structs" },
            flatbuffers::NewSwiftCodeGenerator());

        flatc.RegisterCodeGenerator(
            flatbuffers::FlatCOption{ "T", "ts", "",
                                        "Generate TypeScript code for tables/structs" },
            flatbuffers::NewTsCodeGenerator());

        // We need to inject main's argv[0] into args vector.
        std::vector<const char*> mainArgs(count+1);
        mainArgs[0] = "objectbox-generator";
        std::copy(args, args+count, &mainArgs[1]);

        // Create the FlatC options by parsing the command line arguments.
        const flatbuffers::FlatCOptions &options =
            flatc.ParseFromCommandLineArguments(count+1, mainArgs.data());

        // Compile with the extracted FlatC options.
        code = flatc.Compile(options);

    });
    fflush(stdout);
    fflush(stderr);
    return code;
}
