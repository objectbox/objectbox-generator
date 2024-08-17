# CMake script to test cmake find module with different CMake generators/toolchains.
# 
#   Usage: cmake
#            [ -DPROJECTS="proj1;proj2;..." ]
#            [ -DSINGLE="gen1;gen2;..." ] 
#            [ -DMULTI="gen1;gen2;..." ]
#            [ -DAUTO=TRUE ]  
#            [ -DNO_PREBUILT=TRUE ]
#            -P test-cmake-generators.cmake
#
# PROJECTS specifies the set of projects; defaults to the complete set of CMake Test Projects.
# MULTI specifies Multi-Config Generators such as "Xcode". (Release and Debug builds will be built)
# SINGLE specifies Single-Config Generators. 
# If AUTO is set, Ninja/Make tools will be auto-detected and single-/multi-config generators are appended automatically.
# If NO_PREBUILT is set, a potential pre-built ObjectBox-Generator at top-level will not be used but always downloaded.
#
# Common to all CMake test projects:
# - Boilerplate code for fetching objectbox but using current source-tree find module available from "common.cmake".
# - Check for DO_INSOURCE to run add_obx_schema with INSOURCE variant.

# Declare CMake test project directories which should be tested.
set(PROJECTS cpp-flat;cpp-multiple-targets;cpp-tree;cpp-tree-multiple-targets;cpp-multiple-schema-dirs CACHE STRING "Test projects")
if(NO_PREBUILT)
  set(extraFlags "-DNO_PREBUILT=1")
endif()

# overload execute_process: Exit on cperror code.
function(execute_process)
  _execute_process(${ARGV} RESULT_VARIABLE error)
  if(error)
    message(FATAL_ERROR "")
  endif()
endfunction()

# configureAndBuild all PROJECTS using two variants: default and insource.
# If multi is true, build Release and Debug configurations.
# generator is passed as is to cmake via -G .
function(configureAndBuild multi generator)
    string(REPLACE " " "_" generatorLabel ${generator})
    foreach(project ${PROJECTS})
        set(srcdir ${CMAKE_CURRENT_LIST_DIR}/${project})
        foreach(insource TRUE;FALSE)
            if (insource)
                set(variant "insource")
                set(configureFlags "-DDO_INSOURCE=TRUE")
            else()
                set(variant "default")
                set(configureFlags)
            endif()
            set(builddir ${CMAKE_CURRENT_LIST_DIR}/build/${generatorLabel}/${project}/${variant})
            # Remove all auto-generated files from sources
            file(GLOB_RECURSE auto_generated "*/objectbox-model.*" "*/*.obx.*")
            list(FILTER auto_generated EXCLUDE REGEX "/build/")
            if (auto_generated)
                file(REMOVE ${auto_generated})
            endif()
            message(STATUS "-------------------------------------------------------")
            message(STATUS "**** Test Generator: ${generator} with test project '${project}' and variant '${variant}'.")
            message(STATUS "-------------------------------------------------------")
            execute_process(COMMAND ${CMAKE_COMMAND} -S ${srcdir} -B ${builddir} -G ${generator} ${configureFlags} ${extraFlags})
            if (multi)
                execute_process(COMMAND ${CMAKE_COMMAND} --build ${builddir} --config Release)
                execute_process(COMMAND ${CMAKE_COMMAND} --build ${builddir} --config Debug)
            else()
                execute_process(COMMAND ${CMAKE_COMMAND} --build ${builddir})
            endif()
        endforeach()
    endforeach()
endfunction()

if (AUTO)
    # Auto-detect Generators on Path
    find_program(NINJA ninja)
    if(NINJA)
      list(APPEND SINGLE "Ninja")
      list(APPEND MULTI "Ninja Multi-Config")
    endif()
    find_program(MAKE make)
    if(MAKE)
        list(APPEND SINGLE "Unix Makefiles")
    endif()
endif()

message(STATUS "*******************************************************
Test CMake Generators

Control-Script:     ${CMAKE_CURRENT_LIST_FILE}
Projects:           ${PROJECTS}
Extra flags:        ${extraFlags}
Multi-Generators:   ${MULTI}
Single-Generrators: ${SINGLE}
")


# Run Multi-Config Generators
foreach(generator ${MULTI})
    configureAndBuild(TRUE ${generator})
endforeach()

# Run Single-Config Generators
foreach(generator ${SINGLE})
    configureAndBuild(FALSE ${generator})
endforeach()

message(STATUS "Test CMake Generators: Success ")