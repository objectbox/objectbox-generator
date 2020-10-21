#define CATCH_CONFIG_MAIN

#include "catch2/catch.hpp"
#include "objectbox-cpp.h"
#include "objectbox-model.h"
#include "synced.obx.hpp"

using namespace obx;

namespace {
Store testStore() {
    Store::Options options(create_obx_model());
    options.directory = "c-cpp-tests-db";
    obx_remove_db_files(options.directory.c_str());
    return Store(options);
}
}  // namespace

TEST_CASE("Store can start", "") {
    Store store = testStore();
    // Nothing to check right now, we don't have a server available.
}
