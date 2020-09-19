package httpquery

import (
	"errors"
	"net/url"

	"github.com/gorilla/schema"
)

type Pagination struct {
	Offset    uint64 `schema:"offset"`
	Limit     uint64 `schema:"limit,required"`
	OrderBy   string `schema:"orderby,required"`
	Direction string `schema:"direct"`
}

type Paginator interface {
	PaginatorFromQueryValues(values *url.Values) (*Pagination, error)
}

type paginatorImp struct {
	decoder         *schema.Decoder
	sortableColumns map[string]struct{}
	validDirections map[string]struct{}
}

func NewPaginator(sortableColumns map[string]struct{}) (Paginator, error) {
	if len(sortableColumns) < 1 {
		return nil, errors.New("missing available fields to sort")
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	directions := map[string]struct{}{"asc": struct{}{}, "desc": struct{}{}}
	return &paginatorImp{decoder: decoder, sortableColumns: sortableColumns, validDirections: directions}, nil
}

func (p *paginatorImp) PaginatorFromQueryValues(values *url.Values) (*Pagination, error) {

	pagination := new(Pagination)
	err := p.decoder.Decode(pagination, *values)
	if err != nil {
		return nil, err
	}

	return pagination, nil
}
