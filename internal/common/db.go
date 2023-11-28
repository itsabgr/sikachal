package common

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
)

func exportedFields(fields []reflect.StructField) []reflect.StructField {
	result := make([]reflect.StructField, 0, len(fields))
	for _, f := range fields {
		if f.IsExported() {
			result = append(result, f)
		}
	}
	return result
}
func scanRowIntoStruct[T any](row interface{ Scan(dest ...any) error }) (*T, error) {
	var t T
	tval := reflect.Indirect(reflect.ValueOf(&t))
	fields := exportedFields(reflect.VisibleFields(tval.Type()))
	if len(fields) == 0 {
		panic(errors.New("no exported field"))
	}
	vals, err := scanRow(row, len(fields))
	if err != nil {
		return nil, err
	}
	for i, a := range vals {
		field := tval.FieldByName(fields[i].Name)
		tval.FieldByName(fields[i].Name).Set(reflect.ValueOf(a).Convert(field.Type()))
	}
	return &t, nil
}
func scanRow(row interface{ Scan(dest ...any) error }, count int) ([]any, error) {
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for i := range values {
		valuePtrs[i] = &values[i]
	}
	return values, row.Scan(valuePtrs...)
}

func QueryRow[T any](ctx context.Context, db interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}, query string, args ...any) (t *T, _ error) {
	row := db.QueryRowContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return t, err
	}
	return scanRowIntoStruct[T](row)
}
