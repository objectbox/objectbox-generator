cmake_minimum_required(VERSION 3.14)

project(objectbox-generator-cmake-cpp-tree CXX)

include(../common.cmake)

add_executable(objectbox-test src/main.cpp)
target_link_libraries(objectbox-test objectbox)

set_target_properties(objectbox-test PROPERTIES
    CXX_STANDARD 14
    CXX_STANDARD_REQUIRED YES
)

# Variants:
# - one add_obx_schema call adding two schemas
# - use OUTPUT_DIR and OUTPUT_DIR_HEADERS
# - use EXTRA_OPTIONS
if (DO_INSOURCE)
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            src/schema/person.fbs
            src/schema/task.fbs
            src/schema/monster.fbs
            src/schema/another_monster.fbs
        INSOURCE
        OUTPUT_DIR
            src
        OUTPUT_DIR_HEADERS
            include
        EXTRA_OPTIONS
          -empty-string-as-null
          -nan-as-null
          -optional std::shared_ptr
    )
    target_include_directories(objectbox-test PRIVATE include)
else()
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            src/schema/person.fbs
            src/schema/task.fbs
            src/schema/monster.fbs
            src/schema/another_monster.fbs
        EXTRA_OPTIONS
          -empty-string-as-null
          -nan-as-null
          -optional std::unique_ptr
    )
endif()
