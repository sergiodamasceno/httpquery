package httpquery

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type QueryBuilder struct {
	columns   map[string]struct{}
	table     string
	paginator Paginator
}

func NewQueryBuilder(table string, columns map[string]struct{}) *QueryBuilder {
	p, err := NewPaginator(columns)
	if err != nil {
		log.Fatal("unable to build paginator", err)
	}
	return &QueryBuilder{table: table, columns: columns, paginator: p}
}

func (qb *QueryBuilder) BuildSQLQuery(schema string, values url.Values) (string, []interface{}, error) {

	paginating, err := qb.paginator.PaginatorFromQueryValues(&values)
	if err != nil {
		return "", nil, err
	}

	sqBuilder := sq.Select("*").From(fmt.Sprintf("%s.%s", schema, qb.table))
	for key, value := range values {
		if _, ok := qb.columns[key]; ok {
			sqBuilder = sqBuilder.Where(getCond(key, value[0]))
		}
	}

	sqBuilder = sqBuilder.Limit(paginating.Limit)
	sqBuilder = sqBuilder.Offset(paginating.Offset)
	sqBuilder = sqBuilder.OrderBy(fmt.Sprintf("%s %s", paginating.OrderBy, paginating.Direction))

	return sqBuilder.PlaceholderFormat(sq.Dollar).ToSql()
}

func getCond(column, value string) map[string]interface{} {

	i := strings.Index(value, ".")
	op := value[:i]
	v := value[i+1:]
	log.Println(op, v)

	switch op {
	case "lk":
		return sq.Like{column: v}
	case "gt":
		return sq.Gt{column: v}
	case "ne":
		return sq.NotEq{column: v}
	case "gte":
		return sq.GtOrEq{column: v}
	case "lt":
		return sq.Lt{column: v}
	case "lte":
		return sq.LtOrEq{column: v}
	}

	return sq.Eq{column: v}
}
