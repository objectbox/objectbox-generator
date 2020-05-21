Generated from flatbuffers reflection.fbs:

* cd to flatbuffers repo
* check out the revision you want
* build flatbuffers compiler `flatc`
* change `PATH` with absolute path to `fbsparser directory`
```shell script
path=??????
rm -rv ${path}/reflection
flatc --go -o ${path} reflection/reflection.fbs
```
