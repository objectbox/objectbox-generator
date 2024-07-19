# Variant A: Integrate into project  
# - Fetch 'cmake/FindObjectBoxGenerator.cmake' to local cmake folder 
# - Add cmake folder to CMAKE_MODULE_PATH, e.g.
#   set(CMAKE_MODULE_PATH ${CMAKE_CURRENT_LIST_DIR}/cmake)

cmake_minimum_required(VERSION 3.5)
project(using-cmake)

# Variant B: Get Find-Module via FetchContent
include(FetchContent)
FetchContent_Declare(
  ObjectBoxGenerator
  URL      ${CMAKE_CURRENT_LIST_DIR}/../../cmake
)
FetchContent_MakeAvailable(ObjectBoxGenerator)


# Uncomment the following line to disable fetching from CMake project config.
# set(OBX_GENERATOR_ALLOW_FETCH OFF CACHE BOOL "Disable from CMakeLists.txt file")

find_package(ObjectBoxGenerator 4.0.0 MODULE REQUIRED)

# Test: run command 
message(STATUS "Configure-time test: running ${ObjectBoxGenerator_EXECUTABLE} -version :")
execute_process(COMMAND ${CMAKE_COMMAND} -E env ${ObjectBoxGenerator_EXECUTABLE} -version)
