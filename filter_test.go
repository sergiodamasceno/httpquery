package httpquery

import (
	"database/sql"
	"log"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/schema"
	_ "github.com/lib/pq"
)

var (
	decoder = schema.NewDecoder()
	err     error
	db      *sql.DB
)

var pc = []string{
	"created_at",
	"updated_at",
	"min_sell_price",
	"max_buy_price",
	"name",
	"category_id",
}

func init() {
	decoder.IgnoreUnknownKeys(true)
	if err != nil {
		log.Fatal(err)
	}
	db, err = sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
}

func TestFilterSchema(t *testing.T) {

	builder := NewQueryBuilder("product", TestModelQueryColumns)
	paginating := "&limit=50&orderby=name&direct=ASC"

	t.Run("findAtLeastOne", func(t *testing.T) {
		cases := []struct {
			in, wantContains string
		}{
			{"http://example.com/products?name=eq.TOWING" + paginating, "name ="},
			{"http://example.com/products?updated_at=eq.2020-08-31T01:40:09.158813Z" + paginating, "updated_at ="},
			{"http://example.com/products?name=lk.YUBA" + paginating, "name LIKE"},
			{"http://example.com/products?name=ilk.YUBA" + paginating, "name ILIKE"},
		}
		for _, c := range cases {

			r, err := url.Parse(c.in)
			if err != nil {
				log.Fatal(err)
			}

			sql, args, err := builder.BuildSQLQuery("public", r.Query())
			if err != nil {
				t.Fatal(c.in, err)
			}

			if !strings.Contains(sql, c.wantContains) {
				t.Errorf("expected sql containing: %s. Not found at: %s", c.wantContains, sql)
			}
			t.Log(sql, args)

		}
	})

	t.Run("validate", func(t *testing.T) {
		cases := []string{
			"http://example.com/products",
		}
		for _, strURL := range cases {
			r, err := url.Parse(strURL)
			if err != nil {
				log.Fatal(err)
			}
			_, _, err = builder.BuildSQLQuery("public", r.Query())
			if err == nil {
				t.Errorf("Failed, expected an error validating the missing params URL: %s", strURL)
			}
		}
	})

}

var (
	value, column, operator string
	typ                     reflect.Type
	vals                    reflect.Value
)
