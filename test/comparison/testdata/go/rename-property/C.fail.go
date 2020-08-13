package object

// ERROR = can't merge model information: merging entity C: property New: uid annotation value must not be empty, the property isn't present in the persisted model

// negative test, tag `objectbox:"uid"` on an unknown property
type C struct {
	Id  uint64
	New string `objectbox:"uid"`
}
