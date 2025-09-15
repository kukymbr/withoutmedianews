package domain

const perPageDefault = 10

type PaginationReq struct {
	Page    uint
	PerPage uint
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

func (r PaginationReq) Offset() uint {
	return (r.Page - 1) * r.PerPage
}
