package domain

type PaginationReq struct {
	Page    int
	PerPage int
}

func (r PaginationReq) Offset() int {
	return (r.Page - 1) * r.PerPage
}
