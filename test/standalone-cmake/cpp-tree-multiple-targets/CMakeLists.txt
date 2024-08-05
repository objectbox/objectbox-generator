cmake_minimum_required(VERSION 3.14)

project(objectbox-tree-multiple-targets CXX)

include(../common.cmake)

add_executable(target1 src/main.cpp)
add_executable(target2 src/main.cpp)

target_link_libraries(target1 objectbox)
target_link_libraries(target2 objectbox)

set_target_properties(target1 PROPERTIES
    CXX_STANDARD 14
    CXX_STANDARD_REQUIRED YES
)
set_target_properties(target2 PROPERTIES
    CXX_STANDARD 14
    CXX_STANDARD_REQUIRED YES
)

if (DO_INSOURCE)
    add_obx_schema(
        TARGET 
            target1
        SCHEMA_FILES 
            src/schema/person.fbs
            src/schema/task.fbs
            src/schema/monster.fbs
            src/schema/another_monster.fbs
        INSOURCE
    )
    add_obx_schema(
        TARGET 
            target2
        SCHEMA_FILES 
            src/schema/person.fbs
            src/schema/task.fbs
            src/schema/monster.fbs
            src/schema/another_monster.fbs
        INSOURCE
    )
else()
    add_obx_schema(
        TARGET 
            target1
        SCHEMA_FILES 
            src/schema/person.fbs
            src/schema/task.fbs
            src/schema/monster.fbs
            src/schema/another_monster.fbs
    )
    target_include_directories(target1 PRIVATE ${CMAKE_CURRENT_BINARY_DIR}/src)
    add_obx_schema(
        TARGET 
            target2
        SCHEMA_FILES 
            src/schema/person.fbs
            src/schema/task.fbs
            src/schema/monster.fbs
            src/schema/another_monster.fbs
    )
    target_include_directories(target2 PRIVATE ${CMAKE_CURRENT_BINARY_DIR}/src)
endif()
