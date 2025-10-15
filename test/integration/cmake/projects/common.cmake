if(NOT TOPDIR)
  set(TOPDIR ${CMAKE_CURRENT_LIST_DIR}/../../../..)
endif()

option(DO_INSOURCE "Generate files in-source" FALSE)
option(NO_PREBUILT "Use pre-built generator" FALSE)
include(FetchContent)
FetchContent_Declare(
    objectbox
    GIT_REPOSITORY https://github.com/objectbox/objectbox-c.git
    GIT_TAG        v5.0.0-rc
)
FetchContent_MakeAvailable(objectbox)

# Use find module from source tree.
list(APPEND CMAKE_MODULE_PATH ${TOPDIR}/cmake)

# Use pre-built ObjectBox Generator.
if(NOT NO_PREBUILT)
    set(ObjectBoxGenerator_ROOT ${TOPDIR})
endif()

find_package(ObjectBoxGenerator 4.0.0 REQUIRED)
