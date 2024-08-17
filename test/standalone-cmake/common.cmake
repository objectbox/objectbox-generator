option(DO_INSOURCE "Generate files in-source" FALSE)
option(NO_PREBUILT "Use pre-built generator" FALSE)
include(FetchContent)
FetchContent_Declare(
    objectbox
    GIT_REPOSITORY https://github.com/objectbox/objectbox-c.git
    GIT_TAG        v4.0.1
)
FetchContent_MakeAvailable(objectbox)

# Use find module from source tree.
list(APPEND CMAKE_MODULE_PATH ${CMAKE_CURRENT_LIST_DIR}/../../cmake)

# Use pre-built ObjectBox Generator.
if(NOT NO_PREBUILT)
    set(ObjectBoxGenerator_ROOT ${CMAKE_CURRENT_LIST_DIR}/../..)
endif()

find_package(ObjectBoxGenerator 4.0.0 REQUIRED)
