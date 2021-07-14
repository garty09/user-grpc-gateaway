package pagination

var DefaultPageSize = 10

type Pages struct {
	page    int `json:"page"`
	perPage int `json:"per_page"`
}

func New(page int) *Pages {
	peerPage := DefaultPageSize
	if page < 1 {
		page = 1
	}

	return &Pages{
		page:    page,
		perPage: peerPage,
	}
}

func (p *Pages) Offset() int {
	return (p.page - 1) * p.perPage
}

func (p *Pages) Limit() int {
	return p.perPage
}
