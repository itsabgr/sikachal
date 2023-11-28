package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/itsabgr/sikachal"
	"github.com/itsabgr/sikachal/internal/common"
	_ "modernc.org/sqlite"
)

var flagDB = flag.String("db", "", "sqlite db uri")
var flagCount = flag.Uint64("count", 1_000_000, "insert count")

func init() {
	flag.Parse()
}
func main() {
	db := common.Must(sql.Open("sqlite", *flagDB))
	defer common.Close(db)
	appDB := sikachal.NewDB(db)
	common.Throw(appDB.Create(context.Background()))
	common.Throw(appDB.BulkUserInsert(context.Background(), *flagCount))
}
