package query

//ParseFullTextSearch converts the string received from the client to FullTextSearch struct.
func ParseFullTextSearch(input string) *FullTextSearch {
	if len(input) == 0 {
		return nil
	}

	return &FullTextSearch{Query: input}
}

// GoString implements fmt.GoStringer interface
// Returns string representation of sorting in next form:
// "<name> (ASC|DESC) [, <tag_name> (ASC|DESC)]"
func (fts FullTextSearch) GoString() string {
	return fts.Query
}
