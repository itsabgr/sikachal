package common

import (
	"context"
	"testing"
)
import _ "modernc.org/sqlite"
import "database/sql"

func TestDbBinding(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec("create table test(str text,num int,blb blob);"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("insert into test values (?,?,?)", "hi", 122, []byte("binary")); err != nil {
		t.Fatal(err)
	}
	result, err := QueryRow[struct {
		unexported  int
		Num         int
		Str         string
		unexported2 int
		Blob        []byte
		unexported3 int
	}](context.Background(), db, "select num,str,blb from test limit 1")
	if err != nil {
		t.Fatal(err)
	}
	if result.Str != "hi" {
		t.Fail()
	}
	if result.Num != 122 {
		t.Fail()
	}
	if string(result.Blob) != "binary" {
		t.Failed()
	}
	if _, err := db.Exec("insert into test values (?,?,?)", "ho", 123, []byte("text")); err != nil {
		t.Fatal(err)
	}

}
