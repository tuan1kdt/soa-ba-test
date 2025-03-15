package util

// Paging represents as common pagination logic
// use it as embedded struct
const (
	DefaultPageSize     = 50
	DefaultSort         = "desc"
	DefaultPreloadLimit = 500
)

type Paging struct {
	Cursor    *string
	PerPage   int
	Page      *int
	SortOrder string
}

func (p *Paging) DefaultPaging() {
	if p.PerPage <= 0 || p.PerPage > DefaultPageSize {
		p.PerPage = DefaultPageSize
	}
	if p.SortOrder == "" {
		p.SortOrder = DefaultSort
	}
}

func (p *Paging) CursorFirstPage() bool {
	return p.Cursor == nil || *p.Cursor == "" || *p.Cursor == "\"\""
}
