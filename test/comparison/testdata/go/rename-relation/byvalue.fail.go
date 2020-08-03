package object

/* ERROR:
can't merge model information: merging entity NegTaskRelValue: property Group: uid annotation value must not be empty:
    [rename] apply the current UID 167566062957544642
    [change/reset] apply a new UID 3959279844101328186
*/

type NegTaskRelValue struct {
	Id    uint64
	Group Group `objectbox:"link uid"`
}
