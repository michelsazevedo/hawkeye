package domain

type Searchable interface {
	GetID() string
}

type Course struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Headline string `json:"headline"`
	Author   string `json:"author"`
	Thumb    string `json:"thumb"`
}

func (c *Course) GetID() string {
	return c.Id
}
