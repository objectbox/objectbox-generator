package negative

// ERROR = can't prepare bindings for testdata/go/negative/missing-id.fail.go: no property recognized as an ID on entity MissingId

type MissingId struct {
	Text string
}
