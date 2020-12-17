#pragma once

#include "objectbox.hpp"
#include "objectbox-model.h"

obx::Store testStore(bool removeBeforeOpening, const char* dbDir = nullptr) {
    if (!dbDir) {
    	dbDir = std::getenv("dbDir");
    	if (!dbDir) throw std::invalid_argument("dbDir environment variable not given");
    }
    if (removeBeforeOpening) obx_remove_db_files(dbDir);
    obx::Store::Options options(create_obx_model());
    options.directory(dbDir);
    return obx::Store(options);
}
