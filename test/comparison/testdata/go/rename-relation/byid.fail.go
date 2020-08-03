package object

/* ERROR:
can't merge model information: merging entity NegTaskRelId: property Group: uid annotation value must not be empty:
    [rename] apply the current UID 6745438398739480977
    [change/reset] apply a new UID 3959279844101328186
*/

type NegTaskRelId struct {
	Id    uint64
	Group uint64 `objectbox:"link:Group uid"`
}
