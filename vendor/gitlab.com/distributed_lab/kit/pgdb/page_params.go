package pgdb

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

const (
	// OrderTypeAsc means result should be sorted in ascending order.
	OrderTypeAsc = "asc"
	// OrderTypeDesc means result should be sorted in descending order.
	OrderTypeDesc = "desc"
)

// OffsetPageParams defines page params for offset-based pagination.
type OffsetPageParams struct {
	Limit      uint64 `url:"page[limit]"`
	PageNumber uint64 `url:"page[number]"`
	Order      string `url:"page[order]"`
}

// ApplyTo returns a new SelectBuilder after applying the paging effects of  `p` to `sql`.
func (p *OffsetPageParams) ApplyTo(sql squirrel.SelectBuilder, cols ...string) squirrel.SelectBuilder {
	offset := p.Limit * p.PageNumber

	sql = sql.Limit(p.Limit).Offset(offset)

	switch p.Order {
	case OrderTypeAsc:
		for _, col := range cols {
			sql = sql.OrderBy(fmt.Sprintf("%s %s", col, "asc"))
		}
	case OrderTypeDesc:
		for _, col := range cols {
			sql = sql.OrderBy(fmt.Sprintf("%s %s", col, "desc"))
		}
	default:
		panic(fmt.Errorf("unexpected order type: %v", p.Order))
	}

	return sql
}

//CursorPageParams - page params of the db query
type CursorPageParams struct {
	Cursor uint64
	Order  string
	Limit  uint64
}

// ApplyTo returns a new SelectBuilder after applying the paging effects of
// `p` to `sql`.  This method provides the default case for paging: int64
// cursor-based paging by an id column.
func (p *CursorPageParams) ApplyTo(sql squirrel.SelectBuilder, col string) squirrel.SelectBuilder {
	sql = sql.Limit(p.Limit)

	switch p.Order {
	case OrderTypeAsc:
		sql = sql.
			Where(fmt.Sprintf("%s > ?", col), p.Cursor).
			OrderBy(fmt.Sprintf("%s asc", col))
	case OrderTypeDesc:
		sql = sql.
			Where(fmt.Sprintf("%s < ?", col), p.Cursor).
			OrderBy(fmt.Sprintf("%s desc", col))
	default:
		panic(fmt.Errorf("unexpected order type: %v", p.Order))
	}

	return sql
}
