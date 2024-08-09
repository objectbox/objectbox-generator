<img width="466" src="https://raw.githubusercontent.com/objectbox/objectbox-java/master/logo.png">
<br/>

[![Follow ObjectBox on Twitter](https://img.shields.io/twitter/follow/ObjectBox_io.svg?style=flat-square&logo=twitter&color=fff)](https://twitter.com/ObjectBox_io)

# ObjectBox Generator

Current version: 4.0.0-alpha2

ObjectBox is a superfast cross-platform object-oriented database.
ObjectBox Generator produces binding code for ObjectBox C, C++ and Go APIs (more languages to be supported in the future).
This greatly simplifies using ObjectBox in your favourite language.

## Download

Install the objectbox-generator by downloading the latest binary for your OS from [releases](https://github.com/objectbox/objectbox-generator/releases/latest).
If you want, add it to `$PATH` for convenience.

## Build 

If you prefer to build the generator yourself (vs. downloading), clone this repo and simply run `make`.
This will produce an "objectbox-generator" binary in the main project directory.

Build notes:

* To build yourself, you need Go, Make, CMake and a C++11 tool chain.
* To run test suite, run `make test-depend test`.
* `test-depend` needs to run only once to download objectbox core library and to build flatcc.
* A full test cycle can be triggered by `make clean all test-depend test`.

## Getting started

To get started, have a look at the specific language bindings (the approaches differ):

* C and C++ [repository](https://github.com/objectbox/objectbox-c) and [docs](https://cpp.objectbox.io/).
  In summary, you define a FlatBuffers schema file, and the ObjectBox Generator will create plain C++ data classes
  and helper classes that "glue" the data classes to the ObjectBox runtime library.
* Go [repository](https://github.com/objectbox/objectbox-go) and [docs](https://golang.objectbox.io/).
  Here, you start with Go data structs, for which the Generator generates the glue code directly.

## Development Notes

* Clean test cache: `go clean -testcache`
* Run test suite `test/comparison` with flag `-update` to update expected files.
* Run test suite `test/integration` with flag `-insource` to generate code in source may be helpful (e.g. `cd test/integration && go test ./... -insource`)

# License

```
ObjectBox Generator - a build time tool for ObjectBox
Copyright (C) 2018-2024 ObjectBox Ltd. All rights reserved.
https://objectbox.io
This file is part of ObjectBox Generator.

ObjectBox Generator is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
ObjectBox Generator is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.
You should have received a copy of the GNU Affero General Public License
along with ObjectBox Generator.  If not, see <https://www.gnu.org/licenses/>.
```

Note: AGPL only applies to the Generator itself and not to generated code.
You can regard generated code as "your code", and we impose no limitation on distributing it.
And, just to clarify: as the Generator does not include any warranty, there can be no warranty for the code it generates.       

# Do you ♥️ using ObjectBox?

We want to [hear about your project](https://docs.google.com/forms/d/e/1FAIpQLScIYiOIThcq-AnDVoCvnZOMgxO4S-fBtDSFPQfWldJnhi2c7Q/viewform)!
It will - literally - take just a minute, but help us a lot. Thank you!​ 🙏​
