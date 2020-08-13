package object

// ERROR = can't merge model information: merging entity NegTaskRelManyPtr: relation Groups: uid annotation value must not be empty (model relation UID = 8514850266767180993)

type NegTaskRelManyPtr struct {
	Id     uint64
	Groups []*Group `objectbox:"uid"`
}
