package querybuilder

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type Cond struct {
	selectField []string
	omitField   []string
	limit       *int
	offset      *int
	query       string
	params      []any
	operator    string
	order       interface{}
}

func New() *Cond {
	return new(Cond)
}

func (c *Cond) Build(tx *gorm.DB) *gorm.DB {
	if c.selectField != nil {
		tx = tx.Select(c.selectField)
	}
	if c.omitField != nil {
		tx = tx.Omit(c.omitField...)
	}
	if notEmptyString(c.query) {
		tx = tx.Where(c.query, c.params...)
	}
	if c.limit != nil {
		tx = tx.Limit(*c.limit)
	}
	if c.offset != nil {
		tx = tx.Offset(*c.offset)
	}
	if c.order != nil {
		tx = tx.Order(c.order)
	}

	return tx
}

func (c *Cond) Select(selectField ...string) *Cond {
	if selectField != nil {
		c.selectField = selectField
	}
	return c
}

func (c *Cond) Omit(omitField ...string) *Cond {
	if omitField != nil {
		c.omitField = omitField
	}
	return c
}

func (c *Cond) Limit(limit int) *Cond {
	c.limit = &limit
	return c
}

func (c *Cond) Offset(offset int) *Cond {
	c.offset = &offset
	return c
}

func (c *Cond) Order(order interface{}) *Cond {
	c.order = order
	return c
}

func (c *Cond) Associate(args ...Builder) Builder {
	asso := &associate{
		cond:        c,
		listBuilder: args,
	}
	return asso
}

func (condition *Cond) not() *Cond {
	sb := strings.Builder{}
	sb.WriteString("NOT(")
	sb.WriteString(condition.query)
	sb.WriteString(")")

	condition.query = sb.String()
	condition.operator = ""
	return condition
}

func (condition *Cond) appendCondition(appendOperator string, conditions ...*Cond) *Cond {
	// wrap previous condition in () if needed
	builder := strings.Builder{}
	if condition.operator != appendOperator && !isEmptyString(condition.operator) &&
		!isEmptyString(condition.query) {
		condition.query = wrapQuery(condition.query)
	}
	builder.WriteString(condition.query)
	condition.operator = appendOperator
	for _, c := range conditions {
		if c.selectField != nil {
			condition.selectField = c.selectField
		}
		if c.omitField != nil {
			condition.omitField = c.omitField
		}
		if c.limit != nil {
			condition.limit = c.limit
		}
		if c.offset != nil {
			condition.offset = c.offset
		}

		if !isEmptyString(c.query) {
			if c.operator != appendOperator && !isEmptyString(c.operator) {
				c.query = wrapQuery(c.query)
			}
			if builder.Len() > 0 {
				builder.WriteString(appendOperator)
			}
			builder.WriteString(c.query)
			condition.params = append(condition.params, c.params...)
		}
	}
	condition.query = builder.String()
	return condition
}

func wrapQuery(query string) string {
	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(query)
	sb.WriteString(")")

	return sb.String()
}

func Select(fields ...string) *Cond {
	return &Cond{
		selectField: fields,
	}
}

func Omit(fields ...string) *Cond {
	return &Cond{
		omitField: fields,
	}
}

func Limit(limit int) *Cond {
	return &Cond{limit: &limit}
}

func Offset(offset int) *Cond {
	return &Cond{offset: &offset}
}

func Order(order interface{}) *Cond {
	return &Cond{order: order}
}

// Equal represents "field = value".
func Equal(field string, value interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" = ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// NotEqual represents "field <> value".
func NotEqual(field string, value any) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" <> ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// GreaterThan represents "field > value".
func GreaterThan(field string, value interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" > ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// GreaterEqualThan represents "field >= value".
func GreaterEqualThan(field string, value interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" >= ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// LessThan represents "field < value".
func LessThan(field string, value interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" < ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// LessEqualThan represents "field <= value".
func LessEqualThan(field string, value interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" <= ?")

	return &Cond{
		query:  sb.String(),
		params: []interface{}{value},
	}
}

func appendSliceWhereIN(value interface{}, values ...interface{}) []interface{} {
	var results []interface{}
	if value == nil {
		return values
	}

	rv := reflect.ValueOf(value)
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		for i := range rv.Len() {
			results = append(results, rv.Index(i).Interface())
		}
	} else {
		results = append(results, value)
	}

	results = append(results, values...)

	return results
}

// In represents "field IN (value...)".
// Examples:
// + In("id", []int{1,2,3})
// + In("id", []any{1,2,3})
// + In("id", 1,2,3)
func In(field string, value interface{}, values ...interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" IN (?)")

	return &Cond{
		query:  sb.String(),
		params: []any{appendSliceWhereIN(value, values...)},
	}
}

// NotIn represents "field NOT IN (value...)".
func NotIn(field string, value interface{}, values ...interface{}) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" NOT IN (?)")

	return &Cond{
		query:  sb.String(),
		params: []any{appendSliceWhereIN(value, values...)},
	}
}

// Like represents "field LIKE value".
func Like(field string, value string) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" LIKE ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// NotLike represents "field NOT LIKE value".
func NotLike(field string, value string) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" NOT LIKE ?")

	return &Cond{
		query:  sb.String(),
		params: []any{value},
	}
}

// IsNull represents "field IS NULL".
func IsNull(field string) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" IS NULL")

	return &Cond{
		query:  sb.String(),
		params: []any{},
	}
}

// IsNotNull represents "field IS NOT NULL".
func IsNotNull(field string) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" IS NOT NULL")

	return &Cond{
		query:  sb.String(),
		params: []any{},
	}
}

// Between represents "field BETWEEN lower AND upper".
func Between(field string, lower, upper string) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString(" BETWEEN ? AND ?")

	return &Cond{
		query:  sb.String(),
		params: []any{lower, upper},
	}
}

// NotBetween represents "field NOT BETWEEN lower AND upper".
func NotBetween(field string, lower, upper string) *Cond {
	sb := strings.Builder{}
	sb.WriteString(field)
	sb.WriteString("NOT BETWEEN ? AND ?")

	return &Cond{
		query:  sb.String(),
		params: []any{lower, upper},
	}
}

// And will Join simple a slice of condition into a condition with AND WiseFunc for where statement
func And(conditions ...*Cond) *Cond {
	object := &Cond{}
	return object.appendCondition(" AND ", conditions...)
}

// Raw represents and raw query
func Raw(query string, params []interface{}) *Cond {
	object := &Cond{
		query:  query,
		params: params,
	}
	return object
}

// Or will Join simple a slice of condition into a condition with OR WiseFunc for where statement
func Or(conditions ...*Cond) *Cond {
	object := &Cond{}
	return object.appendCondition(" OR ", conditions...)
}

func Not(condition *Cond) *Cond {
	object := &Cond{
		query:    condition.query,
		params:   condition.params,
		operator: condition.operator,
	}
	return object.not()
}

func isEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}
