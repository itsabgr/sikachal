package main

import (
	"database/sql"
	"flag"
	"github.com/itsabgr/sikachal"
	"github.com/itsabgr/sikachal/internal/common"
	"io"
	"log"
	"log/slog"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"time"
)

var flagDB = flag.String("db", ":memory:", "sqlite db uri")
var flagLogLevel = flag.Int("ll", int(slog.LevelDebug), "log level")
var flagCors = flag.String("cors", "*", "cors origin header value")
var flagKey = flag.String("key", "", "tls key file path")
var flagCert = flag.String("cert", "", "tls certificate file path")
var flagAddr = flag.String("addr", ":8080", "listening address")

func init() {
	flag.Parse()
}
func main() {
	db := common.Must(sql.Open("sqlite", *flagDB))
	defer common.Close(db)
	handler := sikachal.NewServer(
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(*flagLogLevel)})),
		sikachal.NewDB(db),
		*flagCors,
	)
	server := http.Server{
		Addr:              *flagAddr,
		Handler:           handler,
		IdleTimeout:       time.Second * 5,
		WriteTimeout:      time.Second * 5,
		MaxHeaderBytes:    10000, //10KB
		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second * 5,
		ErrorLog:          log.New(io.Discard, "", 0),
	}
	if *flagKey != "" {
		common.Throw(server.ListenAndServeTLS(*flagCert, *flagKey))
	} else {
		common.Throw(server.ListenAndServe())
	}
}
