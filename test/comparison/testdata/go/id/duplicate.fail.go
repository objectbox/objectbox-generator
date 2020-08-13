package object

// ERROR = can't prepare bindings for id/duplicate.fail.go: multiple properties recognized as an ID: Id (0:0) and id (0:0) on entity Duplicate

type Duplicate struct {
	Id uint64
	id uint64
}
