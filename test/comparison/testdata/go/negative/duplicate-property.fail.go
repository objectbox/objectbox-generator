package negative

// ERROR = can't merge model information: merging entity DuplicateProperty: property text: duplicate property name (note that property names are case insensitive)

type DuplicateProperty struct {
	Id   uint64 `objectbox:"id"`
	Text string `objectbox:"name:text"`
	text string
}
