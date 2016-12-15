package main

import (
	"bytes"
	"text/template"
)

type Compilable interface {
	compile() string
}

type query struct {
	description string
	sql         string
	data        interface{}
}

func (q *query) BuildDescription() string {
	buf := bytes.NewBuffer([]byte{})
	t := template.New("Description")
	t.Parse(q.description)
	t.Execute(buf, q.data)
	return string(buf.Bytes())
}

func (q *query) BuildSql() string {
	return ""
}

func initQueries() (queries []query) {
	queries = append(queries, query{
		description: "Вывести SSP, которые показывают банер {{.banner}}",
		sql:         "select ssp from companies,banners where url=companie_owner AND banner_id={{.banner}}",
		data: struct {
			banner string
		}{},
	})

	return
}
