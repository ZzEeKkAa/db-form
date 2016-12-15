package main

import (
	"database/sql"
	"strconv"
	"strings"
)

type form struct {
	method      string
	action      string
	inputs      []input
	selectboxes []selectbox
}

type input struct {
	idescript string
	itype     string
	iname     string
	ivalue    string
}

type selectbox struct {
	idescript string
	iname     string
	iselected string
	ivalue    map[string]string
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
	for _, s := range f.selectboxes {
		res += s.compile()
	}
	inp := input{itype: "submit", ivalue: "Save"}
	res += inp.compile()
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
		depTable, depCol, columnsToShow, ok := set.GetConnection(tableName, column)
		if ok {
			sb := selectbox{idescript: column, iname: column, iselected: string(*resp[i].(*[]byte)), ivalue: map[string]string{}}
			sb.loadMySQL(db, depTable, depCol, columnsToShow...)
			f.selectboxes = append(f.selectboxes, sb)
		} else {
			f.inputs = append(f.inputs, input{idescript: column, iname: column, ivalue: string(*resp[i].(*[]byte)), itype: "text"})
		}
	}

	primary_keys := set.GetTableKey(tableName)
	primary_values := make([]string, len(primary_keys))
	primary_changed := make([]bool, len(primary_keys))
	copy(primary_values, primary_keys)
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

	f.method = "POST"
	f.action = "./" + strconv.Itoa(page)
}

func (s *selectbox) loadMySQL(db *sql.DB, tableName, columnName string, columnsToShow ...string) {
	sql := "SELECT * FROM " + tableName
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

	s.ivalue[""] = "null"
	for rows.Next() {
		err = rows.Scan(resp...)
		key := ""
		var showValues []string
		for _, col := range columnsToShow {
			for j, colResp := range columns {
				if col == colResp {
					colStr := string(*resp[j].(*[]byte))
					showValues = append(showValues, colStr)
				}
			}
		}

		for i, column := range columns {
			if column == columnName {
				key = string(*resp[i].(*[]byte))
			}
		}
		if key != "" {
			s.ivalue[key] = strings.Join(showValues, ", ")
		}
	}
}

func (s *selectbox) compile() string {
	var res string
	res += "<p>"

	if s.idescript != "" {
		res += s.idescript + ":</br>"
	}

	res += "<select class='form-control'"
	if s.iname != "" {
		res += " name='" + s.iname + "'"
	}
	res += ">"

	for key, val := range s.ivalue {

		res += "<option value='" + key + "'"
		if key == s.iselected {
			res += " selected"
		}
		res += ">" + val + "</option>"
	}

	res += "</select></p>"
	return res
}
