package models

const (
	TRUE  = "true"
	FALSE = "false"

	// Default ...
	DefaultPage                 = 1
	DefaultLimit                = 10
	DefaultLimitIsTotalDataTrue = 20

	// OrderBy ...
	OrderByASC  = "ASC"
	OrderByDESC = "DESC"
)

var (
	IsValidOrderBy = map[string]bool{
		OrderByASC:  true,
		OrderByDESC: true,
	}
)

type BasedFilter struct {
	Limit             int    `json:"limit" query:"limit"`
	Offset            int    `json:"offset" query:"offset"`
	Page              int    `json:"page" query:"page"`
	OrderBy           string `json:"order_by" query:"order_by"`
	SortBy            string `json:"sort_by" query:"sort_by"`
	IsNoLimit         bool   `json:"is_no_limit" query:"is_no_limit"`
	IsNotDefaultQuery bool   `json:"is_not_default_query" query:"is_not_default_query"`
	CustomColumns     string `json:"custom_columns" query:"custom_columns"`
}

func (c *BasedFilter) DefaultQuery() BasedFilter {
	if c.Limit <= 0 {
		c.Limit = 10
	}

	if c.Page <= 0 {
		c.Page = 1
	}

	if c.Page > 0 {
		c.Offset = (c.Page - 1) * c.Limit
	}

	return *c
}
