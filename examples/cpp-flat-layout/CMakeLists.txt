cmake_minimum_required(VERSION 3.14)

project("objectbox-generator-cpp-flat-layout" CXX)

include(FetchContent)
FetchContent_Declare(
    objectbox
    GIT_REPOSITORY https://github.com/objectbox/objectbox-c.git
    GIT_TAG        v5.0.0-rc
)
FetchContent_MakeAvailable(objectbox)

list(APPEND CMAKE_MODULE_PATH ${CMAKE_CURRENT_LIST_DIR}/../../cmake)
set(ObjectBoxGenerator_ROOT ${CMAKE_CURRENT_LIST_DIR}/../..)
find_package(ObjectBoxGenerator 4.0.0 REQUIRED)

add_executable(objectbox-test main.cpp)
target_link_libraries(objectbox-test objectbox)

set_target_properties(objectbox-test PROPERTIES
    CXX_STANDARD 14
    CXX_STANDARD_REQUIRED YES
)

option(DO_INSOURCE "" FALSE)

if (DO_INSOURCE)
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            task.fbs
        INSOURCE
    )
else()
    add_obx_schema(
        TARGET 
            objectbox-test
        SCHEMA_FILES 
            task.fbs
    )
endif()
