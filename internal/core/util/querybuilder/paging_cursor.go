package querybuilder

import (
	"fmt"

	"gorm.io/gorm"
)

// cursorForwarder represents as record forwarder in cursor base pagination
type cursorForwarder struct {
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

// genericCursor represents a generic map
type genericCursor map[string]any

func createGenericCursor(field string, value any, pointsNext bool) genericCursor {
	return genericCursor{
		field:         value,
		"points_next": pointsNext,
	}
}

type cursorPaging struct {
	field       string
	limit       int
	sortOrder   sortOrder
	cursor      *string
	isFirstPage bool

	pointsNext    bool
	hasPagination bool
	prevCursor    cursorField
	nextCursor    cursorField
}

func NewCursorPaging(cursor *string, field string, opts ...CursorPagingOpt) *cursorPaging {
	cp := &cursorPaging{
		cursor:      cursor,
		field:       field,
		isFirstPage: cursor == nil,
	}

	for _, opt := range opts {
		opt(cp)
	}

	return cp
}

type CursorPagingOpt func(*cursorPaging)

func WithCursorSortOrder(sortOrder string) CursorPagingOpt {
	return func(cp *cursorPaging) {
		cp.sortOrder = sortOrderFromString(sortOrder)
	}
}

func WithCursorLimit(limit int) CursorPagingOpt {
	return func(cp *cursorPaging) {
		cp.limit = limit
	}
}

type cursorField interface {
	CursorField() any
}

func (c *cursorPaging) Pagination(isFirstPage, hasPagination bool, prev, next cursorField) cursorForwarder {
	if isFirstPage {
		if hasPagination {
			nextCursor := createGenericCursor(c.field, next.CursorField(), true)
			return newCursorPagination(nextCursor, nil)
		}
	} else {
		if c.pointsNext {
			var nextCur genericCursor
			if hasPagination {
				nextCur = createGenericCursor(c.field, next.CursorField(), true)
			}
			prevCur := createGenericCursor(c.field, prev.CursorField(), false)
			return newCursorPagination(nextCur, prevCur)
		} else {
			var prevCur genericCursor
			nextCur := createGenericCursor(c.field, next.CursorField(), true)
			if hasPagination {
				prevCur = createGenericCursor(c.field, prev.CursorField(), false)
			}
			return newCursorPagination(nextCur, prevCur)
		}
	}
	return cursorForwarder{}
}

func (c *cursorPaging) PointNext() bool {
	return c.pointsNext
}

func (c *cursorPaging) Build(tx *gorm.DB) (paginatorTX *gorm.DB) {
	defer func() {
		paginatorTX = tx.Order(fmt.Sprintf("%s %s", c.field, c.sortOrder)).Limit(c.limit + 1) // limit + 1 to detect has next pagination
	}()

	if c.cursor == nil || *(c.cursor) == "\"\"" || *c.cursor == "" {
		return
	}

	decodedCursor, err := decodeCursor(*c.cursor)
	if err != nil {
		return tx
	}
	c.pointsNext = decodedCursor["points_next"] == true
	operator, order := paginationOperator(c.pointsNext, c.sortOrder)

	whereStr := fmt.Sprintf("%s %s ?", c.field, operator)
	tx = tx.Where(whereStr, decodedCursor[c.field])
	if order != "" {
		c.sortOrder = order
	}

	return tx
}

// newCursorPagination generates the CursorPagination
func newCursorPagination(next genericCursor, prev genericCursor) cursorForwarder {
	return cursorForwarder{
		Next: encodeCursor(next),
		Prev: encodeCursor(prev),
	}
}

func paginationOperator(pointsNext bool, sortOrder sortOrder) (string, sortOrder) {
	if pointsNext && sortOrder == sortOrderASC {
		return ">", ""
	}
	if pointsNext && sortOrder == sortOrderDESC {
		return "<", ""
	}
	if !pointsNext && sortOrder == sortOrderASC {
		return "<", sortOrderDESC
	}
	if !pointsNext && sortOrder == sortOrderDESC {
		return ">", sortOrderASC
	}

	return "", ""
}
