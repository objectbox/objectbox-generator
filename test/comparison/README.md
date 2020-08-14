# Code comparison tests

This test suite is based on a premise that we expect an exact generator output (file contents).
Therefore, the test works as follows:
* read a test-case - all "source" files (e.g. *.fbs) in a single folder 
* clean-up any previously generated files
* execute a generator on the test-case (file by file)
* compare the generated files' contents to those stored as ".expected" with the same name
* [optional] compile the generated code

## Test-cases directory structure
* `<source-type>/<test-case>/*.<source-type>` are test case source files, 
    e.g. `fbs/typeful/schema.fbs`
    * negative tests: if the file has a fail infix in the name: `*.fail.<source-type>`, 
      it's considered a negative test case (the generation should fail)
* `<source-type>/<test-case>/<target-type>/*.<target-type-ext>.expected` are files expected to be generated
    e.g. `fbs/typeful/cpp/schema.obx.hpp`
    * there's an exception with `go` source & target type = the target type isn't present in the path
      e.g. `go/typeful/typebuf.obx.go.expected`
* `<source-type>/<test-case>/objectbox-model.json.expected` is the expected model JSON file, it's common for all languages.     
* `<source-type>/<test-case>/<target-type>/objectbox-model.<target-type-ext>.expected` is the expected model JSON file, it's common for all languages.
    * again with an exception to `go` where the target type isn't present in the path
* there can be a `<source-type>/<test-case>/objectbox-model.json.initial` 
    * it would be used as an initial value for the model JSON file before executing the generator,
    * otherwise (if not present), the initial model JSON isn't present (starting new model)
