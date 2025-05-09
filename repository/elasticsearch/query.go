package elasticsearch

type EsQuery interface {
	AddField(field string) EsQuery
	Build() map[string]any
}

type query struct {
	query  string
	fields []string
}

func NewQuery(q string) EsQuery {
	return &query{
		query:  q,
		fields: []string{},
	}
}

func (q *query) AddField(field string) EsQuery {
	q.fields = append(q.fields, field)
	return q
}

func (q *query) Build() map[string]any {
	query := map[string]any{
		"match_all": map[string]any{},
	}

	if q.query != "" {
		query = map[string]any{
			"multi_match": map[string]any{
				"query":  q.query,
				"fields": q.fields,
			},
		}
	}

	return map[string]any{
		"query": query,
		"size":  10,
	}
}
