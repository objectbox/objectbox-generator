package externaltype

// `objectbox:"external-name:MyExtEntityName"`
type Entity struct {
	ID            int64          `objectbox:"id"`
	Name          string         `objectbox:"external-name:MyExtPropName"`
	Email         string         `objectbox:"external-name:MyExtPropName2 external-type:Uuid"`
	ChildEntities []*ChildEntity `objectbox:"external-name:MyExtRelName external-type:MongoId"`
}

type ChildEntity struct {
	Id       uint64
	ImageURL string
}
