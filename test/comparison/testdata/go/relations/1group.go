package object

type Group struct {
	Id   uint64
	Name string

	TaskRelPtrs []*TaskRelPtr `objectbox:"backlink:Group"` // Try backlink with a source property name.
}
