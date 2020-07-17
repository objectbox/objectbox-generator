#pragma once

#include "objectbox-model.h"

obx::Store testStore(bool removeBeforeOpening) {
    obx::Store::Options options(create_obx_model());
    const char* dbDir = std::getenv("dbDir");
    if (!dbDir) throw std::invalid_argument("dbDir environment variable not given");
    REQUIRE(dbDir != nullptr);
    options.directory = dbDir;
    if (removeBeforeOpening) obx_remove_db_files(options.directory.c_str());
    return obx::Store(options);
}
