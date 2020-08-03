package object

/* ERROR:
can't merge model information: merging entity A: property Old: uid annotation value must not be empty:
    [rename] apply the current UID 3390393562759376202
    [change/reset] apply a new UID 6050128673802995827
*/

// negative test, tag `objectbox:"uid"` will cause the build tool to print the UID of the property and fail
type A struct {
	Id  uint64
	Old string `objectbox:"uid"`
}
