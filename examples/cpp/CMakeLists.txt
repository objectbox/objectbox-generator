cmake_minimum_required(VERSION 3.14)

project("objectbox-test" CXX)

include(FetchContent)
FetchContent_Declare(
    objectbox
    GIT_REPOSITORY https://github.com/objectbox/objectbox-c.git
    GIT_TAG        v0.21.0
)
FetchContent_MakeAvailable(objectbox)

list(APPEND CMAKE_MODULE_PATH ${CMAKE_CURRENT_LIST_DIR}/../../cmake)
find_package(ObjectBoxGenerator 0.14 REQUIRED)

add_executable("objectbox-test"
    "src/main.cpp"
)
target_link_libraries("objectbox-test" objectbox)

set_target_properties("objectbox-test" PROPERTIES
    CXX_STANDARD 14
    CXX_STANDARD_REQUIRED YES
)

option(DO_INSOURCE "" TRUE)

if (DO_INSOURCE)
    add_schema_files(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            src/schema/person.fbs
            src/schema/task.fbs
            src/schema/monster.fbs
            src/schema/another_monster.fbs
        INSOURCE
    )
else()
    add_schema_files(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            src/schema/person.fbs
            src/schema/task.fbs
            src/schema/monster.fbs
            src/schema/another_monster.fbs
    )
    target_include_directories(objectbox-test PRIVATE ${CMAKE_CURRENT_BINARY_DIR}/src)
endif()