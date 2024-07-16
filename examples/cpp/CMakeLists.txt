cmake_minimum_required(VERSION 3.14)

project("objectbox-test" CXX)

include(FetchContent)
FetchContent_Declare(
    objectbox
    GIT_REPOSITORY https://github.com/objectbox/objectbox-c.git
    GIT_TAG        v0.21.0
)
FetchContent_MakeAvailable(objectbox)

FetchContent_Declare(
    ObjectBoxGenerator
    URL ${CMAKE_CURRENT_LIST_DIR}/../../cmake
)
FetchContent_MakeAvailable(ObjectBoxGenerator)
find_package(ObjectBoxGenerator 0.14 REQUIRED)


add_executable("objectbox-test"
    "src/main.cpp"
)
target_compile_features("objectbox-test" PRIVATE cxx_std_14)
target_link_libraries("objectbox-test" objectbox)

set_target_properties("objectbox-test" PROPERTIES
    CXX_STANDARD 11
    CXX_STANDARD_REQUIRED YES
)

add_schema_file("objectbox-test" "src/schema/person.fbs")
add_schema_file("objectbox-test" "src/schema/task.fbs")
add_schema_file("objectbox-test" "src/schema/monster.fbs")
add_schema_file("objectbox-test" "src/schema/another_monster.fbs")

generate_schema_bindings("objectbox-test" "src/schema")
