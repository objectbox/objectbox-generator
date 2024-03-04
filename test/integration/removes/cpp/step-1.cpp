#define CATCH_CONFIG_MAIN

#define OBX_CPP_FILE

#include "catch2/catch.hpp"
#include "objectbox.hpp"
#include "schema.obx.hpp"
#include "shared/store-init.h"
using namespace obx;

TEST_CASE("Prepare DB with old entities", "") {
    Store store = testStore(true);
    Box<EntityA> boxA(store);
    Box<EntityB> boxB(store);
    REQUIRE(boxB.put({.id = 0, .name = "bb"}) == 1);
    REQUIRE(boxA.put({.id = 0, .name = "aa", .relId = 1}) == 1);
    REQUIRE(boxA.count() == 1);
    REQUIRE(boxB.count() == 1);
}
