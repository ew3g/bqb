package bqb

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

func dialectReplace(dr DialectReplacer) (string, error) {
	return dr.GetReplacedSQL()
}

func convertArg(text string, arg any) (string, []any, []error) {
	var newArgs []any
	var errs []error

	switch v := arg.(type) {

	case Embedder:
		text = strings.Replace(text, "?", v.RawValue(), 1)

	case driver.Valuer:
		text = strings.Replace(text, "?", paramPh, 1)
		val, err := v.Value()
		if err != nil {
			errs = append(errs, err)
		} else {
			newArgs = append(newArgs, val)
		}
	case []int:
		newPh := []string{}
		for _, i := range v {
			newPh = append(newPh, paramPh)
			newArgs = append(newArgs, i)
		}
		text = strings.Replace(text, "?", strings.Join(newPh, ","), 1)

	case []*int:
		newPh := []string{}
		for _, i := range v {
			newPh = append(newPh, paramPh)
			newArgs = append(newArgs, i)
		}
		if len(newPh) > 0 {
			text = strings.Replace(text, "?", strings.Join(newPh, ","), 1)
		} else {
			text = strings.Replace(text, "?", paramPh, 1)
			newArgs = append(newArgs, nil)
		}

	case []string:
		newPh := []string{}
		for _, s := range v {
			newPh = append(newPh, paramPh)
			newArgs = append(newArgs, s)
		}
		text = strings.Replace(text, "?", strings.Join(newPh, ","), 1)

	case []*string:
		newPh := []string{}
		for _, s := range v {
			newPh = append(newPh, paramPh)
			newArgs = append(newArgs, s)
		}
		if len(newPh) > 0 {
			text = strings.Replace(text, "?", strings.Join(newPh, ","), 1)
		} else {
			text = strings.Replace(text, "?", paramPh, 1)
			newArgs = append(newArgs, nil)
		}

	case []any:
		newPh := []string{}
		for _, s := range v {
			newPh = append(newPh, paramPh)
			newArgs = append(newArgs, s)
		}
		text = strings.Replace(text, "?", strings.Join(newPh, ","), 1)

	case *Query:
		if v == nil {
			text = strings.Replace(text, "?", paramPh, 1)
			newArgs = append(newArgs, nil)
			return text, newArgs, errs
		}
		sql, params, err := v.toSql()
		text = strings.Replace(text, "?", sql, 1)
		if err != nil {
			errs = append(errs, err)
		} else {
			newArgs = append(newArgs, params...)
		}

	case JsonMap, JsonList:
		text = strings.Replace(text, "?", paramPh, 1)
		bytes, err := json.Marshal(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("cann jsonify struct: %v", err))
		} else {
			newArgs = append(newArgs, string(bytes))
		}

	case *JsonMap, *JsonList:
		text = strings.Replace(text, "?", paramPh, 1)
		bytes, err := json.Marshal(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("cann jsonify struct: %v", err))
		} else {
			newArgs = append(newArgs, string(bytes))
		}

	case Embedded:
		text = strings.Replace(text, "?", string(v), 1)

	default:
		text = strings.Replace(text, "?", paramPh, 1)
		newArgs = append(newArgs, v)
	}

	return text, newArgs, errs
}

func checkParamCounts(text, original string, args []any) error {
	extraCount := strings.Count(text, "?")
	if extraCount > 0 {
		return fmt.Errorf("extra ? in text: %v (%d args)", original, len(args))
	}

	paramCount := strings.Count(text, paramPh)
	if paramCount < len(args) {
		return fmt.Errorf("missing ? in text: %v (%d args)", original, len(args))
	}
	return nil
}

func makePart(text string, args ...any) QueryPart {
	tempPh := "XXX___XXX"
	originalText := text
	text = strings.ReplaceAll(text, "??", tempPh)

	var newArgs []any
	errs := make([]error, 0)

	for _, arg := range args {
		argText, fArgs, argErrs := convertArg(text, arg)
		if len(argErrs) > 0 {
			errs = append(errs, argErrs...)
		}
		newArgs = append(newArgs, fArgs...)
		text = strings.ReplaceAll(argText, "??", tempPh)
	}

	if err := checkParamCounts(text, originalText, newArgs); err != nil {
		errs = append(errs, err)
	}

	text = strings.ReplaceAll(text, tempPh, "??")

	return QueryPart{
		Text:   text,
		Params: newArgs,
		Errs:   errs,
	}
}
