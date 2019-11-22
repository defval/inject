package di

// Extractor
type Extractor interface {
	Extract(params ExtractParams) error
}
