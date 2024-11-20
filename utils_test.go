package bqb

import (
	"testing"
)

func Test_dialectReplace_unknown_dialect(t *testing.T) {
	const (
		testSql = "test-sql"
	)
	params := []any{1, 2, "a", "c"}
	sql, err := dialectReplace(Dialect("unknown"), testSql, params)

	if sql != "test-sql" {
		t.Errorf("unexpected sql statement: want %s got %s", testSql, sql)
	}

	if err != nil {
		t.Error("unknown dialect should not return an error")
	}
}

func Test_dialectReplace_raw_dialect(t *testing.T) {
	wantSql := "SELECT * FROM t WHERE id = 1"

	r := RawReplacer{
		Sql:    "SELECT * FROM t WHERE id = {{xX_PARAM_Xx}}",
		Params: []interface{}{1},
	}
	sql, err := dialectReplace2(&r)

	if err != nil {
		t.Error("valid sql statement should not return error")
	}

	if sql != wantSql {
		t.Errorf("unexpected sql statement: want %s got %s", wantSql, sql)
	}
}
