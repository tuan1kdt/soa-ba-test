package querybuilder

import (
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

// offsetForwarder represents as record forwarder in offset base pagination
type offsetForwarder struct {
	CurrentPage int `json:"current_page,omitempty"`
	PerPage     int `json:"per_page,omitempty"`
	TotalPage   int `json:"total_page,omitempty"`
}

type offsetPaging struct {
	orderFields []string
	limit       int
	page        int
	sortOrder   sortOrder
	totalPage   int
}

// NewOffsetPaging use limit offset for pagination
// WARN: this solution is not recommended, as it can produce slow queries with a large offset.
// Use cursor-based pagination instead
// If the limit offset is required and offset is a large value, please read the following article
// https://hackmysql.com/deferred-join-deep-dive/
func NewOffsetPaging(page, totalRecords int, orderFields []string, opts ...OffsetPagingOpt) *offsetPaging {
	op := &offsetPaging{
		orderFields: orderFields,
		page:        page,
	}

	for _, opt := range opts {
		opt(op)
	}
	op.totalPage = int(math.Ceil(float64(totalRecords) / float64(op.limit)))

	return op
}

type OffsetPagingOpt func(*offsetPaging)

func WithOffsetSortOrder(sortOrder string) OffsetPagingOpt {
	return func(op *offsetPaging) {
		op.sortOrder = sortOrderFromString(sortOrder)
	}
}

func WithOffsetLimit(limit int) OffsetPagingOpt {
	return func(op *offsetPaging) {
		op.limit = limit
	}
}

func (c *offsetPaging) Build(tx *gorm.DB) *gorm.DB {
	orders := make([]string, len(c.orderFields))
	for i := range c.orderFields {
		orders[i] = fmt.Sprintf("%s %s", c.orderFields[i], c.sortOrder)
	}
	tx = tx.Order(strings.Join(orders, ",")).Limit(c.limit).Offset((c.page - 1) * c.limit)

	return tx
}

func (c *offsetPaging) Pagination() offsetForwarder {
	return offsetForwarder{
		CurrentPage: c.page,
		PerPage:     c.limit,
		TotalPage:   c.totalPage,
	}
}
