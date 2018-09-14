package thesaurus

// Thesaurus is ...
type Thesaurus interface {
	Synonyms(term string) ([]string, error)
}
