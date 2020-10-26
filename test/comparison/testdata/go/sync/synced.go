package object

// Tests "sync" entity annotation and its interaction with standalone and property relations

// `objectbox:"sync"`
type SyncedEntity struct {
	Id            uint64
	PropertyRel   SyncedRelTarget `objectbox:"link"`
	StandaloneRel []SyncedRelTarget
}

// `objectbox:"sync"`
type SyncedRelTarget struct {
	Id uint64
}
