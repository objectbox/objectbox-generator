package object

// ERROR = can't prepare bindings for cycles/embedding.fail.go: property EmbeddingChainA.EmbeddingChainB.EmbeddingChainC.EmbeddingChainA: embedded struct cycle detected: EmbeddingChainA.EmbeddingChainB.EmbeddingChainC

type EmbeddingChainA struct {
	Id uint64
	*EmbeddingChainB
}

type EmbeddingChainB struct {
	*EmbeddingChainC
}

type EmbeddingChainC struct {
	*EmbeddingChainA
}
