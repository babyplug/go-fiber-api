package repo

type PaginationMetadata struct {
	Page       uint `json:"page"`
	PerPage    uint `json:"per_page"`
	TotalPages uint `json:"total_pages"`
	TotalItems uint `json:"total_items"`
}