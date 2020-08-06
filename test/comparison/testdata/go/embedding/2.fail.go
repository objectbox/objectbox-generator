package object

// ERROR = can't prepare bindings for embedding/2.fail.go: multiple properties recognized as an ID: Id (0:0) and Id (0:0) on entity Negative2

// duplicate field
type Negative2 struct {
	Id                `objectbox:"inline"`
	IdAndFloat64Value `objectbox:"inline"`
}
