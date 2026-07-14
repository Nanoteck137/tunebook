package types

type QueryParams struct {
	Filter string
	Sort   string
}

type PageParams struct {
	PerPage int
	Page    int
}
