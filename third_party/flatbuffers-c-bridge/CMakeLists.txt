cmake_minimum_required(VERSION 3.5)

set(CMAKE_C_STANDARD 11)
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_POSITION_INDEPENDENT_CODE ON)

include(third_party/flatbuffers.cmake)

project(flatbuffers-c-bridge C)
add_library(${PROJECT_NAME} STATIC
        src/flatbuffersc.cpp
        )
target_include_directories(${PROJECT_NAME} PUBLIC include)
target_link_libraries(${PROJECT_NAME} PRIVATE flatbuffers)

# Define flatbuffers-compiler library which we need to provide fbs_flatc() function
# see FlatBuffers_Compiler_SRCS in flatbuffers CmakeLists.txt
project(flatbuffers-c-bridge-flatc)
add_library(${PROJECT_NAME} STATIC
        src/flatc_main.cpp
        ${FLATBUFFERS_SRC_DIR}/src/annotated_binary_text_gen.cpp
        ${FLATBUFFERS_SRC_DIR}/src/bfbs_gen_lua.cpp
        ${FLATBUFFERS_SRC_DIR}/src/bfbs_gen_nim.cpp
        ${FLATBUFFERS_SRC_DIR}/src/binary_annotator.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_binary.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_cpp.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_csharp.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_dart.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_kotlin.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_go.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_java.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_php.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_python.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_lobster.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_rust.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_fbs.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_grpc.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_json_schema.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_swift.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_text.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_ts.cpp
        ${FLATBUFFERS_SRC_DIR}/src/flatc.cpp
        ${FLATBUFFERS_SRC_DIR}/include/flatbuffers/code_generators.h
        ${FLATBUFFERS_SRC_DIR}/src/code_generators.cpp
        ${FLATBUFFERS_SRC_DIR}/grpc/src/compiler/schema_interface.h
        ${FLATBUFFERS_SRC_DIR}/grpc/src/compiler/cpp_generator.h
        ${FLATBUFFERS_SRC_DIR}/grpc/src/compiler/cpp_generator.cc
        ${FLATBUFFERS_SRC_DIR}/grpc/src/compiler/go_generator.h
        ${FLATBUFFERS_SRC_DIR}/grpc/src/compiler/go_generator.cc
        ${FLATBUFFERS_SRC_DIR}/grpc/src/compiler/java_generator.h
        ${FLATBUFFERS_SRC_DIR}/grpc/src/compiler/java_generator.cc
        ${FLATBUFFERS_SRC_DIR}/grpc/src/compiler/python_generator.h
        ${FLATBUFFERS_SRC_DIR}/grpc/src/compiler/python_generator.cc
        ${FLATBUFFERS_SRC_DIR}/grpc/src/compiler/swift_generator.h
        ${FLATBUFFERS_SRC_DIR}/grpc/src/compiler/swift_generator.cc
        ${FLATBUFFERS_SRC_DIR}/grpc/src/compiler/ts_generator.cc

        )
target_link_libraries(${PROJECT_NAME} PUBLIC flatbuffers)
target_include_directories(${PROJECT_NAME} PRIVATE include ${FLATBUFFERS_SRC_DIR}/src ${FLATBUFFERS_SRC_DIR}/grpc)

include(test/CMakeLists.txt)