# Find Module for tool `objectbox-generator` 
# with "opt-in" automatic fetch and prepare for local execution, controlled by variable "ObjectBoxGenerator_AllowFetch" (defaults to true).

option(OBX_GENERATOR_ALLOW_FETCH "Opt-in automatic download and prepare for local execution" ON)

# Updated by maintainer to latest available version:
set(ObjectBoxGenerator_FETCH_VERSION 0.14.0)
set(ObjectBoxGenerator_FETCH_HASH_Linux SHA256=2a06d06567524c9fa8bd815007a9c8a304792b7c0e49fade775a7553026b419a)
# TODO: update for Windows/macOS
set(ObjectBoxGenerator_FETCH_HASH_Windows SHA256=2a06d06567524c9fa8bd815007a9c8a304792b7c0e49fade775a7553026b419a)
set(ObjectBoxGenerator_FETCH_HASH_macOS SHA256=2a06d06567524c9fa8bd815007a9c8a304792b7c0e49fade775a7553026b419a)
set(ObjectBoxGenerator_FETCH_DIR ${CMAKE_BINARY_DIR}/ObjectBoxGenerator-Fetch)
set(ObjectBoxGenerator_FETCH_BASEURL "https://github.com/objectbox/objectbox-generator/releases/download")

include(FindPackageHandleStandardArgs)

# Find program also in fetch dir.
find_program(ObjectBoxGenerator_EXECUTABLE objectbox-generator PATHS ${ObjectBoxGenerator_FETCH_DIR} NO_CACHE)

if(ObjectBoxGenerator_EXECUTABLE)
  execute_process(
    COMMAND ${ObjectBoxGenerator_EXECUTABLE} -version 
    OUTPUT_VARIABLE Output 
    OUTPUT_STRIP_TRAILING_WHITESPACE 
    ECHO_ERROR_VARIABLE 
    COMMAND_ERROR_IS_FATAL ANY)
  string(REGEX REPLACE ".* v\([0-9\.]*\).*" "\\1" ObjectBoxGenerator_VERSION ${Output})
  find_package_check_version(${ObjectBoxGenerator_VERSION} ObjectBoxGenerator_VERSION_OK)
  if(NOT ${ObjectBoxGenerator_VERSION_OK}) 
    set(ObjectBoxGenerator_FETCH_REQUIRED TRUE)
  endif()
else()
  set(ObjectBoxGenerator_FETCH_REQUIRED TRUE)
endif()
  
if(OBX_GENERATOR_ALLOW_FETCH AND ObjectBoxGenerator_FETCH_REQUIRED)
  message(STATUS "ObjectBox-Generator Fetch: Executable not found, attempting to download to build directory and prepare for execution (to disable behaviour set OBX_GENERATOR_ALLOW_FETCH to OFF)")
  if (CMAKE_HOST_LINUX AND CMAKE_HOST_SYSTEM_PROCESSOR STREQUAL "x86_64")
    set(ObjectBoxGenerator_FETCH_ARCH Linux)
  elseif (CMAKE_HOST_WIN32 AND CMAKE_HOST_SYSTEM_PROCESSOR STREQUAL "x86_64")
    # TODO: check windows
    set(ObjectBoxGenerator_FETCH_ARCH Windows)
  elseif (CMAKE_HOST_APPLE)
    set(ObjectBoxGenerator_FETCH_ARCH macOS)
    # TODO: check apple
  else()
    # TODO: clarify details
    message(FATAL_ERROR "ObjectBoxGenerator Fetch failed: unsupported platform (must be Linux/x86-64, Windows/x86-64 or macOS)")
  endif()
  if(ObjectBoxGenerator_FIND_VERSION)
    if(${ObjectBoxGenerator_FIND_VERSION} VERSION_GREATER ${ObjectBoxGenerator_FETCH_VERSION})
      message(FATAL_ERROR "Requested ObjectBox-Generator version '${ObjectBoxGenerator_FIND_VERSION}' not available. Latest version is ${ObjectBoxGenerator_FETCH_VERSION}")
    endif()
  endif()
  message(STATUS "ObjectBox-Generator Fetch: Downloading version ${ObjectBoxGenerator_FETCH_VERSION}")
  set(ObjectBoxGenerator_FETCH_FILE objectbox-generator-${ObjectBoxGenerator_FETCH_ARCH}.zip)
  set(ObjectBoxGenerator_FETCH_URL ${ObjectBoxGenerator_FETCH_BASEURL}/v${ObjectBoxGenerator_FETCH_VERSION}/${ObjectBoxGenerator_FETCH_FILE})
  set(ObjectBoxGenerator_FETCH_PATH ${ObjectBoxGenerator_FETCH_DIR}/${ObjectBoxGenerator_FETCH_FILE})
  set(ObjectBoxGenerator_UNPACK_FILE ${ObjectBoxGenerator_FETCH_DIR}/objectbox-generator)
  message(STATUS "ObjectBox-Generator Fetch: Downloading archive from ${ObjectBoxGenerator_FETCH_URL} to ${ObjectBoxGenerator_FETCH_PATH}")
  file(DOWNLOAD "${ObjectBoxGenerator_FETCH_URL}" "${ObjectBoxGenerator_FETCH_PATH}"
        TLS_VERIFY ON 
        EXPECTED_HASH "${ObjectBoxGenerator_FETCH_HASH_${ObjectBoxGenerator_FETCH_ARCH}}"
  )
  message(STATUS "ObjectBox-Generator Fetch: Unpacking archive file ${ObjectBoxGenerator_FETCH_PATH} to ${ObjectBoxGenerator_UNPACK_FILE}")
  file(ARCHIVE_EXTRACT INPUT ${ObjectBoxGenerator_FETCH_PATH} DESTINATION ${ObjectBoxGenerator_FETCH_DIR} VERBOSE)
  message(STATUS "ObjectBox-Generator Fetch: Make file ${ObjectBoxGenerator_UNPACK_FILE} executable, and retry location..")
  file(CHMOD ${ObjectBoxGenerator_UNPACK_FILE} PERMISSIONS OWNER_READ OWNER_EXECUTE GROUP_READ GROUP_EXECUTE WORLD_READ WORLD_EXECUTE)
  find_program(ObjectBoxGenerator_EXECUTABLE objectbox-generator PATHS ${ObjectBoxGenerator_FETCH_DIR})
endif()

find_package_handle_standard_args(ObjectBoxGenerator 
  REQUIRED_VARS ObjectBoxGenerator_EXECUTABLE
  VERSION_VAR ObjectBoxGenerator_VERSION
)

if (ObjectBoxGenerator_FOUND)
  include(obxGenerator)
endif()

