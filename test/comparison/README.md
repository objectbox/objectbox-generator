# Code comparison tests

For a few basic scenarios **only(!)**, we compare the generated code with the expected code.
Rule of thumb: **avoid any redundancy in the generated code**.
For example, a string property should be generated **only once** for each possible annotation value (e.g., index).
If you see redundancy, you are welcome to remove it.

Use with care as this is a messy way to test in the sense of generating a lot of text;
sometimes even for simple changes.
The signal-to-noise ratio in commits is usually bad as generated reference files "drown out" the actual changes.
One direct consequence is that developers may miss errors in the generated code. 

Thus, in general, **we prefer other ways of testing**.
E.g., actually compile the generated code, run it and check its behavior as part of a test.
This can still be combined with a few targeted "text contains" checks in the generated code.

## Expected UIDs

Note: **the generated UIDs do not consider the existing ones in the JSON file!**

In tests, UIDs are generated using a preseeded random number generator (given via Options to the generator).
While this gives us a fixed sequence of numbers, certain changes "avalanche" different UIDs to lots of properties.
The reason for this is that UIDs are first assigned to all entities, and then to all properties.
Thus, if a new class is added, the UIDs of all properties are changed.
This makes the commits "noisier" and harder to review than necessary.

## Update expected files

If you read the disclaimer above, and you still want to extend the "expected files" tests (are you sure?):

* Edit generator_test.go: set the Flag value of overwriteExpected to `true`
* Run generator_test.go (TestCompare) or run `go test test`: this overwrites the existing expected files
* Review the changes (the expected files that are generated) and make changes if needed 
* Edit generator_test.go: set the Flag value of overwriteExpected back to `false`
* Run generator_test.go (TestCompare) or run `go test test` to verify
* Commit

## Basic outline

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
