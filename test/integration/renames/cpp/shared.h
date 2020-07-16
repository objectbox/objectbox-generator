#pragma once

#include "catch2/catch.hpp"
#include "objectbox-model.h"
#include "schema-cpp.obx.h"

obx::Store testStore(bool removeBeforeOpening) {
    obx::Store::Options options(create_obx_model());
    const char* dbDir = std::getenv("dbDir");
    REQUIRE(dbDir != nullptr);
    options.directory = dbDir;
    if (removeBeforeOpening) obx_remove_db_files(options.directory.c_str());
    return obx::Store(options);
}
