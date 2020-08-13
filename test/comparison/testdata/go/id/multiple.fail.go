package object

// ERROR = model finalization failed: entity Multiple 6:8325060299420976708 is invalid: multiple properties marked as ID: Id (1:7837839688282259259) and id2 (2:2518412263346885298)

type Multiple struct {
	Id  uint64 `objectbox:"id"`
	id2 uint64 `objectbox:"id"`
}
