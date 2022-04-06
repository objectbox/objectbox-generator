package object

// ERROR = can't prepare bindings for cycles/embedding-named.fail.go: property EmbeddingNamedChainA.BPtr.CPtr.APtr: embedded struct cycle detected: EmbeddingNamedChainA.BPtr.CPtr

type EmbeddingNamedChainA struct {
	Id   uint64
	BPtr *EmbeddingNamedChainB
}

type EmbeddingNamedChainB struct {
	CPtr *EmbeddingNamedChainC
}

type EmbeddingNamedChainC struct {
	APtr *EmbeddingNamedChainA
}
