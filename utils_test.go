package bqb

import (
	"testing"
)

type TestReplacer struct {
	SQL string
}

func (tr TestReplacer) GetReplacedSQL() (string, error) {
	return tr.SQL, nil
}

func Test_dialectReplacer_dialect(t *testing.T) {
	tr := TestReplacer{
		SQL: "my test sql",
	}

	wantSQL := "my test sql"

	sql, err := dialectReplace(tr)

	if err != nil {
		t.Error("valid SQL statement should not return error")
	}

	if sql != wantSQL {
		t.Errorf("unexpected SQL statement: want %s got %s", wantSQL, sql)
	}
}
