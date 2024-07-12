set(CMAKE_MODULE_PATH ${CMAKE_CURRENT_LIST_DIR}/../../cmake)

cmake_minimum_required(VERSION 3.5)
project(using-cmake)

# Uncomment the following line to disable fetching from CMake project config.
# set(OBX_GENERATOR_ALLOW_FETCH OFF CACHE BOOL "Disable from CMakeLists.txt file")

find_package(ObjectBoxGenerator 0.14 MODULE REQUIRED)

# Test: run command 
message(STATUS "Configure-time test: running ${ObjectBoxGenerator_EXECUTABLE} -version :")
execute_process(COMMAND ${CMAKE_COMMAND} -E env ${ObjectBoxGenerator_EXECUTABLE} -version)