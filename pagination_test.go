package pagination_test

import (
	"testing"
	"time"

	"github.com/AndiVS/pagination"
)

type MyTestFilter struct {
	Name      string
	StartDate time.Time
	Age       int
}

func (flt *MyTestFilter) ToSQL() (string, []any) {
	q := "WHERE name = $1 AND start_date > $2 AND age > $3"
	args := []any{flt.Name, flt.StartDate, flt.Age}
	return q, args
}

func TestPaginationQuery(t *testing.T) {
	// define filter which implements pagination.Filter interface
	flt := &MyTestFilter{
		Name:      "John",
		StartDate: time.Date(2022, 10, 10, 0, 0, 0, 0, time.UTC),
		Age:       22,
	}

	// build pagination.Pagination struct providing necessary options
	p := pagination.Pagination[*MyTestFilter]{
		Sort: pagination.Sort{
			OrderBy: "name",
			Asc:     false,
		},
		Limit:  20,
		Offset: 100,
		Filter: flt,
	}

	// call ToSQL method which returns sql clause and args for placeholders
	actualQuery, actualArgs := p.ToSQL()

	expectedQuery := "WHERE name = $1 AND start_date > $2 AND age > $3 ORDER BY name DESC LIMIT 20 OFFSET 100"
	if expectedQuery != actualQuery {
		t.Fatalf("incorrect query has been built, expected is %s, but got %s", expectedQuery, actualQuery)
		t.FailNow()
	}

	expectedArgs := []any{flt.Name, flt.StartDate, flt.Age}
	for i := range expectedArgs {
		if expectedArgs[i] != actualArgs[i] {
			t.Fatalf("incorrect placeholder param on index position %d, expected is %v, but got %v", i+1, expectedQuery, actualQuery)
			t.FailNow()
		}
	}
}
