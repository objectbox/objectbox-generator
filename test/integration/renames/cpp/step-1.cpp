#define CATCH_CONFIG_MAIN

#include "catch2/catch.hpp"
#include "objectbox-cpp.h"
#include "shared.h"
using namespace obx;

TEST_CASE("Prepare DB with old names", "") {
    Store store = testStore(true);
    Box<OldEntityName> box(store);  // one way
    REQUIRE(box.put({.id = 0, .oldPropertyName = 11}) == 1);
    REQUIRE(box.put({.id = 0, .oldPropertyName = 22}) == 2);
    REQUIRE(box.count() == 2);
}
