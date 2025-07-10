package models

type PaginationSortBy string

const (
	SortByDesc PaginationSortBy = "desc"
	SortByAsc  PaginationSortBy = "asc"
)

type PaginationFiltering struct {
	Page    int
	PerPage int
	Search  string
	SortBy  PaginationSortBy
}
