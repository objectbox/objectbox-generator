cmake_minimum_required(VERSION 3.14)

project(objectbox-generator-cmake-cpp-flat CXX)

include(../common.cmake)

add_executable(objectbox-test main.cpp)
target_link_libraries(objectbox-test objectbox)

set_target_properties(objectbox-test PROPERTIES
    CXX_STANDARD 14
    CXX_STANDARD_REQUIRED YES
)

# Also use two different variants here:
# - multiple add_obx_schema calls on same target
# - one add_obx_schema call adding two schemas
if (DO_INSOURCE)
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            task.fbs
        INSOURCE
    )
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            monster.fbs
        INSOURCE
    )
else()
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            task.fbs
            monster.fbs
    )
    target_include_directories(objectbox-test PRIVATE ${CMAKE_CURRENT_BINARY_DIR})
endif()
