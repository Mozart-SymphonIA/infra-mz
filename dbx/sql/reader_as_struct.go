package sql

import (
	"context"
	dbsql "database/sql"
	"errors"
	"reflect"
	"strings"
)

func (s *sqlReader) QueryStruct(ctx context.Context, query string, dest any, args ...any) error {
	return s.scan(ctx, query, dest, true, args...)
}

func (s *sqlReader) QueryStructs(ctx context.Context, query string, dest any, args ...any) error {
	return s.scan(ctx, query, dest, false, args...)
}

func (s *sqlReader) scan(ctx context.Context, query string, dest any, single bool, args ...any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("dest must be pointer")
	}
	v = v.Elem()

	elemType, isSlice, isPtr := analyzeType(v.Type(), v.Kind())

	rows, err := s.c.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	fieldMap := buildFieldMap(elemType)

	cnt := 0
	for rows.Next() {
		cnt++
		if err := processRow(rows, cols, fieldMap, elemType, isSlice, isPtr, &v); err != nil {
			return err
		}
		if !isSlice {
			break
		}
	}

	if single && cnt == 0 {
		return dbsql.ErrNoRows
	}
	return rows.Err()
}

func analyzeType(t reflect.Type, k reflect.Kind) (reflect.Type, bool, bool) {
	isSlice := k == reflect.Slice
	if isSlice {
		t = t.Elem()
	}
	isPtr := t.Kind() == reflect.Ptr
	if isPtr {
		t = t.Elem()
	}
	return t, isSlice, isPtr
}

func buildFieldMap(elemType reflect.Type) map[string]int {
	fieldMap := make(map[string]int)
	for i := 0; i < elemType.NumField(); i++ {
		f := elemType.Field(i)
		if f.PkgPath != "" {
			continue
		}
		tag := f.Tag.Get("db")
		if tag == "" {
			tag = f.Name
		}
		fieldMap[strings.ToLower(tag)] = i
	}
	return fieldMap
}

func processRow(rows *dbsql.Rows, cols []string, fieldMap map[string]int, elemType reflect.Type, isSlice, isPtr bool, v *reflect.Value) error {
	newVal := reflect.New(elemType).Elem()
	scans := make([]any, len(cols))
	for i, col := range cols {
		if idx, ok := fieldMap[strings.ToLower(col)]; ok {
			scans[i] = newVal.Field(idx).Addr().Interface()
		} else {
			var skip any
			scans[i] = &skip
		}
	}
	if err := rows.Scan(scans...); err != nil {
		return err
	}

	if isSlice {
		appendVal := newVal
		if isPtr {
			appendVal = newVal.Addr()
		}
		v.Set(reflect.Append(*v, appendVal))
	} else {
		v.Set(newVal)
	}
	return nil
}
