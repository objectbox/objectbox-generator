package object

// ERROR = can't merge model information: merging entity Negative1: property Value: duplicate property name (note that property names are case insensitive)

// both contain Value field but of two distinct types
type Negative1 struct {
	Id           uint64
	Float64Value `objectbox:"inline"`
	BytesValue   `objectbox:"inline"`
}
