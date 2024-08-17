#define CATCH_CONFIG_MAIN

#define OBX_CPP_FILE

#include "annotated.obx.hpp"
#include "catch2/catch.hpp"
#include "objectbox.hpp"
#include "objectbox-model.h"
#include "schema.obx.hpp"
#include "shared/store-init.h"

using namespace obx;

TEST_CASE("CRUD", "") {
    Store store = testStore(true,  "c-cpp-tests-db");

    // Box<Typeful_> box(store); // one way
    auto box = store.box<Typeful>();  // another way to get the same box

    REQUIRE(box.count() == 0);
    REQUIRE(box.cPtr() != nullptr);

    // put(`const Typeful&`) must be accepted for insert and update, just
    // returning an ID
    REQUIRE(box.put({.id = 0, .int_ = 11}) == 1);
    Typeful obj;
    obj.id = 1;
    obj.int_ = 99;
    obj.floatvector = std::vector<float> {-23.456f, 42.109f};
    REQUIRE(box.put(obj) == 1); 

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
    REQUIRE(read->floatvector.size() == 2);
    REQUIRE(read->floatvector[0] == -23.456f);
    REQUIRE(read->floatvector[1] ==  42.109f);

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

    REQUIRE(box.put({.id = 0, .int_ = 23}) == 3);  

    std::unique_ptr<Typeful> nullVectors = box.get(3);
    REQUIRE(nullVectors->floatvector.size() == 0);
}

TEST_CASE("Self assigned IDs", "") {
    Store store = testStore(true,  "c-cpp-tests-db");
    auto box = store.box<Annotated>();

    REQUIRE(box.isEmpty());

	// can't use brace initializer here on some older compilers ¯\_(ツ)_/¯
	Annotated item{};
    item.identifier = 0;
	item.time = 11;
    REQUIRE(box.put(item) == 1);

    item.identifier = 25;
    item.time = 99;
    REQUIRE(box.put(item) == 25);

    std::unique_ptr<Annotated> read = box.get(1);
    REQUIRE(read != nullptr);
    REQUIRE(read->time == 11);

    read = box.get(25);
    REQUIRE(read != nullptr);
    REQUIRE(read->time == 99);
}

TEST_CASE("update-entity", "") {
    Store store = testStore(true, "c-cpp-tests-db");
    auto box = store.box<Typeful>();
    Typeful myObj1{.int_ = 23};
    obx_id id = box.put(myObj1);

    Typeful myObj2{.int_ = 42, .string = "foobar"};
    box.get(id, myObj2);
    REQUIRE(myObj2.int_ == 23);
    REQUIRE(myObj2.string.empty());
}