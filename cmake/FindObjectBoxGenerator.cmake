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

This function adds ObjectBox schema files (.fbs) to a C++ target which
implies a build task to auto-generate C++ source/header files 
from schema file(s) (with dependency rule tracking) and
adds them as sources to the target for compilation::

     add_obx_schema(
       TARGET <target>
       SCHEMA_FILES <schemafile>..
       [INSOURCE]
       [OUTPUT_DIR <path>]
       [OUTPUT_DIR_HEADERS <path>]
       [OUTPUT_DIR_MODEL_JSON <path>]
       [CXX_STANDARD 11|14]
       [EXTRA_OPTIONS <options>..]
     )
  
Note: The parameters ``TARGET`` and ``SCHEMA_FILES`` are required.

``TARGET`` specifies the CMake target to which the generated sources shall be assigned to.

``SCHEMA_FILES`` takes one or multiple ObjectBox schema file(s).
A schema file is the input for the ObjectBox Generator and defines classes and their members.
For details on the schema file please refer to the documentation.
Schema files have the pattern ``<name>.fbs``.
For each schema file, the generator creates a C++ source and header file
using the pattern ``<name>.obx.cpp`` and ``<name>.obx.hpp``, respectively.
(Also, two model files are generated: objectbox-model.h and a objectbox-model.json.)

The option ``INSOURCE`` tells the generator to place all generated files in the source tree (directory).
Note, that by default, the generator writes the generated C/C++ sources to the CMake build dir.
It's often preferable to use INSOURCE, as it can have several advantages:
  
* It makes the generated sources more "visible" to developers.
* It allows checking in generated sources to version control.
* It does not require a generator setup for consumers, e.g. after checkout.

``OUTPUT_DIR`` specifies the location for auto-generated files in the source tree
(default: current source directory).
For in-source (INSOURCE) builds, this affects all generated files.
For out-of-source builds, it only affects the ``objectbox-model.json`` file, because must be be kept in-source.
The given directory can be relative to current source directory or can be given as absolute path.

``OUTPUT_DIR_HEADERS`` sets the output directory for generated header files for INSOURCE builds.
It can be used alongside OUTPUT_DIR and then "overwrites" the directory for headers (only).

The option ``OUTPUT_DIR_MODEL_JSON`` specifies the location of the generated ``objectbox-model.json`` file.
It defaults to current source directory, or ``OUTPUT_DIR`` if it is given.
This generated file must be maintained under version source control
since it is essential maintain database schema changes over time.

Supply the option ``CXX_STANDARD`` to generate sources complying to a lower C++ standard, i.e. ``11`` for C++11.
By default, and when ``14`` is given, the generator creates sources compatible with C++14 and higher versions.

The option ``EXTRA_OPTIONS`` may pass additional arguments directly to the
code generator executable (e.g. "-empty-string-as-null -optional std::shared_ptr")

Out-of-source configuration notes:
per default, generated files (except the model JSON file) are written relative to the current binary (build) directory.
Generated headers and sources are written to the sub-directories ``ObjectBoxGenerator-include`` and
``ObjectBoxGenerator-src``, respectively.

.. _ObjectBox: https://objectbox.io

.. _ObjectBoxGenerator: https://github.com/objectbox/objectbox-generator



#]=======================================================================]

option(OBX_GENERATOR_ALLOW_FETCH "Opt-in automatic download and prepare for local execution" ON)

# Updated by maintainer to latest available version:
set(ObjectBoxGenerator_FETCH_VERSION 4.0.0-alpha2) # Note: must be safe for filename
# Using the version in the directories used for the generator executable to cleanly support multiple versions.
# This is e.g. relevant when updating to ensure fetching the new version.
set(ObjectBoxGenerator_FETCH_DIR ${CMAKE_BINARY_DIR}/ObjectBoxGenerator-Fetch/${ObjectBoxGenerator_FETCH_VERSION})
set(ObjectBoxGenerator_FETCH_BASEURL "https://github.com/objectbox/objectbox-generator/releases/download")
set(ObjectBoxGenerator_INSTALL_DIR ${CMAKE_BINARY_DIR}/ObjectBoxGenerator-Install/${ObjectBoxGenerator_FETCH_VERSION})

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

  # Search for "v1.2..." pattern in the output, e.g. v4.0.0-alpha0
  string(REGEX REPLACE ".* v\([0-9]*\.[0-9]*[^ \t\n\r]*\).*" "\\1" ObjectBoxGenerator_VERSION ${Output})
  message(STATUS "ObjectBox Generator version ${ObjectBoxGenerator_VERSION} found (${ObjectBoxGenerator_EXECUTABLE})")
  string(REGEX REPLACE "([0-9]+)\\..*" "\\1" ObjectBoxGenerator_VERSION_MAJOR "${ObjectBoxGenerator_VERSION}")
  string(REGEX REPLACE "[0-9]+\\.([0-9]+)\\..*" "\\1" ObjectBoxGenerator_VERSION_MINOR "${ObjectBoxGenerator_VERSION}")
  string(REGEX REPLACE "[0-9]+\\.[0-9]+\\.([0-9]+).*" "\\1" ObjectBoxGenerator_VERSION_PATCH "${ObjectBoxGenerator_VERSION}")

  set(ObjectBoxGenerator_VERSION ${ObjectBoxGenerator_VERSION} PARENT_SCOPE)
  set(ObjectBoxGenerator_VERSION_MAJOR ${ObjectBoxGenerator_VERSION_MAJOR} PARENT_SCOPE)
  set(ObjectBoxGenerator_VERSION_MINOR ${ObjectBoxGenerator_VERSION_MINOR} PARENT_SCOPE)
  set(ObjectBoxGenerator_VERSION_PATCH ${ObjectBoxGenerator_VERSION_PATCH} PARENT_SCOPE)
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
  set(oneValueArgs TARGET;OUTPUT_DIR;OUTPUT_DIR_HEADERS;OUTPUT_DIR_MODEL_JSON;CXX_STANDARD)
  set(multiValueArgs SCHEMA_FILES;EXTRA_OPTIONS)
  cmake_parse_arguments(ARG "${options}" "${oneValueArgs}" "${multiValueArgs}" ${ARGN})
  
  if(NOT ARG_TARGET)
    message(FATAL_ERROR "cmake_parse_arguments(): Missing target (Argument TARGET is empty or not set).")
  endif()
  if(NOT ARG_SCHEMA_FILES)
    message(FATAL_ERROR "cmake_parse_arguments(): Missing schema file(s) (Argument SCHEMA_FILES is empty or not set).")
  endif()

  # Prepare OBX_GEN_OUTPUT_DIR
  if (ARG_OUTPUT_DIR)
    if(IS_ABSOLUTE ${ARG_OUTPUT_DIR})
      set(OBX_GEN_OUTPUT_DIR ${ARG_OUTPUT_DIR})
    else()
      set(OBX_GEN_OUTPUT_DIR ${CMAKE_CURRENT_SOURCE_DIR}/${ARG_OUTPUT_DIR})
    endif()
    file(MAKE_DIRECTORY ${OBX_GEN_OUTPUT_DIR})
  endif()
  
  # Prepare OBX_GEN_OUTPUT_DIR_MODEL_JSON
  if (ARG_OUTPUT_DIR_MODEL_JSON)
    if(NOT IS_ABSOLUTE ${ARG_OUTPUT_DIR_MODEL_JSON})
      set(OBX_GEN_OUTPUT_DIR_MODEL_JSON ${CMAKE_CURRENT_SOURCE_DIR}/${ARG_OUTPUT_DIR_MODEL_JSON})
    else()
      set(OBX_GEN_OUTPUT_DIR_MODEL_JSON ${ARG_OUTPUT_DIR_MODEL_JSON})
    endif()
    file(MAKE_DIRECTORY ${OBX_GEN_OUTPUT_DIR_MODEL_JSON})
  else()
    if (OBX_GEN_OUTPUT_DIR)
      set(OBX_GEN_OUTPUT_DIR_MODEL_JSON ${OBX_GEN_OUTPUT_DIR})
    else()
      set(OBX_GEN_OUTPUT_DIR_MODEL_JSON ${CMAKE_CURRENT_SOURCE_DIR})
    endif()
  endif()

  # Prepare OBX_GEN_OUTPUT_DIR_SRC and OBX_GEN_OUTPUT_DIR_HEADERS
  if (ARG_INSOURCE)
    if (OBX_GEN_OUTPUT_DIR)
      set(OBX_GEN_OUTPUT_DIR_SRC ${OBX_GEN_OUTPUT_DIR})
    else()
      set(OBX_GEN_OUTPUT_DIR_SRC)
    endif()
    if (ARG_OUTPUT_DIR_HEADERS)
      if(IS_ABSOLUTE ${ARG_OUTPUT_DIR_HEADERS})
        set(OBX_GEN_OUTPUT_DIR_HEADERS ${ARG_OUTPUT_DIR_HEADERS})
      else()
        set(OBX_GEN_OUTPUT_DIR_HEADERS ${CMAKE_CURRENT_SOURCE_DIR}/${ARG_OUTPUT_DIR_HEADERS})
      endif()
      file(MAKE_DIRECTORY ${OBX_GEN_OUTPUT_DIR_HEADERS})
    else()
      set(OBX_GEN_OUTPUT_DIR_HEADERS)
    endif()
  else()
      if (ARG_OUTPUT_DIR_HEADERS)
        message(FATAL_ERROR "cmake_parse_arguments(): Option 'OUTPUT_DIR_HEADERS' only available for INSOURCE mode")
      endif()
    set(OBX_GEN_OUTPUT_DIR_SRC ${CMAKE_CURRENT_BINARY_DIR}/ObjectBoxGenerator-src)
    set(OBX_GEN_OUTPUT_DIR_HEADERS ${CMAKE_CURRENT_BINARY_DIR}/ObjectBoxGenerator-include)
  endif()

  set(sources)

  set(lang -cpp)
  if(ARG_CXX_STANDARD)
    if(ARG_CXX_STANDARD EQUAL 11)
      set(lang -cpp${ARG_CXX_STANDARD})
    elseif(NOT ARG_CXX_STANDARD EQUAL 14)
      message(FATAL_ERROR "cmake_parse_arguments(): CXX_STANDARD ${ARG_CXX_STANDARD} is not a known value. Set the it to 11 to generate C++11, otherwise it defaults to minimum level 14.")
    endif()
  endif()

  foreach(SCHEMA_FILE ${ARG_SCHEMA_FILES})
    
    # 3.20: cmake_path(ABSOLUTE_PATH SCHEMA_FILE OUTPUT_VARIABLE schema_filepath)
    if(IS_ABSOLUTE ${SCHEMA_FILE})
      set(schema_filepath ${SCHEMA_FILE})
    else()
      set(schema_filepath ${CMAKE_CURRENT_SOURCE_DIR}/${SCHEMA_FILE})
    endif()
    if(NOT EXISTS "${schema_filepath}")
      # Assuming it will be generated during build time, we do not produce an error, but just a status log:
      message(STATUS
              "cmake_parse_arguments(): Given schema file \"${schema_filepath}\" was not found at configuration time.")
    endif()
    
    string(REGEX REPLACE "\\.fbs$" "" basefile ${SCHEMA_FILE})
    string(REGEX REPLACE ".*/"     "" basefile ${basefile})
    
    if (ARG_INSOURCE)
      if (NOT OBX_GEN_OUTPUT_DIR_SRC) # no output directory is set, so we take the parent directory
        string(REGEX REPLACE "\.fbs$" ".obx.cpp" cppfile ${schema_filepath})
        string(REGEX REPLACE "\.fbs$" ".obx.hpp" hppfile ${schema_filepath})
        
        # 3.20: cmake_path(GET cppfile PARENT_PATH out_dir)
        string(REGEX REPLACE "/[^/]*$" "" parent_dir ${cppfile}) 
        set(obxGenOutOptions -out ${parent_dir})
      else() # output directory is set
        set(obxGenOutOptions -out ${OBX_GEN_OUTPUT_DIR_SRC})
        set(cppfile ${OBX_GEN_OUTPUT_DIR_SRC}/${basefile}.obx.cpp)
        set(hppfile ${OBX_GEN_OUTPUT_DIR_SRC}/${basefile}.obx.hpp)
      endif()
      if (OBX_GEN_OUTPUT_DIR_HEADERS)
        list(APPEND obxGenOutOptions -out-headers ${OBX_GEN_OUTPUT_DIR_HEADERS})
        set(hppfile ${OBX_GEN_OUTPUT_DIR_HEADERS}/${basefile}.obx.hpp)
      endif()
    else()
      set(cppfile ${OBX_GEN_OUTPUT_DIR_SRC}/${basefile}.obx.cpp)
      set(hppfile ${OBX_GEN_OUTPUT_DIR_HEADERS}/${basefile}.obx.hpp)
      set(obxGenOutOptions
        -out         ${OBX_GEN_OUTPUT_DIR_SRC}
        -out-headers ${OBX_GEN_OUTPUT_DIR_HEADERS}
      )
    endif()

    # We explicitly do not add "objectbox-model.json" file as output 
    # or byproduct here as it would otherwise be removed by a Makefile clean.
    # In addition, "objectbox-model.h" is also not mentioned as
    # it also fixes ninja build issues.
    # ObjectBox Model JSON file will be always maintained in OBX_GEN_OUTPUT_DIR_MODEL_JSON.
    add_custom_command(
      OUTPUT 
        ${cppfile} 
        ${hppfile} 
      COMMAND 
        ${ObjectBoxGenerator_EXECUTABLE}
      ARGS
          ${obxGenOutOptions}
          -model ${OBX_GEN_OUTPUT_DIR_MODEL_JSON}/objectbox-model.json 
          ${lang} 
          ${ARG_EXTRA_OPTIONS} 
          ${schema_filepath}
      DEPENDS 
        ${schema_filepath}
      USES_TERMINAL # Needed for ninja
    )
    list(APPEND sources ${cppfile} ${hppfile})
  endforeach()
    
  target_sources(${ARG_TARGET} PRIVATE ${sources}) 
  if (NOT ARG_INSOURCE)
    target_include_directories(${ARG_TARGET} PRIVATE ${OBX_GEN_OUTPUT_DIR_HEADERS})
  endif()
endfunction()
