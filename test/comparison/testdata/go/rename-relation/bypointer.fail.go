package object

/* ERROR:
can't merge model information: merging entity NegTaskRelPtr: property Group: uid annotation value must not be empty:
    [rename] apply the current UID 5974317550424871033
    [change/reset] apply a new UID 3959279844101328186
*/

type NegTaskRelPtr struct {
	Id    uint64
	Group *Group `objectbox:"link uid"`
}
