package db

const perPageDefault = 10

type PaginationReq struct {
	Page    int
	PerPage int
}

func (r PaginationReq) GetNormalized() PaginationReq {
	normalized := r

	if normalized.PerPage == 0 {
		normalized.PerPage = perPageDefault
	}

	if normalized.Page <= 0 {
		normalized.Page = 1
	}

	return normalized
}

func (r PaginationReq) Offset() int {
	return (r.Page - 1) * r.PerPage
}
