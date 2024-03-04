#define CATCH_CONFIG_MAIN

#define OBX_CPP_FILE

#include "catch2/catch.hpp"
#include "objectbox.hpp"
#include "schema.obx.hpp"
#include "shared/store-init.h"
using namespace obx;

TEST_CASE("Check DB with new names", "") {
    Store store = testStore(false);
    Box<NewEntityName> box(store);  // one way
    REQUIRE(box.count() == 2);
    REQUIRE(box.get(1)->newPropertyName == 11);
    REQUIRE(box.get(2)->newPropertyName == 22);
}
