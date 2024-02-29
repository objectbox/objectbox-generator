#define CATCH_CONFIG_MAIN

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
    REQUIRE(read->double_.value() == src->double_.value());
    REQUIRE(read->relId.value() == src->relId.value());
}
} // namespace

TEST_CASE("Optional") {
    SECTION("Values") {
        testOptionalValues<Optional>();
    }
    SECTION("Empty/NaN as values") {
        Store store = testStore(true, "c-cpp-tests-db");
        Box<Optional> box(store);
        // no values inserted -> no values loaded
        obx_id id = box.put(Optional());

        std::unique_ptr<Optional> src = box.get(id);
        REQUIRE(src);

        src->string = "";
        src->float32 = std::nanf("1");
        src->float64 = std::nan("2");
        src->double_ = std::nan("-3");
        src->float_  = std::nanf("-4");

        box.put(*src);
        std::unique_ptr<Optional> read = box.get(id);
        REQUIRE(read);

        REQUIRE(read->string.has_value());
        REQUIRE(read->float32.has_value());
        REQUIRE(read->float64.has_value());
        REQUIRE(read->double_.has_value());
        REQUIRE(read->float_.has_value());
    }
}

TEST_CASE("OptionalAsNull") {
    SECTION("Values") {
        testOptionalValues<OptionalAsNull>();
    }
    SECTION("Empty/NaN as null") {
        Store store = testStore(true, "c-cpp-tests-db");
        Box<OptionalAsNull> box(store);
        // no values inserted -> no values loaded
        obx_id id = box.put(OptionalAsNull());
        std::unique_ptr<OptionalAsNull> src = box.get(id);;
        REQUIRE(src);

        src->string  = "";
        src->float32 = std::nanf("1");
        src->float64 = std::nan("2");
        src->double_ = std::nan("-3");
        src->float_  = std::nanf("-4");

        box.put(*src);
        std::unique_ptr<OptionalAsNull> read = box.get(id);
        REQUIRE(read);

        REQUIRE_FALSE(read->string.has_value());
        REQUIRE_FALSE(read->float32.has_value());
        REQUIRE_FALSE(read->float64.has_value());
        REQUIRE_FALSE(read->double_.has_value());
        REQUIRE_FALSE(read->float_.has_value());
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
    REQUIRE_FALSE(read->double_.operator bool());
    REQUIRE_FALSE(read->relId.operator bool());

    // values inserted -> values loaded
    std::unique_ptr<Entity> src = std::move(read);
    src->int_.reset(new int32_t(__LINE__));
    src->int8.reset(new int8_t(__LINE__));
    src->int16.reset(new int16_t(__LINE__));
    src->int32.reset(new int32_t(__LINE__));
    src->int64.reset(new int64_t(__LINE__));
    src->uint.reset(new uint32_t(__LINE__));
    src->uint8.reset(new uint8_t(__LINE__));
    src->uint16.reset(new uint16_t(__LINE__));
    src->uint32.reset(new uint32_t(__LINE__));
    src->uint64.reset(new uint64_t(__LINE__));
    src->bool_.reset(new bool(__LINE__));
    src->string.reset(new std::string("foo"));
    src->stringvector.reset(new std::vector<std::string>{"foo", "bar"});
    src->byte.reset(new int8_t(__LINE__));
    src->ubyte.reset(new uint8_t(__LINE__));
    src->bytevector.reset(new std::vector<int8_t>{-13, 30});
    src->ubytevector.reset(new std::vector<uint8_t>{5, 6});
    src->float32.reset(new float(__LINE__));
    src->float64.reset(new double(__LINE__));
    src->float_.reset(new float(__LINE__));
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
    REQUIRE(*read->double_ == *src->double_);
    REQUIRE(*read->relId == *src->relId);
}
}  // namespace

TEST_CASE("UniquePtr") {
    
    SECTION("Values") {
        static_assert(std::is_same<decltype(UniquePtr::int_), std::unique_ptr<int32_t>>::value, "must be std::unique_ptr");
        testPtrValues<UniquePtr>();
    }

    SECTION("Empty/NaN") {
        Store store = testStore(true, "c-cpp-tests-db");
        Box<UniquePtr> box(store);

        // no values inserted -> no values loaded
        obx_id id = box.put(UniquePtr());
        std::unique_ptr<UniquePtr> read = box.get(id);
        REQUIRE(read);
        REQUIRE_FALSE(read->string.operator bool());
        REQUIRE_FALSE(read->float32.operator bool());
        REQUIRE_FALSE(read->float64.operator bool());
        REQUIRE_FALSE(read->float_.operator bool());
        REQUIRE_FALSE(read->double_.operator bool());
        std::unique_ptr<UniquePtr> src = std::move(read);

        src->string.reset(new std::string(""));
        src->float32.reset(new float(std::nanf("1")));
        src->float64.reset(new double(std::nan("2")));
        src->float_.reset(new float(std::nanf("3")));
        src->double_.reset(new double(std::nan("4")));

        box.put(*src);
        read = box.get(id);
        REQUIRE(read);

        REQUIRE(read->string.operator bool());
        REQUIRE(read->float32.operator bool());
        REQUIRE(read->float64.operator bool());
        REQUIRE(read->float_.operator bool());
        REQUIRE(read->double_.operator bool());
    }
}

TEST_CASE("UniquePtrAsNull") {

    SECTION("Values") {
        static_assert(std::is_same<decltype(UniquePtrAsNull::int_), std::unique_ptr<int32_t>>::value, "must be std::unique_ptr");
        testPtrValues<UniquePtrAsNull>();
    }

    SECTION("Empty/NaN") {
        Store store = testStore(true, "c-cpp-tests-db");
        Box<UniquePtrAsNull> box(store);

        // no values inserted -> no values loaded
        obx_id id = box.put(UniquePtrAsNull());
        std::unique_ptr<UniquePtrAsNull> read = box.get(id);
        REQUIRE(read);
        REQUIRE_FALSE(read->string.operator bool());
        REQUIRE_FALSE(read->float32.operator bool());
        REQUIRE_FALSE(read->float64.operator bool());
        REQUIRE_FALSE(read->float_.operator bool());
        REQUIRE_FALSE(read->double_.operator bool());
        std::unique_ptr<UniquePtrAsNull> src = std::move(read);

        src->string.reset(new std::string(""));
        src->float32.reset(new float(std::nanf("1")));
        src->float64.reset(new double(std::nan("2")));
        src->float_.reset(new float(std::nanf("3")));
        src->double_.reset(new double(std::nan("4")));

        box.put(*src);
        read = box.get(id);
        REQUIRE(read);

        REQUIRE_FALSE(read->string.operator bool());
        REQUIRE_FALSE(read->float32.operator bool());
        REQUIRE_FALSE(read->float64.operator bool());
        REQUIRE_FALSE(read->float_.operator bool());
        REQUIRE_FALSE(read->double_.operator bool());
    }

}

TEST_CASE("SharedPtr") {

    SECTION("Values") {
        static_assert(std::is_same<decltype(SharedPtr::int_), std::shared_ptr<int32_t>>::value, "must be std::shared_ptr");
        testPtrValues<SharedPtr>();
    }

    SECTION("Empty/NaN") {
        Store store = testStore(true, "c-cpp-tests-db");
        Box<SharedPtr> box(store);

        // no values inserted -> no values loaded
        obx_id id = box.put(SharedPtr());
        std::unique_ptr<SharedPtr> read = box.get(id);
        REQUIRE(read);
        REQUIRE_FALSE(read->string.operator bool());
        REQUIRE_FALSE(read->float32.operator bool());
        REQUIRE_FALSE(read->float64.operator bool());
        REQUIRE_FALSE(read->float_.operator bool());
        REQUIRE_FALSE(read->double_.operator bool());
        std::unique_ptr<SharedPtr> src = std::move(read);

        src->string.reset(new std::string(""));
        src->float32.reset(new float(std::nanf("1")));
        src->float64.reset(new double(std::nan("2")));
        src->float_.reset(new float(std::nanf("3")));
        src->double_.reset(new double(std::nan("4")));

        box.put(*src);
        read = box.get(id);
        REQUIRE(read);

        REQUIRE(read->string.operator bool());
        REQUIRE(read->float32.operator bool());
        REQUIRE(read->float64.operator bool());
        REQUIRE(read->float_.operator bool());
        REQUIRE(read->double_.operator bool());
    }
}

TEST_CASE("SharedPtrAsNull") {

    SECTION("Values") {
        static_assert(std::is_same<decltype(SharedPtrAsNull::int_), std::shared_ptr<int32_t>>::value, "must be std::shared_ptr");
        testPtrValues<SharedPtrAsNull>();
    }

    SECTION("Empty/NaN") {
        Store store = testStore(true, "c-cpp-tests-db");
        Box<SharedPtrAsNull> box(store);

        // no values inserted -> no values loaded
        obx_id id = box.put(SharedPtrAsNull());
        std::unique_ptr<SharedPtrAsNull> read = box.get(id);
        REQUIRE(read);
        REQUIRE_FALSE(read->string.operator bool());
        REQUIRE_FALSE(read->float32.operator bool());
        REQUIRE_FALSE(read->float64.operator bool());
        REQUIRE_FALSE(read->float_.operator bool());
        REQUIRE_FALSE(read->double_.operator bool());
        std::unique_ptr<SharedPtrAsNull> src = std::move(read);

        src->string.reset(new std::string(""));
        src->float32.reset(new float(std::nanf("1")));
        src->float64.reset(new double(std::nan("2")));
        src->float_.reset(new float(std::nanf("3")));
        src->double_.reset(new double(std::nan("4")));

        box.put(*src);
        read = box.get(id);
        REQUIRE(read);

        REQUIRE_FALSE(read->string.operator bool());
        REQUIRE_FALSE(read->float32.operator bool());
        REQUIRE_FALSE(read->float64.operator bool());
        REQUIRE_FALSE(read->float_.operator bool());
        REQUIRE_FALSE(read->double_.operator bool());
    }
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
    REQUIRE(*read->double_ == *src.double_);
    REQUIRE(*read->relId == *src.relId);

    PlainCPtr_free_pointers(&src);
    PlainCPtr_free(read);
}