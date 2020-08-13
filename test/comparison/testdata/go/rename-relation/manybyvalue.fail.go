package object

// ERROR = can't merge model information: merging entity NegTaskRelManyValue: relation Groups: uid annotation value must not be empty (model relation UID = 4345851588384648695)

type NegTaskRelManyValue struct {
	Id     uint64
	Groups []Group `objectbox:"uid"`
}
