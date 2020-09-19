package httpquery

import (
	"net/url"
	"strings"
	"testing"
)

func TestPaginator(t *testing.T) {

	productColumns := map[string]struct{}{
		"created_at":     struct{}{},
		"updated_at":     struct{}{},
		"min_sell_price": struct{}{},
		"max_buy_price":  struct{}{},
		"name":           struct{}{},
	}

	paginator, err := NewPaginator(productColumns)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		in, want string // for want != "" we check error string, therefor we expect an error
	}{
		{"http://example.com/foo?orderby=created_at&offset=10&limit=50&direct=asc", ""},
		// {"http://example.com/foo?orderby=wrong_FIELD&offset=10&direct=asc&limit=50", "invalid field wrong_FIELD"},
	}
	for i, c := range cases {

		r, err := url.Parse(c.in)
		if err != nil {
			t.Fatal(err)
		}
		vals := r.Query()

		pagination, err := paginator.PaginatorFromQueryValues(&vals)

		if c.want != "" && err == nil {
			t.Fatalf("expected an error for case %d", i)
		}

		if c.want == "" && err != nil {
			t.Fatalf("unexpected error for case %d: %s", i, err)
		}

		if c.want != "" && !strings.Contains(err.Error(), c.want) {
			t.Errorf("expected error containing message %s, got: %s", c.want, err.Error())
		}

		if c.want == "" {
			if pagination.Direction == "" {
				t.Error("unmarshal missing direction")
			}
			if pagination.OrderBy == "" {
				t.Error("unmarshal missing sort by")
			}
		}
	}

}
