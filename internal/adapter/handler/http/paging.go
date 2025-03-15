package http

// paging represents as common pagination logic
// use it as embedded struct
const (
	DefaultPageSize     = 50
	DefaultSort         = "desc"
	DefaultPreloadLimit = 500
)

type paging struct {
	Cursor    *string `json:"cursor" form:"cursor" validate:"omitnil"`  // use cursor pagination
	PerPage   int     `json:"per_page" form:"per_page" validate:"gt=0"` // records per page (aka limit)
	Page      *int    `json:"page" form:"page" validate:"omitnil,gt=0"` // use page pagination
	SortOrder string  `json:"sort_order" form:"sort_order" validate:"oneof=asc desc"`
}

func (p *paging) DefaultPaging() {
	if p.PerPage <= 0 || p.PerPage > DefaultPageSize {
		p.PerPage = DefaultPageSize
	}
	if p.SortOrder == "" {
		p.SortOrder = DefaultSort
	}
}

func (p *paging) CursorFirstPage() bool {
	return p.Cursor == nil || *p.Cursor == "" || *p.Cursor == "\"\""
}
