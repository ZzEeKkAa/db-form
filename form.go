package main

import (
	"database/sql"
	"strconv"
	"strings"
)

type form struct {
	method string
	action string
	inputs []input
}

type input struct {
	idescript string
	itype     string
	iname     string
	ivalue    string
}

func (i *input) compile() string {
	var res string
	res += "<p>"

	if i.idescript != "" {
		res += i.idescript + ":</br>"
	}

	res += "<input class='form-control' type='" + i.itype + "'"
	if i.ivalue != "" {
		res += " value='" + i.ivalue + "'"
	}
	if i.iname != "" {
		res += " name='" + i.iname + "'"
	}
	res += " /></p>"
	return res
}

func (f *form) compile() string {
	var res string
	res += "<form class='form-group'"
	if f.method != "" {
		res += " method='" + f.method + "'"
	}
	if f.action != "" {
		res += " action='" + f.action + "'"
	}
	res += ">"
	for _, i := range f.inputs {
		res += i.compile()
	}
	res += "</form>"
	return res
}

func (f *form) loadMySQL(db *sql.DB, tableName string, page int) {
	sql := "SELECT * FROM " + tableName + " LIMIT " + strconv.Itoa(page-1) + ",1"

	stmtGet, err := db.Prepare(sql)

	if err != nil {
		panic(err)
	}

	rows, err := stmtGet.Query()
	defer rows.Close()

	columns, err := rows.Columns()

	resp := []interface{}{}
	for range columns {
		resp = append(resp, new([]byte))
	}

	if rows.Next() {
		err = rows.Scan(resp...)
	}

	for i, column := range columns {
		f.inputs = append(f.inputs, input{idescript: column, iname: column, ivalue: string(*resp[i].(*[]byte)), itype: "text"})
	}

	primary_keys := set.GetTableKey(tableName)
	primary_values := make([]string, len(primary_keys))
	primary_changed := make([]bool, len(primary_keys))
	for i, column := range columns {
		for j, col := range primary_values {
			if col == column && !primary_changed[j] {
				primary_changed[j] = true
				primary_values[j] = string(*resp[i].(*[]byte))
			}
		}
	}
	f.inputs = append(f.inputs, input{itype: "hidden", iname: "primary_keys", ivalue: strings.Join(primary_keys, ",")})
	f.inputs = append(f.inputs, input{itype: "hidden", iname: "primary_values", ivalue: strings.Join(primary_values, ",")})

	f.inputs = append(f.inputs, input{itype: "submit", ivalue: "Save"})
	f.method = "POST"
	f.action = "./" + strconv.Itoa(page)
}

func (f *form) updateMySQL(db *sql.DB, tableName string, page int) {
	//stmtGet, err := db.Prepare("SELECT * FROM " + tableName + " LIMIT " + strconv.Itoa(page-1) + ",1")
}
