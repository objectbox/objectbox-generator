package object

type Teacher struct {
	Id       uint64
	Students []*Student `objectbox:"backlink"` // Try backlink without a source property name.
}

type Student struct {
	Id      uint64
	Teacher *Teacher `objectbox:"link"`
}
