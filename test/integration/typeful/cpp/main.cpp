#define CATCH_CONFIG_MAIN

#include "catch2/catch.hpp"
#include "objectbox-cpp.h"

#include "annotated-cpp.obx.h"
#include "objectbox-model.h"
#include "schema-cpp.obx.h"

using namespace obx;

namespace {
Store testStore() {
  Store::Options options(create_obx_model());
  options.directory = "c-cpp-tests-db";
  obx_remove_db_files(options.directory.c_str());
  return Store(options);
}
} // namespace

TEST_CASE("CRUD", "") {
  Store store = testStore();

  // Box<Typeful_> box(store); // one way
  auto box = store.box<Typeful>(); // another way to get the same box

  REQUIRE(box.count() == 0);
  REQUIRE(box.cPtr() != nullptr);

  // put(`const Typeful&`) must be accepted for insert and update, just
  // returning an ID
  REQUIRE(box.put({.id = 0, .int_ = 11}) == 1);
  REQUIRE(box.put({.id = 1, .int_ = 99}) == 1); // NOTE: .int_ is set to 0 now

  // put(`Typeful&`) must set the ID on the given object
  Typeful object2 = {};
  object2.string = "foo";
  REQUIRE(box.put(object2) == 2);
  REQUIRE(object2.id == 2);

  // pointer returning get()
  REQUIRE(box.get(3) == nullptr);
  std::unique_ptr<Typeful> read = box.get(1);
  REQUIRE(read != nullptr);
  REQUIRE(read->id == 1);
  REQUIRE(read->int_ == 99);

  // out-param get()
  REQUIRE(!box.get(3, *read));
  REQUIRE(box.get(2, *read));
  REQUIRE(read->id == 2);
  REQUIRE(read->int_ == 0);
  REQUIRE(read->int8 == 0);
  REQUIRE(read->string == "foo");

  REQUIRE(box.count() == 2);
  REQUIRE(box.count(1) == 1);
  REQUIRE(box.count(5) == 2);
  REQUIRE(!box.isEmpty());

  REQUIRE(!box.contains(3));
  REQUIRE(box.contains(1));
  REQUIRE(box.contains(1));

  REQUIRE(!box.remove(3));
  REQUIRE(box.remove(1));
  REQUIRE(box.count() == 1);
  REQUIRE(!box.contains(1));
  REQUIRE(box.contains(2));

  REQUIRE(box.removeAll() == 1);
  REQUIRE(!box.contains(2));
  REQUIRE(box.count() == 0);
  REQUIRE(box.isEmpty());
}
