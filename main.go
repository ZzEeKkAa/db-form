package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strings"

	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

func main() {
	flag.Parse()

	db, err := sql.Open("mysql", "root:toortoortoor@tcp(localhost:3306)/rtb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO companies VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	getTables, err := db.Prepare("show tables;")
	if err != nil {
		panic(err.Error())
	}
	defer getTables.Close()

	h := func(ctx *fasthttp.RequestCtx) {
		path := strings.Split(string(ctx.Path()), "/")
		if path[1] == "info" {
			fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
			fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
			fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
			fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
			fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
			fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
			fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
			fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
			fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
			fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

			fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)

			ctx.SetContentType("text/plain; charset=UTF-8")
		} else if path[1] == "insert" {
			_, err = stmtIns.Exec("http://pupkin.com/", nil, "Pupkin SSP", "Pupkin DSP")

			if err != nil {
				panic(err.Error())
			}
		} else if path[1] == "table" {
			var f form
			f.loadMySQL(db, path[2])
			fmt.Fprintf(ctx, f.compile())
			ctx.SetContentType("text/html; charset=UTF-8")
		} else {
			var table string
			tables, err := getTables.Query()
			if err != nil {
				panic(err.Error())
			}
			for tables.Next() {
				tables.Scan(&table)
				fmt.Fprintf(ctx, "<p><a href='/table/"+table+"'>"+table+"</a></p>")
			}

			fmt.Fprintf(ctx, "<p><a href='/insert/'>Insert</a></p>")
			ctx.SetContentType("text/html; charset=UTF-8")
		}
	}
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
