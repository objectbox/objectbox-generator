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
set(ObjectBoxGenerator_INSTALL_DIR ${CMAKE_BINARY_DIR}/ObjectBoxGenerator-Install)

include(FindPackageHandleStandardArgs)

# Find program also in fetch dir.
find_program(ObjectBoxGenerator_EXECUTABLE objectbox-generator PATHS ${ObjectBoxGenerator_INSTALL_DIR} NO_CACHE)

if(ObjectBoxGenerator_EXECUTABLE)
  execute_process(
    COMMAND ${ObjectBoxGenerator_EXECUTABLE} -version 
    OUTPUT_VARIABLE Output 
    OUTPUT_STRIP_TRAILING_WHITESPACE 
    RESULT_VARIABLE ExitStatus
    # 3.19: COMMAND_ERROR_IS_FATAL ANY
  )
  if(NOT ${ExitStatus} EQUAL 0) 
    message(WARNING "Unable to query version on ${ObjectBoxGenerator_EXECUTABLE}")
    return()
  endif()

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
  if (CMAKE_HOST_SYSTEM_NAME STREQUAL "Linux" AND CMAKE_HOST_SYSTEM_PROCESSOR STREQUAL "x86_64")
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
  execute_process(COMMAND ${CMAKE_COMMAND} -E tar xf ${ObjectBoxGenerator_FETCH_PATH} WORKING_DIRECTORY ${ObjectBoxGenerator_FETCH_DIR})
  # 3.18: file(ARCHIVE_EXTRACT INPUT ${ObjectBoxGenerator_FETCH_PATH} DESTINATION ${ObjectBoxGenerator_FETCH_DIR} VERBOSE)
  message(STATUS "ObjectBox-Generator Fetch: Make file ${ObjectBoxGenerator_UNPACK_FILE} executable, and retry location..")
  file(COPY ${ObjectBoxGenerator_UNPACK_FILE} DESTINATION ${ObjectBoxGenerator_INSTALL_DIR} FILE_PERMISSIONS OWNER_READ OWNER_EXECUTE GROUP_READ GROUP_EXECUTE WORLD_READ WORLD_EXECUTE)
  # 3.19: file(CHMOD ${ObjectBoxGenerator_UNPACK_FILE} PERMISSIONS OWNER_READ OWNER_EXECUTE GROUP_READ GROUP_EXECUTE WORLD_READ WORLD_EXECUTE)
  find_program(ObjectBoxGenerator_EXECUTABLE objectbox-generator PATHS ${ObjectBoxGenerator_INSTALL_DIR})
endif()

find_package_handle_standard_args(ObjectBoxGenerator 
  REQUIRED_VARS ObjectBoxGenerator_EXECUTABLE
  VERSION_VAR ObjectBoxGenerator_VERSION
)

if (NOT ObjectBoxGenerator_FOUND)
  return()
endif()

function (add_schema_files)

  set(options INSOURCE)
  set(oneValueArgs TARGET)
  set(multiValueArgs SCHEMA_FILES)
  cmake_parse_arguments(ARG "${options}" "${oneValueArgs}" "${multiValueArgs}" ${ARGN})

  if (ARG_INSOURCE)	
    set(base_dir ${CMAKE_CURRENT_SOURCE_DIR})
  else()
    set(base_dir ${CMAKE_CURRENT_BINARY_DIR})
  endif()

  set(sources)

  foreach(SCHEMA_FILE ${ARG_SCHEMA_FILES})
    
    # 3.20: cmake_path(ABSOLUTE_PATH SCHEMA_FILE OUTPUT_VARIABLE schema_filepath)
    set(schema_filepath ${CMAKE_CURRENT_SOURCE_DIR}/${SCHEMA_FILE})

    string(REGEX REPLACE "\.fbs$" ".obx.cpp" cppfile ${base_dir}/${SCHEMA_FILE})
    string(REGEX REPLACE "\.fbs$" ".obx.hpp" hppfile ${base_dir}/${SCHEMA_FILE})
    
    # 3.20: cmake_path(GET cppfile PARENT_PATH out_dir)
    string(REGEX REPLACE "/[^/]*$" "" out_dir ${cppfile}) 

    if (NOT ARG_INSOURCE)
      file(MAKE_DIRECTORY ${out_dir})
    endif()
    add_custom_command(
      OUTPUT 
        ${cppfile} 
        ${hppfile} 
      COMMAND 
        ${ObjectBoxGenerator_EXECUTABLE} ARGS -out ${out_dir} -cpp ${schema_filepath}
      BYPRODUCTS 
        ${out_dir}/objectbox-model.h
        ${out_dir}/objectbox-model.json
      DEPENDS 
        ${schema_filepath}
    )
    list(APPEND sources ${cppfile} ${hppfile})
  endforeach()
    
  target_sources(${ARG_TARGET} PRIVATE ${sources}) 
endfunction()
