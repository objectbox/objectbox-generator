cmake_minimum_required(VERSION 3.14)

project(objectbox-generator-cmake-cpp-flat CXX)

include(../common.cmake)

add_executable(objectbox-test main.cpp)
target_link_libraries(objectbox-test objectbox)

set_target_properties(objectbox-test PROPERTIES
    CXX_STANDARD 14
    CXX_STANDARD_REQUIRED YES
)

# Variants:
# - multiple add_obx_schema calls on same target
# - single add_obx_schema call adding two schemas
# - use OUTPUT_DIR_HEADERS (without OUTPUT_DIR)
if (DO_INSOURCE)
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            task.fbs
        INSOURCE
        OUTPUT_DIR
            schema1
    )
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            monster.fbs
        INSOURCE
        OUTPUT_DIR
            schema2
    )
    target_include_directories(objectbox-test PRIVATE schema1 schema2)
else()
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            task.fbs
            monster.fbs
    )
endif()
