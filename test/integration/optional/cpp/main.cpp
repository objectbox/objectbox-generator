#define CATCH_CONFIG_MAIN

#define OBX_CPP_FILE

#include <type_traits>

#include "c-ptr.obx.h"
#include "catch2/catch.hpp"
#include "shared/store-init.h"
#include "std-optional.obx.hpp"
#include "std-optional-as-null.obx.hpp"
#include "std-shared_ptr.obx.hpp"
#include "std-shared_ptr-as-null.obx.hpp"
#include "std-unique_ptr.obx.hpp"
#include "std-unique_ptr-as-null.obx.hpp"
#include "as-null.obx.hpp"

using namespace obx;

namespace {
template<typename Entity> void testOptionalValues() {
    static_assert(std::is_same<decltype(Entity::int_), std::optional<int>>::value, "must be std::optional");

    Store store = testStore(true, "c-cpp-tests-db");
    Box<Entity> box(store);

    // no values inserted -> no values loaded
    obx_id id = box.put(Entity());
    std::unique_ptr<Entity> read = box.get(id);
    REQUIRE(read);

    REQUIRE_FALSE(read->int_.has_value());
    REQUIRE_FALSE(read->int8.has_value());
    REQUIRE_FALSE(read->int16.has_value());
    REQUIRE_FALSE(read->int32.has_value());
    REQUIRE_FALSE(read->int64.has_value());
    REQUIRE_FALSE(read->uint.has_value());
    REQUIRE_FALSE(read->uint8.has_value());
    REQUIRE_FALSE(read->uint16.has_value());
    REQUIRE_FALSE(read->uint32.has_value());
    REQUIRE_FALSE(read->uint64.has_value());
    REQUIRE_FALSE(read->bool_.has_value());
    REQUIRE_FALSE(read->string.has_value());
    REQUIRE_FALSE(read->stringvector.has_value());
    REQUIRE_FALSE(read->byte.has_value());
    REQUIRE_FALSE(read->ubyte.has_value());
    REQUIRE_FALSE(read->bytevector.has_value());
    REQUIRE_FALSE(read->ubytevector.has_value());
    REQUIRE_FALSE(read->float32.has_value());
    REQUIRE_FALSE(read->float64.has_value());
    REQUIRE_FALSE(read->float_.has_value());
    REQUIRE_FALSE(read->floatvector.has_value());
    REQUIRE_FALSE(read->double_.has_value());
    REQUIRE_FALSE(read->relId.has_value());

    // values inserted -> values loaded
    std::unique_ptr<Entity> src = std::move(read);
    src->int_ = __LINE__;
    src->int8 = __LINE__;
    src->int16 = __LINE__;
    src->int32 = __LINE__;
    src->int64 = __LINE__;
    src->uint = __LINE__;
    src->uint8 = __LINE__;
    src->uint16 = __LINE__;
    src->uint32 = __LINE__;
    src->uint64 = __LINE__;
    src->bool_ = __LINE__;
    src->string = "foo";
    src->stringvector = std::vector<std::string>{"foo", "bar"};
    src->byte = __LINE__;
    src->ubyte = __LINE__;
    src->bytevector = std::vector<int8_t>{-13, 30};
    src->ubytevector = std::vector<uint8_t>{5, 6};
    src->float32 = __LINE__;
    src->float64 = __LINE__;
    src->float_ = __LINE__;
    src->floatvector = std::vector<float>{1.0f, 2.0f};
    src->double_ = __LINE__;
    src->relId = __LINE__;

    box.put(*src);
    read = box.get(id);
    REQUIRE(read);

    REQUIRE(read->int_.has_value());
    REQUIRE(read->int8.has_value());
    REQUIRE(read->int16.has_value());
    REQUIRE(read->int32.has_value());
    REQUIRE(read->int64.has_value());
    REQUIRE(read->uint.has_value());
    REQUIRE(read->uint8.has_value());
    REQUIRE(read->uint16.has_value());
    REQUIRE(read->uint32.has_value());
    REQUIRE(read->uint64.has_value());
    REQUIRE(read->bool_.has_value());
    REQUIRE(read->string.has_value());
    REQUIRE(read->stringvector.has_value());
    REQUIRE(read->byte.has_value());
    REQUIRE(read->ubyte.has_value());
    REQUIRE(read->bytevector.has_value());
    REQUIRE(read->ubytevector.has_value());
    REQUIRE(read->float32.has_value());
    REQUIRE(read->float64.has_value());
    REQUIRE(read->float_.has_value());
    REQUIRE(read->floatvector.has_value());
    REQUIRE(read->double_.has_value());
    REQUIRE(read->relId.has_value());

    REQUIRE(read->int_.value() == src->int_.value());
    REQUIRE(read->int8.value() == src->int8.value());
    REQUIRE(read->int16.value() == src->int16.value());
    REQUIRE(read->int32.value() == src->int32.value());
    REQUIRE(read->int64.value() == src->int64.value());
    REQUIRE(read->uint.value() == src->uint.value());
    REQUIRE(read->uint8.value() == src->uint8.value());
    REQUIRE(read->uint16.value() == src->uint16.value());
    REQUIRE(read->uint32.value() == src->uint32.value());
    REQUIRE(read->uint64.value() == src->uint64.value());
    REQUIRE(read->bool_.value() == src->bool_.value());
    REQUIRE(read->string.value() == src->string.value());
    REQUIRE(read->stringvector.value() == src->stringvector.value());
    REQUIRE(read->byte.value() == src->byte.value());
    REQUIRE(read->ubyte.value() == src->ubyte.value());
    REQUIRE(read->bytevector.value() == src->bytevector.value());
    REQUIRE(read->ubytevector.value() == src->ubytevector.value());
    REQUIRE(read->float32.value() == src->float32.value());
    REQUIRE(read->float64.value() == src->float64.value());
    REQUIRE(read->float_.value() == src->float_.value());
    REQUIRE(read->floatvector.value() == src->floatvector.value());
    REQUIRE(read->double_.value() == src->double_.value());
    REQUIRE(read->relId.value() == src->relId.value());
}

template<typename Entity> void testOptionalNull() {
    static_assert(std::is_same<decltype(Entity::int_), std::optional<int>>::value, "must be std::optional");

    Store store = testStore(true, "c-cpp-tests-db");
    Box<Entity> box(store);

    // no values inserted -> no values loaded
    obx_id id = box.put(Entity());

    // values inserted 
    Entity update;
    update.int_ = __LINE__;
    update.int8  = __LINE__;
    update.int16 = __LINE__;
    update.int32 = __LINE__;
    update.int64 = __LINE__;
    update.uint =  __LINE__;
    update.uint8 = __LINE__;
    update.uint16 = __LINE__;
    update.uint32 = __LINE__;
    update.uint64 = __LINE__;
    update.bool_ = __LINE__;
    update.string = "foo";
    update.stringvector = std::vector<std::string>{"foo", "bar"};
    update.byte = __LINE__;
    update.ubyte = __LINE__;
    update.bytevector = std::vector<int8_t>{-13, 30};
    update.ubytevector = std::vector<uint8_t>{5, 6};
    update.float32 = __LINE__;
    update.float64 = __LINE__;
    update.float_ = __LINE__;
    update.floatvector = std::vector<float>{-23.456f, 42.109f};
    update.double_ = __LINE__;
    update.relId = __LINE__;

    // overwritten by all fields null
    box.get(id, update);
 
    REQUIRE_FALSE(update.int_.has_value());
    REQUIRE_FALSE(update.int8.has_value());
    REQUIRE_FALSE(update.int16.has_value());
    REQUIRE_FALSE(update.int32.has_value());
    REQUIRE_FALSE(update.int64.has_value());
    REQUIRE_FALSE(update.uint.has_value());
    REQUIRE_FALSE(update.uint8.has_value());
    REQUIRE_FALSE(update.uint16.has_value());
    REQUIRE_FALSE(update.uint32.has_value());
    REQUIRE_FALSE(update.uint64.has_value());
    REQUIRE_FALSE(update.bool_.has_value());
    REQUIRE_FALSE(update.string.has_value());
    REQUIRE_FALSE(update.stringvector.has_value());
    REQUIRE_FALSE(update.byte.has_value());
    REQUIRE_FALSE(update.ubyte.has_value());
    REQUIRE_FALSE(update.bytevector.has_value());
    REQUIRE_FALSE(update.ubytevector.has_value());
    REQUIRE_FALSE(update.float32.has_value());
    REQUIRE_FALSE(update.float64.has_value());
    REQUIRE_FALSE(update.float_.has_value());
    REQUIRE_FALSE(update.floatvector.has_value());
    REQUIRE_FALSE(update.double_.has_value());
    REQUIRE_FALSE(update.relId.has_value());
}

template<typename Entity> void testOptionalEmptyNaN(bool asValue) {

    Store store = testStore(true, "c-cpp-tests-db");
    Box<Entity> box(store);

    obx_id id = box.put(Entity());
    std::unique_ptr<Entity> src = box.get(id);;
    REQUIRE(src);

    src->string  = "";
    src->float32 = std::nanf("1");
    src->float64 = std::nan("2");
    src->double_ = std::nan("-3");
    src->float_  = std::nanf("-4");
    src->int_    = 23;

    box.put(*src);

    // Update entity by empty-string/NaN 
    Entity updated;
    updated.string = "foo";
    updated.float32 = 1.0f;
    updated.float64 = 2.0;
    updated.double_ = 4.0;
    updated.float_  = 4.0f;
    updated.int_    = 42;

    box.get(id,updated);

    if(asValue) {
        REQUIRE(updated.string.has_value());
        REQUIRE(updated.float32.has_value());
        REQUIRE(updated.float64.has_value());
        REQUIRE(updated.double_.has_value());
        REQUIRE(updated.float_.has_value());
    } else { // as-null
        REQUIRE_FALSE(updated.string.has_value());
        REQUIRE_FALSE(updated.float32.has_value());
        REQUIRE_FALSE(updated.float64.has_value());
        REQUIRE_FALSE(updated.double_.has_value());
        REQUIRE_FALSE(updated.float_.has_value());
    }

    // Others always update as value:
    REQUIRE(updated.int_.has_value());
    REQUIRE(*updated.int_ == 23);
}

} // namespace

TEST_CASE("Optional") {
    SECTION("Values") {
        testOptionalValues<Optional>();
    }
    SECTION("Null") {
        testOptionalNull<Optional>();
    }
    SECTION("Empty/NaN as value") {
        testOptionalEmptyNaN<Optional>(true);
    }
}

TEST_CASE("OptionalAsNull") {
    SECTION("Values") {
        testOptionalValues<OptionalAsNull>();
    }
    SECTION("Nulls") {
        testOptionalNull<OptionalAsNull>();
    }
    SECTION("Empty/NaN as null") {
        testOptionalEmptyNaN<OptionalAsNull>(false);
    }
}

namespace {
template <typename Entity>
void testPtrValues() {
    Store store = testStore(true, "c-cpp-tests-db");
    Box<Entity> box(store);

    // no values inserted -> no values loaded
    obx_id id = box.put(Entity());
    std::unique_ptr<Entity> read = box.get(id);
    REQUIRE(read);

    REQUIRE_FALSE(read->int_.operator bool());
    REQUIRE_FALSE(read->int8.operator bool());
    REQUIRE_FALSE(read->int16.operator bool());
    REQUIRE_FALSE(read->int32.operator bool());
    REQUIRE_FALSE(read->int64.operator bool());
    REQUIRE_FALSE(read->uint.operator bool());
    REQUIRE_FALSE(read->uint8.operator bool());
    REQUIRE_FALSE(read->uint16.operator bool());
    REQUIRE_FALSE(read->uint32.operator bool());
    REQUIRE_FALSE(read->uint64.operator bool());
    REQUIRE_FALSE(read->bool_.operator bool());
    REQUIRE_FALSE(read->string.operator bool());
    REQUIRE_FALSE(read->stringvector.operator bool());
    REQUIRE_FALSE(read->byte.operator bool());
    REQUIRE_FALSE(read->ubyte.operator bool());
    REQUIRE_FALSE(read->bytevector.operator bool());
    REQUIRE_FALSE(read->ubytevector.operator bool());
    REQUIRE_FALSE(read->float32.operator bool());
    REQUIRE_FALSE(read->float64.operator bool());
    REQUIRE_FALSE(read->float_.operator bool());
    REQUIRE_FALSE(read->floatvector.operator bool());
    REQUIRE_FALSE(read->double_.operator bool());
    REQUIRE_FALSE(read->relId.operator bool());

    // values inserted -> values loaded
    std::unique_ptr<Entity> src = std::move(read);
    src->int_.reset(new int32_t(__LINE__));
    src->int8.reset(new int8_t(1));
    src->int16.reset(new int16_t(__LINE__));
    src->int32.reset(new int32_t(__LINE__));
    src->int64.reset(new int64_t(__LINE__));
    src->uint.reset(new uint32_t(__LINE__));
    src->uint8.reset(new uint8_t(2));
    src->uint16.reset(new uint16_t(__LINE__));
    src->uint32.reset(new uint32_t(__LINE__));
    src->uint64.reset(new uint64_t(__LINE__));
    src->bool_.reset(new bool(true));
    src->string.reset(new std::string("foo"));
    src->stringvector.reset(new std::vector<std::string>{"foo", "bar"});
    src->byte.reset(new int8_t(3));
    src->ubyte.reset(new uint8_t(4));
    src->bytevector.reset(new std::vector<int8_t>{-13, 30});
    src->ubytevector.reset(new std::vector<uint8_t>{5, 6});
    src->float32.reset(new float(__LINE__));
    src->float64.reset(new double(__LINE__));
    src->float_.reset(new float(__LINE__));
    src->floatvector.reset(new std::vector<float>{-23.456f,42.109f});
    src->double_.reset(new double(__LINE__));
    src->relId.reset(new obx_id(__LINE__));

    box.put(*src);
    read = box.get(id);
    REQUIRE(read);

    REQUIRE(read->int_.operator bool());
    REQUIRE(read->int8.operator bool());
    REQUIRE(read->int16.operator bool());
    REQUIRE(read->int32.operator bool());
    REQUIRE(read->int64.operator bool());
    REQUIRE(read->uint.operator bool());
    REQUIRE(read->uint8.operator bool());
    REQUIRE(read->uint16.operator bool());
    REQUIRE(read->uint32.operator bool());
    REQUIRE(read->uint64.operator bool());
    REQUIRE(read->bool_.operator bool());
    REQUIRE(read->string.operator bool());
    REQUIRE(read->stringvector.operator bool());
    REQUIRE(read->byte.operator bool());
    REQUIRE(read->ubyte.operator bool());
    REQUIRE(read->bytevector.operator bool());
    REQUIRE(read->ubytevector.operator bool());
    REQUIRE(read->float32.operator bool());
    REQUIRE(read->float64.operator bool());
    REQUIRE(read->float_.operator bool());
    REQUIRE(read->floatvector.operator bool());
    REQUIRE(read->double_.operator bool());
    REQUIRE(read->relId.operator bool());

    REQUIRE(*read->int_ == *src->int_);
    REQUIRE(*read->int8 == *src->int8);
    REQUIRE(*read->int16 == *src->int16);
    REQUIRE(*read->int32 == *src->int32);
    REQUIRE(*read->int64 == *src->int64);
    REQUIRE(*read->uint == *src->uint);
    REQUIRE(*read->uint8 == *src->uint8);
    REQUIRE(*read->uint16 == *src->uint16);
    REQUIRE(*read->uint32 == *src->uint32);
    REQUIRE(*read->uint64 == *src->uint64);
    REQUIRE(*read->bool_ == *src->bool_);
    REQUIRE(*read->string == *src->string);
    REQUIRE(*read->stringvector == *src->stringvector);
    REQUIRE(*read->byte == *src->byte);
    REQUIRE(*read->ubyte == *src->ubyte);
    REQUIRE(*read->bytevector == *src->bytevector);
    REQUIRE(*read->ubytevector == *src->ubytevector);
    REQUIRE(*read->float32 == *src->float32);
    REQUIRE(*read->float64 == *src->float64);
    REQUIRE(*read->float_ == *src->float_);
    REQUIRE(*read->floatvector == *src->floatvector);
    REQUIRE(*read->double_ == *src->double_);
    REQUIRE(*read->relId == *src->relId);
}

template<typename Entity> void testPtrNull() {

    Store store = testStore(true, "c-cpp-tests-db");
    Box<Entity> box(store);

    // no values inserted -> no values loaded
    obx_id id = box.put(Entity());

    // Update object populated with value 
    Entity update;
    update.int_.reset(new int(__LINE__));
    update.int8.reset(new int8_t(1));
    update.int16.reset(new int16_t(__LINE__));
    update.int32.reset(new int32_t(__LINE__));
    update.int64.reset(new int64_t(__LINE__));
    update.uint.reset(new unsigned int(__LINE__));
    update.uint8.reset(new uint8_t(2));
    update.uint16.reset(new uint16_t(__LINE__));
    update.uint32.reset(new uint32_t(__LINE__));
    update.uint64.reset(new uint64_t(__LINE__));
    update.bool_ .reset(new bool(true));
    update.string.reset(new std::string("foo"));
    update.stringvector.reset(new std::vector<std::string>{"foo", "bar"});
    update.byte.reset(new int8_t(3));
    update.ubyte.reset(new uint8_t(4));
    update.bytevector.reset(new std::vector<int8_t>{-13, 30});
    update.ubytevector.reset(new std::vector<uint8_t>{5, 6});
    update.float32.reset(new float(__LINE__));
    update.float64.reset(new double(__LINE__));
    update.float_.reset(new float(__LINE__));
    update.floatvector.reset(new std::vector<float>{-23.456f, 42.109f});
    update.double_.reset(new double(__LINE__));
    update.relId.reset(new uint64_t(__LINE__));;

    // Updated by null values
    box.get(id, update);

    REQUIRE_FALSE(update.int_.operator bool()); 
    REQUIRE_FALSE(update.int8.operator bool());
    REQUIRE_FALSE(update.int16.operator bool());
    REQUIRE_FALSE(update.int32.operator bool());
    REQUIRE_FALSE(update.int64.operator bool());
    REQUIRE_FALSE(update.uint.operator bool());
    REQUIRE_FALSE(update.uint8.operator bool());
    REQUIRE_FALSE(update.uint16.operator bool());
    REQUIRE_FALSE(update.uint32.operator bool());
    REQUIRE_FALSE(update.uint64.operator bool());
    REQUIRE_FALSE(update.bool_.operator bool());
    REQUIRE_FALSE(update.string.operator bool());
    REQUIRE_FALSE(update.stringvector.operator bool());
    REQUIRE_FALSE(update.byte.operator bool());
    REQUIRE_FALSE(update.ubyte.operator bool());
    REQUIRE_FALSE(update.bytevector.operator bool());
    REQUIRE_FALSE(update.ubytevector.operator bool());
    REQUIRE_FALSE(update.float32.operator bool());
    REQUIRE_FALSE(update.float64.operator bool());
    REQUIRE_FALSE(update.float_.operator bool());
    REQUIRE_FALSE(update.floatvector.operator bool());
    REQUIRE_FALSE(update.double_.operator bool());
    REQUIRE_FALSE(update.relId.operator bool());
}

template <typename Entity>
void testPtrEmptyNaN(bool asValue) {
    // Check for optional pointer types with "..-as-null" flag.
    
    // Create an asNull Empty/Null object in store
    Store store = testStore(true, "c-cpp-tests-db");
    Box<Entity> box(store);

    obx_id id = box.put(Entity());
    std::unique_ptr<Entity> nulled = box.get(id);
    REQUIRE(nulled);

    // Reset to Empty/NaN
    nulled->string.reset( new std::string("") );
    nulled->float32.reset( new float( std::nanf("1") ) );
    nulled->float64.reset( new double( std::nan("2") ) );
    nulled->double_.reset( new double( std::nan("-3") ) );
    nulled->float_.reset( new float( std::nanf("-4") ) );
    // Reset one other field
    nulled->int_.reset( new int(23) );

    box.put(*nulled);

    // Update
    Entity update;

    // Set Empty/NaN values
    update.string.reset(new std::string("foo"));
    update.float32.reset(new float(1.0f));
    update.float64.reset(new double(2.0));
    update.double_.reset(new double(4.0));
    update.float_.reset(new float(4.0f));

    // Set other and new-field 
    update.int_.reset(new int(42));
    update.uint.reset(new unsigned int(311));

    // Update
    box.get(id,update);

    if(asValue) {
        REQUIRE(update.string.operator bool());
        REQUIRE(update.float32.operator bool());
        REQUIRE(update.float64.operator bool());
        REQUIRE(update.double_.operator bool());
        REQUIRE(update.float_.operator bool());    
    } else { // as-null
        REQUIRE_FALSE(update.string.operator bool());
        REQUIRE_FALSE(update.float32.operator bool());
        REQUIRE_FALSE(update.float64.operator bool());
        REQUIRE_FALSE(update.double_.operator bool());
        REQUIRE_FALSE(update.float_.operator bool());    
    }
    // Updated by optional value:
    REQUIRE(update.int_.operator bool());
    REQUIRE(*update.int_ == 23);
    // Overwritten by optional null value:
    REQUIRE_FALSE(update.uint.operator bool());
}   
}  // namespace

TEST_CASE("UniquePtr") {
    SECTION("Values") {
        static_assert(std::is_same<decltype(UniquePtr::int_), std::unique_ptr<int32_t>>::value, "must be std::unique_ptr");
        testPtrValues<UniquePtr>();
    }
    SECTION("Null") {
        testPtrNull<UniquePtr>();
    }
    SECTION("Empty/NaN as value") {
        testPtrEmptyNaN<UniquePtr>(true);
    }
}

TEST_CASE("UniquePtrAsNull") {
    SECTION("Values") {
        static_assert(std::is_same<decltype(UniquePtrAsNull::int_), std::unique_ptr<int32_t>>::value, "must be std::unique_ptr");
        testPtrValues<UniquePtrAsNull>();
    }
    SECTION("Null") {
        testPtrNull<UniquePtrAsNull>();
    }
    SECTION("Empty/NaN as null") {
        testPtrEmptyNaN<UniquePtrAsNull>(false);
    }
}

TEST_CASE("SharedPtr") {
    SECTION("Values") {
        static_assert(std::is_same<decltype(SharedPtr::int_), std::shared_ptr<int32_t>>::value, "must be std::shared_ptr");
        testPtrValues<SharedPtr>();
    }
    SECTION("Null") {
        testPtrNull<SharedPtr>();
    }
    SECTION("Empty/NaN as value") {
        testPtrEmptyNaN<SharedPtr>(true);
    }
}

TEST_CASE("SharedPtrAsNull") {
    SECTION("Values") {
        static_assert(std::is_same<decltype(SharedPtrAsNull::int_), std::shared_ptr<int32_t>>::value, "must be std::shared_ptr");
        testPtrValues<SharedPtrAsNull>();
    }
    SECTION("Null") {
        testPtrNull<SharedPtrAsNull>();
    }
    SECTION("EmptyNaN as null") {
        testPtrEmptyNaN<SharedPtrAsNull>(false);
    }
}

TEST_CASE("TypefulAsNull") {
    Store store = testStore(true, "c-cpp-tests-db");
    Box<AsNull> box(store);

    obx_id id = box.put(AsNull());

    std::unique_ptr<AsNull> src = box.get(id);
    REQUIRE(src);

    REQUIRE(src->bytevector.size() == 0);
    REQUIRE(src->stringvector.size() == 0);
    REQUIRE(src->floatvector.size() == 0);

    // Initialize strings as "empty", float/double as NaN
    src->string = "";
    src->float32 = std::nanf("1");
    src->float64 = std::nan("2");
    src->double_ = std::nan("-3");
    src->float_  = std::nanf("-4");
    src->int_    = 23;

    src->bytevector.push_back(1);
    src->bytevector.push_back(2);
    src->bytevector.push_back(3);

    src->stringvector.push_back("foo");
    src->stringvector.push_back("bar");
    src->stringvector.push_back("");

    src->floatvector.push_back(23.456f);
    src->floatvector.push_back(-42.109f);
    src->floatvector.push_back(1.234f);
    
    // Put Entity
    box.put(*src);
    std::unique_ptr<AsNull> read = box.get(id);
    REQUIRE(read);

    REQUIRE(read->bytevector.size() == 3);
    REQUIRE(read->stringvector.size() == 3);
    REQUIRE(read->floatvector.size() == 3);

    // Initialize "updated" entity with values ("non-empty" and valid floating-point numbers)
    AsNull updated;
    updated.string = "foo";
    updated.float32 = 1.0f;
    updated.float64 = 2.0;
    updated.double_ = 4.0;
    updated.float_  = 4.0f;
    updated.int_    = 42;
    updated.uint    = 311;
    updated.floatvector = { 1.234f, 2.345f };

    // Update 
    box.get(id,updated);

    // Empty/NaN fields update fields to null values
    REQUIRE(updated.string.empty());
    REQUIRE(updated.float32 == 0.0f);
    REQUIRE(updated.float64 == 0.0);
    REQUIRE(updated.double_ == 0.0);
    REQUIRE(updated.float_  == 0.0f);

    // Two others are set to value
    REQUIRE(updated.int_ == 23);
    REQUIRE(updated.uint == 0);

    REQUIRE(updated.bytevector.size() == 3);
    REQUIRE(updated.stringvector.size() == 3);
    REQUIRE(updated.floatvector.size() == 3);
}

namespace {
template <typename T>
void cSetOptionalField(T*& outField, T value) {
    outField = (T*) malloc(sizeof(value));
    *outField = value;
}

template <typename T>
void cSetOptionalFieldVector(T*& outField, size_t& outLen, std::vector<T> value) {
    outLen = value.size();
    outField = (T*) malloc(sizeof(T*) * outLen);
    for (int i = 0; i < outLen; i++) {
        outField[i] = value[i];
    }
}
template <typename T>
void cSetOptionalFieldVector(T*& outField, size_t& outLen, std::vector<std::string> value) {
    outLen = value.size();
    outField = (T*) malloc(sizeof(T*) * outLen);
    for (int i = 0; i < outLen; i++) {
        outField[i] = (char*) malloc(sizeof(char) * (value[i].size() + 1));
        strcpy(outField[i], value[i].c_str());
    }
}
}  // namespace

TEST_CASE("c") {
    Store store = testStore(true, "c-cpp-tests-db");
    OBX_box* box = obx_box(store.cPtr(), PlainCPtr_ENTITY_ID);

    // no values inserted -> no values loaded
    PlainCPtr object{};
    obx_id id = PlainCPtr_put(box, &object);
    PlainCPtr* read = PlainCPtr_get(box, id);
    REQUIRE(read);

    REQUIRE(read->int_ == nullptr);
    REQUIRE(read->int8 == nullptr);
    REQUIRE(read->int16 == nullptr);
    REQUIRE(read->int32 == nullptr);
    REQUIRE(read->int64 == nullptr);
    REQUIRE(read->uint == nullptr);
    REQUIRE(read->uint8 == nullptr);
    REQUIRE(read->uint16 == nullptr);
    REQUIRE(read->uint32 == nullptr);
    REQUIRE(read->uint64 == nullptr);
    REQUIRE(read->bool_ == nullptr);
    REQUIRE(read->string == nullptr);
    REQUIRE(read->stringvector == nullptr);
    REQUIRE(read->byte == nullptr);
    REQUIRE(read->ubyte == nullptr);
    REQUIRE(read->bytevector == nullptr);
    REQUIRE(read->ubytevector == nullptr);
    REQUIRE(read->float32 == nullptr);
    REQUIRE(read->float64 == nullptr);
    REQUIRE(read->float_ == nullptr);
    REQUIRE(read->floatvector == nullptr);
    REQUIRE(read->double_ == nullptr);
    REQUIRE(read->relId == nullptr);
    PlainCPtr_free(read);

    // values inserted -> values loaded
    PlainCPtr src = {0};
    cSetOptionalField(src.int_, int32_t(__LINE__));
    cSetOptionalField(src.int8, int8_t(__LINE__));
    cSetOptionalField(src.int16, int16_t(__LINE__));
    cSetOptionalField(src.int32, int32_t(__LINE__));
    cSetOptionalField(src.int64, int64_t(__LINE__));
    cSetOptionalField(src.uint, uint32_t(__LINE__));
    cSetOptionalField(src.uint8, uint8_t(__LINE__));
    cSetOptionalField(src.uint16, uint16_t(__LINE__));
    cSetOptionalField(src.uint32, uint32_t(__LINE__));
    cSetOptionalField(src.uint64, uint64_t(__LINE__));
    cSetOptionalField(src.bool_, bool(__LINE__));
    src.string = (char*) malloc(sizeof(char) * 4);
    strcpy(src.string, "foo");
    cSetOptionalFieldVector(src.stringvector, src.stringvector_len, std::vector<std::string>{"foo", "bar"});
    cSetOptionalField(src.byte, int8_t(__LINE__));
    cSetOptionalField(src.ubyte, uint8_t(__LINE__));
    cSetOptionalFieldVector(src.bytevector, src.bytevector_len, std::vector<int8_t>{-13, 30});
    cSetOptionalFieldVector(src.ubytevector, src.ubytevector_len, std::vector<uint8_t>{5, 6});
    cSetOptionalField(src.float32, float(__LINE__));
    cSetOptionalField(src.float64, double(__LINE__));
    cSetOptionalField(src.float_, float(__LINE__));
    cSetOptionalFieldVector(src.floatvector, src.floatvector_len, std::vector<float>{-23.456f, 42.109f});
    cSetOptionalField(src.double_, double(__LINE__));
    cSetOptionalField(src.relId, obx_id(__LINE__));

    id = PlainCPtr_put(box, &src);
    read = PlainCPtr_get(box, id);
    REQUIRE(read);

    REQUIRE(read->int_ != nullptr);
    REQUIRE(read->int8 != nullptr);
    REQUIRE(read->int16 != nullptr);
    REQUIRE(read->int32 != nullptr);
    REQUIRE(read->int64 != nullptr);
    REQUIRE(read->uint != nullptr);
    REQUIRE(read->uint8 != nullptr);
    REQUIRE(read->uint16 != nullptr);
    REQUIRE(read->uint32 != nullptr);
    REQUIRE(read->uint64 != nullptr);
    REQUIRE(read->bool_ != nullptr);
    REQUIRE(read->string != nullptr);
    REQUIRE(read->stringvector != nullptr);
    REQUIRE(read->byte != nullptr);
    REQUIRE(read->ubyte != nullptr);
    REQUIRE(read->bytevector != nullptr);
    REQUIRE(read->ubytevector != nullptr);
    REQUIRE(read->float32 != nullptr);
    REQUIRE(read->float64 != nullptr);
    REQUIRE(read->float_ != nullptr);
    REQUIRE(read->floatvector != nullptr);
    REQUIRE(read->double_ != nullptr);
    REQUIRE(read->relId != nullptr);

    REQUIRE(*read->int_ == *src.int_);
    REQUIRE(*read->int8 == *src.int8);
    REQUIRE(*read->int16 == *src.int16);
    REQUIRE(*read->int32 == *src.int32);
    REQUIRE(*read->int64 == *src.int64);
    REQUIRE(*read->uint == *src.uint);
    REQUIRE(*read->uint8 == *src.uint8);
    REQUIRE(*read->uint16 == *src.uint16);
    REQUIRE(*read->uint32 == *src.uint32);
    REQUIRE(*read->uint64 == *src.uint64);
    REQUIRE(*read->bool_ == *src.bool_);
    REQUIRE(*read->string == *src.string);
    REQUIRE(read->stringvector_len == src.stringvector_len);
    REQUIRE(read->stringvector[0] == std::string(src.stringvector[0]));
    REQUIRE(read->stringvector[1] == std::string(src.stringvector[1]));
    REQUIRE(*read->byte == *src.byte);
    REQUIRE(*read->ubyte == *src.ubyte);
    REQUIRE(read->bytevector_len == src.bytevector_len);
    REQUIRE(0 == memcmp(read->bytevector, src.bytevector, sizeof(*src.bytevector) * src.bytevector_len));
    REQUIRE(read->ubytevector_len == src.ubytevector_len);
    REQUIRE(0 == memcmp(read->ubytevector, src.ubytevector, sizeof(*src.ubytevector) * src.ubytevector_len));
    REQUIRE(*read->float32 == *src.float32);
    REQUIRE(*read->float64 == *src.float64);
    REQUIRE(*read->float_ == *src.float_);
    REQUIRE(*read->floatvector == *src.floatvector);
    REQUIRE(*read->double_ == *src.double_);
    REQUIRE(*read->relId == *src.relId);

    PlainCPtr_free_pointers(&src);
    PlainCPtr_free(read);
}