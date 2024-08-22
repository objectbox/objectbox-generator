package object

import "github.com/objectbox/objectbox-generator/v4/test/comparison/testdata/go/embedding/other"

type E struct {
	other.Trackable `objectbox:"inline"`
	id              uint64
	other.ForeignAlias
	other.ForeignNamed
}
