cmake_minimum_required(VERSION 3.14)

project(objectbox-generator-cmake-cpp-multiple-schema-dirs CXX)

include(../common.cmake)

add_executable(objectbox-test src/main.cpp)
target_link_libraries(objectbox-test objectbox)

set_target_properties(objectbox-test PROPERTIES
    CXX_STANDARD 14
    CXX_STANDARD_REQUIRED YES
)

# Variants:
# - use OUTPUT_DIR 
if (DO_INSOURCE)
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            src/schema1/monster.fbs
            src/schema2/task.fbs
        INSOURCE
        OUTPUT_DIR
            src/generated
        OUTPUT_DIR_MODEL_JSON
            src/model
    )
    target_include_directories(objectbox-test PRIVATE src/generated)
else()
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            src/schema1/monster.fbs
        OUTPUT_DIR_MODEL_JSON
            src/model
    )
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            src/schema2/task.fbs
        OUTPUT_DIR_MODEL_JSON
            src/model
    )
endif()
