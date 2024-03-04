#define CATCH_CONFIG_MAIN

#define OBX_CPP_FILE

#include "catch2/catch.hpp"
#include "objectbox.hpp"
#include "schema.obx.hpp"
#include "shared/store-init.h"
using namespace obx;

TEST_CASE("Check DB after removal", "") {
    Store store = testStore(false);
    REQUIRE(Box<EntityB>(store).count() == 1);
}
