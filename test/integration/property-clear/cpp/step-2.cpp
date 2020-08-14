#define CATCH_CONFIG_MAIN

#include "catch2/catch.hpp"
#include "objectbox-cpp.h"
#include "schema.obx.hpp"
#include "shared/store-init.h"
using namespace obx;

TEST_CASE("Check DB after change", "") {
    Store store = testStore(false);
    Box<EntityName> box(store);  // one way
    REQUIRE(box.count() == 2);
    // "value" property is now empty (reset)
    REQUIRE(box.get(1)->value == 0);
    REQUIRE(box.get(2)->value == 0);
}
