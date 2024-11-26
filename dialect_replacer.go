package bqb

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	questionMark                = "?"
	doubleQuestionMarkDelimiter = "??"
)

// DialectReplacer defines the required behavior for implementing
// a replacer of parameters in SQL statements.
type DialectReplacer interface {
	GetReplacedSQL() (string, error)
}

type RawReplacer struct {
	SQL    string
	Params []any
}

func (r RawReplacer) GetReplacedSQL() (string, error) {
	sql := r.SQL
	for _, param := range r.Params {
		p, err := paramToRaw(param)
		if err != nil {
			return "", err
		}
		sql = strings.Replace(sql, paramPh, p, 1)
	}
	return sql, nil
}

func paramToRaw(param any) (string, error) {
	switch p := param.(type) {
	case bool:
		return fmt.Sprintf("%v", p), nil
	case float32, float64, int, int8, int16, int32, int64,
		uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", p), nil
	case *int:
		if p == nil {
			return "NULL", nil
		}
		return fmt.Sprintf("%v", *p), nil
	case string:
		return fmt.Sprintf("'%v'", p), nil
	case *string:
		if p == nil {
			return "NULL", nil
		}
		return fmt.Sprintf("'%v'", *p), nil
	case nil:
		return "NULL", nil
	default:
		return "", fmt.Errorf("unsupported type for Raw query: %T", p)
	}
}

type SQLReplacer struct {
	SQL string
}

func (s SQLReplacer) GetReplacedSQL() (string, error) {
	return replaceQuestionMarksInSQL(s.SQL)
}

type MySQLReplacer struct {
	SQL string
}

func (m MySQLReplacer) GetReplacedSQL() (string, error) {
	return replaceQuestionMarksInSQL(m.SQL)
}

func replaceQuestionMarksInSQL(sql string) (string, error) {
	return strings.ReplaceAll(sql, paramPh, questionMark), nil
}

type PgSQLReplacer struct {
	SQL    string
	Params []any
}

func (p PgSQLReplacer) GetReplacedSQL() (string, error) {
	p.SQL = strings.ReplaceAll(p.SQL, doubleQuestionMarkDelimiter, questionMark)
	parts := strings.Split(p.SQL, paramPh)
	var builder strings.Builder
	for i := range p.Params {
		_, _ = builder.WriteString(parts[i] + "$" + strconv.Itoa(i+1))
	}
	builder.WriteString(parts[len(parts)-1])
	return builder.String(), nil
}
