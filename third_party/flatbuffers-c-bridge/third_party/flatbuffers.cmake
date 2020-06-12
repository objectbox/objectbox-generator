set(FLATBUFFERS_VERSION 1.12.0)
set(FLATBUFFERS_CHECKSUM ba0a75fd12dbef8f6557a74e611b7a3d0c5fe7bf) # `sha1sum *.zip` from github releases page

include(FetchContent)
FetchContent_Declare(
        flatbuffers-${FLATBUFFERS_VERSION}
        URL "https://github.com/google/flatbuffers/archive/v${FLATBUFFERS_VERSION}.zip"
        URL_HASH SHA1=${FLATBUFFERS_CHECKSUM}
)

FetchContent_MakeAvailable(flatbuffers-${FLATBUFFERS_VERSION})
