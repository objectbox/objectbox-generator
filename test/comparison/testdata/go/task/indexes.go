package object

type TaskIndexed struct {
	Id        uint64 `objectbox:"id"`
	Uid       string `objectbox:"unique"` // uses HASH as default
	UidValue  string `objectbox:"unique index:value"`
	UidHash   string `objectbox:"unique index:hash"`
	UidHash64 string `objectbox:"unique index:hash64"`
	UidInt    uint64 `objectbox:"unique"` // uses VALUE as default
	Name      string `objectbox:"index"`  // uses HASH as default
	Priority  int    `objectbox:"index"`  // uses VALUE as default
	Group     string `objectbox:"index:value"`
	Place     string `objectbox:"index:hash"`
	Source    string `objectbox:"index:hash64"`
}
