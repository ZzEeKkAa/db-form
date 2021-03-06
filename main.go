package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strings"

	"log"

	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

var set Settings
var queries []Query
var zvits []Zvit

func main() {
	initSettings()
	zvits = initZvits()
	flag.Parse()

	//db, err := sql.Open("mysql", "root:toortoortoor@tcp(localhost:3306)/rtb")
	db, err := sql.Open("mysql", "root@tcp(192.168.2.84:3306)/rtb")
	//db, err := sql.Open("mysql", "root@tcp(localhost:3306)/rtb")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	queries = initQueries(db)

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

	updateTable := func(postArgs *fasthttp.Args, table string) error {
		var err error

		var sqlUpdate []string
		var primaryKeys []string
		var primaryValues []string
		//var sqlArgs []string
		postArgs.VisitAll(func(key, value []byte) {
			val := string(value)
			switch string(key) {
			case "primary_keys":
				primaryKeys = strings.Split(val, ",")
			case "primary_values":
				primaryValues = strings.Split(val, ",")
			default:
				if string(value) != "" {
					sqlUpdate = append(sqlUpdate, "`"+string(key)+"` = '"+string(value)+"'")
				} else {
					sqlUpdate = append(sqlUpdate, "`"+string(key)+"` = null")
				}
			}
		})

		sql := ""
		insert := true
		for _, v := range primaryValues {
			if v != "" {
				insert = false
			}
		}
		if insert {
			sql += "INSERT INTO `" + table + "` SET " + strings.Join(sqlUpdate, ", ")
		} else {
			sql += "UPDATE `" + table + "` SET " + strings.Join(sqlUpdate, ", ") + " WHERE "
			var whereRule []string
			for i := range primaryKeys {
				whereRule = append(whereRule, " `"+primaryKeys[i]+"`='"+primaryValues[i]+"'")
			}
			sql += strings.Join(whereRule, " AND ")
		}
		if len(sqlUpdate) > 0 {
			fmt.Println(sql)
			_, err = db.Query(sql)
		}
		return err
	}

	h := func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, `<html><head><link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
  <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.1/jquery.min.js"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script><style>body{text-align:center} form{margin: 0 auto; width:600px}</style></head><body>`)

		path := strings.Split(string(ctx.Path()), "/")
		if path[1] == "insert" {
			_, err = stmtIns.Exec("http://pupkin.com/", nil, "Pupkin SSP", "Pupkin DSP")

			if err != nil {
				panic(err.Error())
			}
		} else if path[1] == "table" {
			err := updateTable(ctx.PostArgs(), path[2])
			fmt.Fprintf(ctx, "<h3>Form `"+path[2]+"`</h3>")
			if err != nil {
				fmt.Fprintf(ctx, "<div class='alert alert-danger'>"+err.Error()+"</div>")
			}
			var f form
			page64, err := strconv.ParseInt(path[3], 10, 10)
			if err != nil {
				panic(err.Error())
			}
			page := int(page64)
			f.loadMySQL(db, path[2], page)
			fmt.Fprintf(ctx, f.compile())
			if page > 1 {
				fmt.Fprintf(ctx, "<a class='btn btn-info' role='button' href='/table/"+path[2]+"/"+strconv.Itoa(page-1)+"'> Prev </a>")
			}
			fmt.Fprintf(ctx, "<a class='btn btn-info' role='button' href='/table/"+path[2]+"/"+strconv.Itoa(page+1)+"'> Next </a>")
		} else if path[1] == "zvit" {
			zvit_id, err := strconv.ParseInt(path[2], 10, 64)
			if err != nil || zvit_id < 0 || int(zvit_id) >= len(zvits) {
				fmt.Fprintf(ctx, "<h3>Unknown zvit!</h3>")
			} else {
				zvit := &zvits[zvit_id]

				fmt.Fprintf(ctx, "<h3>"+zvit.name+"</h3>")

				sql := zvit.sql

				q, err := db.Prepare(sql)
				if err != nil {
					panic(err)
				}
				rows, err := q.Query()
				columns, err := rows.Columns()

				resp := []interface{}{}
				for range columns {
					resp = append(resp, new([]byte))
				}
				fmt.Fprintf(ctx, "<table class='table'><thead><tr><th>"+strings.Join(columns, "</th><th>")+"</th></tr><tbody>")
				for rows.Next() {
					rows.Scan(resp...)
					var vals []string
					for _, i := range resp {
						vals = append(vals, string(*i.(*[]byte)))
					}
					fmt.Fprintf(ctx, "<tr><td>"+strings.Join(vals, "</td><td>")+"</td></tr>")
				}

				fmt.Fprintf(ctx, "</tbody></table>")
			}
		} else if path[1] == "query" {
			query_id, err := strconv.ParseInt(path[2], 10, 64)
			if err != nil || query_id < 0 || int(query_id) >= len(queries) {
				fmt.Fprintf(ctx, "<h3>Unknown query!</h3>")
			} else {
				query := &queries[query_id]

				showTable := false
				ctx.PostArgs().VisitAll(func(key, value []byte) {
					_, ok := query.data[string(key)]
					if ok {
						showTable = true
						switch query.data[string(key)].(type) {
						case input:
							inp := query.data[string(key)].(input)
							inp.ivalue = string(value)
							query.data[string(key)] = inp
						case selectbox:
							sb := query.data[string(key)].(selectbox)
							sb.iselected = string(value)
							query.data[string(key)] = sb
						}
					}
				})

				fmt.Fprintf(ctx, "<h3>"+query.name+"</h3>")
				fmt.Fprintf(ctx, query.BuildDescription("/query/"+path[2], ""))
				if showTable {
					sql := query.BuildSql()
					fmt.Println(sql)
					q, err := db.Prepare(sql)
					if err != nil {
						panic(err)
					}
					rows, err := q.Query()
					columns, err := rows.Columns()

					resp := []interface{}{}
					for range columns {
						resp = append(resp, new([]byte))
					}
					fmt.Fprintf(ctx, "<table class='table'><thead><tr><th>"+strings.Join(columns, "</th><th>")+"</th></tr><tbody>")
					for rows.Next() {
						rows.Scan(resp...)
						var vals []string
						for _, i := range resp {
							vals = append(vals, string(*i.(*[]byte)))
						}
						fmt.Fprintf(ctx, "<tr><td>"+strings.Join(vals, "</td><td>")+"</td></tr>")
					}

					fmt.Fprintf(ctx, "</tbody></table>")
				}
			}
		} else {
			var table string
			tables, err := getTables.Query()
			if err != nil {
				panic(err.Error())
			}
			fmt.Fprintf(ctx, "<h3>Forms</h3><div class='list-group'>")
			for tables.Next() {
				tables.Scan(&table)
				fmt.Fprintf(ctx, "<a href='/table/"+table+"/1' class='list-group-item'>"+table+"</a>")
			}
			fmt.Fprintf(ctx, "</div>")

			fmt.Fprintf(ctx, "<h3>Queries</h3><div class='list-group'>")
			for i, q := range queries {
				fmt.Fprintf(ctx, "<a href='/query/"+strconv.FormatInt(int64(i), 10)+"' class='list-group-item'>"+q.name+"</a>")
			}
			fmt.Fprintf(ctx, "</div>")

			fmt.Fprintf(ctx, "<h3>Zvits</h3><div class='list-group'>")
			for i, z := range zvits {
				fmt.Fprintf(ctx, "<a href='/zvit/"+strconv.FormatInt(int64(i), 10)+"' class='list-group-item'>"+z.name+"</a>")
			}
			fmt.Fprintf(ctx, "</div>")
		}

		fmt.Fprintf(ctx, `</body></html>`)
		ctx.SetContentType("text/html; charset=UTF-8")
	}
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
