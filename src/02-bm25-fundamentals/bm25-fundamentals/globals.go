package main

const (
	// Application name
	AppName = "bm25-fundamentals"
	
	// BM25 Parameters (hardcoded in SQLite FTS5, provided here for reference/education)
	BM25_K1 = 1.2  // Term frequency saturation parameter
	BM25_B  = 0.75 // Document length normalization parameter
)