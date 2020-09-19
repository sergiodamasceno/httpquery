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

			valueWithPoint := value[0]
			i := strings.Index(valueWithPoint, ".")
			op := valueWithPoint[:i]
			v := valueWithPoint[i+1:]

			switch op {
			case "lk":
				sqBuilder = sqBuilder.Where(sq.Like{key: v})
			case "gt":
				sqBuilder = sqBuilder.Where(sq.Gt{key: v})
			case "ne":
				sqBuilder = sqBuilder.Where(sq.NotEq{key: v})
			case "gte":
				sqBuilder = sqBuilder.Where(sq.GtOrEq{key: v})
			case "lt":
				sqBuilder = sqBuilder.Where(sq.Lt{key: v})
			case "lte":
				sqBuilder = sqBuilder.Where(sq.LtOrEq{key: v})
			default:
				sqBuilder = sqBuilder.Where(sq.Eq{key: v})
			}

		}
	}

	sqBuilder = sqBuilder.Limit(paginating.Limit)
	sqBuilder = sqBuilder.Offset(paginating.Offset)
	sqBuilder = sqBuilder.OrderBy(fmt.Sprintf("%s %s", paginating.OrderBy, paginating.Direction))

	return sqBuilder.PlaceholderFormat(sq.Dollar).ToSql()
}
