#[=======================================================================[.rst:
FindObjectBoxGenerator
----------------------

ObjectBox-Generator (ObjectBoxGenerator_) is a code generator tool to support 
C/C++ and Go development with ObjectBox (ObjectBox_), a superfast 
cross-platform object-oriented database.

This find module automatically locates a local installation of the 
executable ``objectbox-generator`` and checks requested version with found one.
In addition, it can automatically download and unpack a version into the build directory
and make it executable. Automatic download is enabled by default via the option 
``OBX_GENERATOR_ALLOW_FETCH``.

Currently supported platforms are Linux/x86-64, macOS and Windows/x86-64.

.. code-block:: cmake

  find_package(ObjectBoxGenerator 4.0.0 REQUIRED)

The following variables are defined by this module:

.. variable:: ObjectBoxGenerator_FOUND

  Whether objectbox-generator was successfully found.  

.. variable:: ObjectBoxGenerator_EXECUTABLE

  If found, this variable comprises the full path to executable.  

Functions
^^^^^^^^^

.. command:: add_obx_schema

This function adds ObjectBox schema files to a C++ target which 
implies a build task to auto-generate C++ source/header files 
from schema file(s) (with dependency rule tracking) and
adds them as sources to the target for compilation::

     add_obx_schema(
       TARGET <target>
       SCHEMA_FILES <schemafile>..
       [INSOURCE]
       [CXX_STANDARD 11|14]
       [EXTRA_OPTIONS <options>..]
     )
  
ObjectBox schema files have the filename pattern ``<name>.fbs`` 
which yields the name of auto-generated C++ source and header file 
using the pattern ``<name>.obx.cpp`` and ``<name>.obx.hpp``, respectively.

If the option ``INSOURCE`` is set then generated files are 
written relative to the current source directory, otherwise the
current binary directory is taken as base directory. 

In additon the generator also creates and updates the files 
``objectbox-model.h`` and  ``objectbox-model.json`` next to the 
generated C++ source/header files.

The option ``CXX_STANDARD`` may be set to ``11`` or ``14`` (default) to specify 
the base-line C++ language standard the code generator supports for generated
code.

The option ``EXTRA_OPTIONS`` may pass additional arguments when invoking the code generator (e.g. "-empty-string-as-null -optional std::shared_ptr")

.. _ObjectBox: https://objectbox.io

.. _ObjectBoxGenerator: https://github.com/objectbox/objectbox-generator



#]=======================================================================]

option(OBX_GENERATOR_ALLOW_FETCH "Opt-in automatic download and prepare for local execution" ON)

# Updated by maintainer to latest available version:
set(ObjectBoxGenerator_FETCH_VERSION 4.0.0-alpha1)
set(ObjectBoxGenerator_FETCH_DIR ${CMAKE_BINARY_DIR}/ObjectBoxGenerator-Fetch)
set(ObjectBoxGenerator_FETCH_BASEURL "https://github.com/objectbox/objectbox-generator/releases/download")
set(ObjectBoxGenerator_INSTALL_DIR ${CMAKE_BINARY_DIR}/ObjectBoxGenerator-Install)

include(FindPackageHandleStandardArgs)

# read version from objectbox-generator executable
function(_get_version)
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
  set(ObjectBoxGenerator_VERSION ${ObjectBoxGenerator_VERSION} PARENT_SCOPE)
endfunction()

# if already fetched/downloaded, path will available from CMake cache.
# (No need to pass in ObjectBoxGenerator_INSTALL_DIR)
find_program(ObjectBoxGenerator_EXECUTABLE objectbox-generator) 

if(ObjectBoxGenerator_EXECUTABLE)
  _get_version()
  if(${CMAKE_VERSION} VERSION_GREATER_EQUAL 3.19)
    find_package_check_version(${ObjectBoxGenerator_VERSION} ObjectBoxGenerator_VERSION_OK)
  else()
    if(ObjectBoxGenerator_FIND_VERSION)
      if(${ObjectBoxGenerator_VERSION} VERSION_GREATER_EQUAL ${ObjectBoxGenerator_FIND_VERSION}) 
        set(ObjectBoxGenerator_VERSION_OK TRUE)
      else()
        set(ObjectBoxGenerator_VERSION_OK FALSE)
      endif()
    endif()
  endif()
  if(NOT ${ObjectBoxGenerator_VERSION_OK}) 
    message(STATUS "ObjectBoxGenerator: Found version ${ObjectBoxGenerator_VERSION}, but which is not suitable with requested one ${ObjectBoxGenerator_FIND_VERSION}")
    set(ObjectBoxGenerator_FETCH_REQUIRED TRUE)
  endif()
else()
  set(ObjectBoxGenerator_FETCH_REQUIRED TRUE)
endif()
  
if(OBX_GENERATOR_ALLOW_FETCH AND ObjectBoxGenerator_FETCH_REQUIRED)
  message(STATUS "ObjectBox-Generator Fetch: Executable not found, attempting to download to build directory and prepare for execution (to disable behaviour set OBX_GENERATOR_ALLOW_FETCH to OFF)")
  set(_ObjectBoxGenerator_do_fetch TRUE)
  if (CMAKE_HOST_SYSTEM_NAME STREQUAL "Linux" AND CMAKE_HOST_SYSTEM_PROCESSOR STREQUAL "x86_64")
    set(ObjectBoxGenerator_FETCH_ARCH Linux)
  elseif (CMAKE_HOST_SYSTEM_NAME STREQUAL "Windows" AND (
            CMAKE_HOST_SYSTME_PROCESSOR STREQUAL "x86" OR
            CMAKE_HOST_SYSTEM_PROCESSOR STREQUAL "x86_64" OR 
            CMAKE_HOST_SYSTEM_PROCESSOR STREQUAL "AMD64"))
    set(ObjectBoxGenerator_FETCH_ARCH Windows)
  elseif (CMAKE_HOST_SYSTEM_NAME STREQUAL "Darwin")
    set(ObjectBoxGenerator_FETCH_ARCH macOS)
  else()
    message(WARNING "ObjectBoxGenerator Fetch aborted: unsupported platform (must be Linux/x86-64, Windows/x86-64 or macOS)")
    set(_ObjectBoxGenerator_do_fetch FALSE)
  endif()
  if(ObjectBoxGeenrator_FIND_VERSION_EXACT)
    message(WARNING "Exact not possible")
      set(_ObjectBoxGenerator_do_fetch FALSE)
  elseif(ObjectBoxGenerator_FIND_VERSION)
    if(${ObjectBoxGenerator_FIND_VERSION} VERSION_GREATER ${ObjectBoxGenerator_FETCH_VERSION})
      message(WARNING "Requested ObjectBox-Generator version '${ObjectBoxGenerator_FIND_VERSION}' not available. Latest version is ${ObjectBoxGenerator_FETCH_VERSION}")
      set(_ObjectBoxGenerator_do_fetch FALSE)
    endif()
  endif()
  if(_ObjectBoxGenerator_do_fetch)
    message(STATUS "ObjectBox-Generator Fetch: Downloading version ${ObjectBoxGenerator_FETCH_VERSION}")
    set(ObjectBoxGenerator_FETCH_FILE objectbox-generator-${ObjectBoxGenerator_FETCH_ARCH}.zip)
    set(ObjectBoxGenerator_FETCH_URL ${ObjectBoxGenerator_FETCH_BASEURL}/v${ObjectBoxGenerator_FETCH_VERSION}/${ObjectBoxGenerator_FETCH_FILE})
    set(ObjectBoxGenerator_FETCH_PATH ${ObjectBoxGenerator_FETCH_DIR}/${ObjectBoxGenerator_FETCH_FILE})
    set(ObjectBoxGenerator_UNPACK_FILE ${ObjectBoxGenerator_FETCH_DIR}/objectbox-generator)
    if (WIN32)
      set(ObjectBoxGenerator_UNPACK_FILE ${ObjectBoxGenerator_UNPACK_FILE}.exe)
    endif()
    message(STATUS "ObjectBox-Generator Fetch: Downloading archive from ${ObjectBoxGenerator_FETCH_URL} to ${ObjectBoxGenerator_FETCH_PATH}")
    file(DOWNLOAD "${ObjectBoxGenerator_FETCH_URL}" "${ObjectBoxGenerator_FETCH_PATH}"
          TLS_VERIFY ON 
          STATUS DownloadStatus
    )
    if (DownloadStatus EQUAL 0)
    message(STATUS "ObjectBox-Generator Fetch: Unpacking archive file ${ObjectBoxGenerator_FETCH_PATH} to ${ObjectBoxGenerator_UNPACK_FILE}")
    execute_process(COMMAND ${CMAKE_COMMAND} -E tar xf ${ObjectBoxGenerator_FETCH_PATH} WORKING_DIRECTORY ${ObjectBoxGenerator_FETCH_DIR})
    # 3.18: file(ARCHIVE_EXTRACT INPUT ${ObjectBoxGenerator_FETCH_PATH} DESTINATION ${ObjectBoxGenerator_FETCH_DIR} VERBOSE)
    message(STATUS "ObjectBox-Generator Fetch: Make file ${ObjectBoxGenerator_UNPACK_FILE} executable, and retry location..")
    file(COPY ${ObjectBoxGenerator_UNPACK_FILE} DESTINATION ${ObjectBoxGenerator_INSTALL_DIR} FILE_PERMISSIONS OWNER_READ OWNER_EXECUTE GROUP_READ GROUP_EXECUTE WORLD_READ WORLD_EXECUTE)
    # 3.19: file(CHMOD ${ObjectBoxGenerator_UNPACK_FILE} PERMISSIONS OWNER_READ OWNER_EXECUTE GROUP_READ GROUP_EXECUTE WORLD_READ WORLD_EXECUTE)
    # unset/find_program(..) fails with 3.20, so we redefine ObjectBoxGenerator_EXECUTABLE directly
    if(WIN32)
      set(ObjectBoxGenerator_EXECUTABLE ${ObjectBoxGenerator_INSTALL_DIR}/objectbox-generator.exe)
      set(ObjectBoxGenerator_EXECUTABLE ${ObjectBoxGenerator_INSTALL_DIR}/objectbox-generator.exe CACHE FILEPATH "" FORCE)
    else()
      set(ObjectBoxGenerator_EXECUTABLE ${ObjectBoxGenerator_INSTALL_DIR}/objectbox-generator)
      set(ObjectBoxGenerator_EXECUTABLE ${ObjectBoxGenerator_INSTALL_DIR}/objectbox-generator CACHE FILEPATH "" FORCE)
    endif()
    _get_version()
    else()
      message(WARNING "ObjectBox-Generator Fetch: Failed to download, reason: ${DownloadStatus}")
    endif()
  endif()
endif()

find_package_handle_standard_args(ObjectBoxGenerator 
  REQUIRED_VARS ObjectBoxGenerator_EXECUTABLE
  VERSION_VAR ObjectBoxGenerator_VERSION
)

if (NOT ObjectBoxGenerator_FOUND)
  return()
endif()

function (add_obx_schema)

  set(options INSOURCE)
  set(oneValueArgs TARGET;CXX_STANDARD)
  set(multiValueArgs SCHEMA_FILES;EXTRA_OPTIONS)
  cmake_parse_arguments(ARG "${options}" "${oneValueArgs}" "${multiValueArgs}" ${ARGN})

  if (ARG_INSOURCE)	
    set(base_dir ${CMAKE_CURRENT_SOURCE_DIR})
  else()
    set(base_dir ${CMAKE_CURRENT_BINARY_DIR}/ObjectBoxGenerator-output)
  endif()

  set(sources)

  set(lang -cpp)
  if(ARG_CXX_STANDARD)
    if(ARG_CXX_STANDARD EQUAL 11)
      set(lang -cpp${ARG_CXX_STANDARD})
    elseif(NOT ARG_CXX_STANDARD EQUAL 14)
      message(WARNING "ObjectBoxGenerator: CXX_STANDARD ${ARG_CXX_STANDARD} not supported, available are: 11 and 14. Defaults to 14.")
    endif()
  endif()

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
    # We explicitly do not add "objectbox-model.json" file here 
    # as it would otherwise be removed by a Makefile clean.
    # In addition, "objectbox-model.h" is also not mentioned as
    # it also fixes ninja build issues.
    # ObjectBox Model JSON file will be always maintained in current top-level source.
    add_custom_command(
      OUTPUT 
        ${cppfile} 
        ${hppfile} 
      COMMAND 
        ${ObjectBoxGenerator_EXECUTABLE} ARGS -out ${out_dir} -model ${CMAKE_CURRENT_SOURCE_DIR}/objectbox-model.json ${lang} ${ARG_EXTRA_OPTIONS} ${schema_filepath}
      DEPENDS 
        ${schema_filepath}
      USES_TERMINAL # Needed for ninja
    )
    list(APPEND sources ${cppfile} ${hppfile})
  endforeach()
    
  target_sources(${ARG_TARGET} PRIVATE ${sources}) 
  if (NOT ARG_INSOURCE)
    target_include_directories(objectbox-test PRIVATE ${base_dir})
  endif() 
endfunction()
