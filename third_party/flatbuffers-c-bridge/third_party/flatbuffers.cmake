# Cmake FetchContent is nice to use but not supported in pre-v3.11 cmake.
# Instead we keep a copy of some FB sources in this repository.
# Current version: 1.12.0

set(FLATBUFFERS_SRC_DIR ${CMAKE_CURRENT_LIST_DIR}/flatbuffers)

# see FlatBuffers_Library_SRCS in flatbuffers cmake
project(flatbuffers)
add_library(${PROJECT_NAME} STATIC
        ${FLATBUFFERS_SRC_DIR}/include/flatbuffers/base.h
        ${FLATBUFFERS_SRC_DIR}/include/flatbuffers/flatbuffers.h
        ${FLATBUFFERS_SRC_DIR}/include/flatbuffers/hash.h
        ${FLATBUFFERS_SRC_DIR}/include/flatbuffers/idl.h
        ${FLATBUFFERS_SRC_DIR}/include/flatbuffers/util.h
        ${FLATBUFFERS_SRC_DIR}/include/flatbuffers/reflection.h
        ${FLATBUFFERS_SRC_DIR}/include/flatbuffers/reflection_generated.h
        ${FLATBUFFERS_SRC_DIR}/include/flatbuffers/stl_emulation.h
        ${FLATBUFFERS_SRC_DIR}/include/flatbuffers/flexbuffers.h
        ${FLATBUFFERS_SRC_DIR}/include/flatbuffers/registry.h
        ${FLATBUFFERS_SRC_DIR}/include/flatbuffers/minireflect.h
        ${FLATBUFFERS_SRC_DIR}/src/idl_parser.cpp
        ${FLATBUFFERS_SRC_DIR}/src/idl_gen_text.cpp
        ${FLATBUFFERS_SRC_DIR}/src/reflection.cpp
        ${FLATBUFFERS_SRC_DIR}/src/util.cpp
        )
target_include_directories(${PROJECT_NAME} PUBLIC ${FLATBUFFERS_SRC_DIR}/include)

