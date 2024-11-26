package bqb

import "testing"

func Test_dialectReplace_raw_dialect(t *testing.T) {
	wantSQL := "SELECT * FROM t WHERE id = 1 AND name = 'foo'"

	r := RawReplacer{
		SQL:    "SELECT * FROM t WHERE id = {{xX_PARAM_Xx}} AND name = {{xX_PARAM_Xx}}",
		Params: []interface{}{1, "foo"},
	}
	sql, err := r.GetReplacedSQL()

	if err != nil {
		t.Error("valid SQL statement should not return error")
	}

	if sql != wantSQL {
		t.Errorf("unexpected SQL statement: want %s got %s", wantSQL, sql)
	}
}

func Test_dialectReplace_mySQL_replacer(t *testing.T) {
	wantSQL := "SELECT * FROM t WHERE id = ? AND name = ?"

	m := MySQLReplacer{
		SQL: "SELECT * FROM t WHERE id = {{xX_PARAM_Xx}} AND name = {{xX_PARAM_Xx}}",
	}

	sql, err := m.GetReplacedSQL()

	if err != nil {
		t.Error("valid SQL statement should not return error")
	}

	if sql != wantSQL {
		t.Errorf("unexpected statement: want %s got %s", wantSQL, sql)
	}
}

func Test_dialectReplace_Sql_replacer(t *testing.T) {
	wantSQL := "SELECT * FROM t WHERE id = ? AND name = ?"

	s := SQLReplacer{
		SQL: "SELECT * FROM t WHERE id = {{xX_PARAM_Xx}} AND name = {{xX_PARAM_Xx}}",
	}

	sql, err := s.GetReplacedSQL()

	if err != nil {
		t.Error("valid SQL statement should not return error")
	}

	if sql != wantSQL {
		t.Errorf("unexpected statement: want %s got %s", wantSQL, sql)
	}
}

func Test_dialectReplace_PgSQLReplacer(t *testing.T) {
	wantSQL := "SELECT * FROM t WHERE id = $1 AND name = $2"

	p := PgSQLReplacer{
		SQL:    "SELECT * FROM t WHERE id = {{xX_PARAM_Xx}} AND name = {{xX_PARAM_Xx}}",
		Params: []interface{}{1, "foo"},
	}

	sql, err := p.GetReplacedSQL()

	if err != nil {
		t.Errorf("valid SQL statement should not return error")
	}

	if sql != wantSQL {
		t.Errorf("unexpected statement: want %s got %s", wantSQL, sql)
	}
}
