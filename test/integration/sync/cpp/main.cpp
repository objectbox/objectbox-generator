#define CATCH_CONFIG_MAIN

#define OBX_CPP_FILE

#include "catch2/catch.hpp"
#include "objectbox.hpp"
#include "objectbox-model.h"
#include "synced.obx.hpp"
#include "shared/store-init.h"

using namespace obx;

TEST_CASE("CRUD", "") {
    Store store = testStore(true,  "c-cpp-tests-db");

    // Nothing to check right now, we don't have a server available.
}
