Limited C-API for FlatBuffers library and FlatBuffers compiler.

This library is linked by `internal/generator/flatbuffersc/fbsc.go` for
two purposes:

1. Parsing a schema file to flatbuffers serialized bytes block. This function 
   is used by Go (`flatbuffersc.ParseSchemaFile`) to build a tree of reflection 
   objects (passing the serialization to `reflection.GetRootAsSchema`). Modules under 
   `reflection` are auto-generated by `scripts/update-flatbuffersc-reflection.sh`.
2. Embedding the flatc compiler to be launched via `objectbox-generator FLATC <Flatc-Args>..`
