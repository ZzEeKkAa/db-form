package main

import "database/sql"

type form struct {
	method string
	inputs []input
}

type input struct {
	idescript string
	itype     string
	ivalue    string
}

func (i *input) compile() string {
	var res string
	res += "<p>"

	if i.idescript != "" {
		res += i.idescript + ":</br>"
	}

	res += "<input type='" + i.itype + "'"
	if i.ivalue != "" {
		res += " value='" + i.ivalue + "'"
	}
	res += " /></p>"
	return res
}

func (f *form) compile() string {
	var res string
	res += "<form>"
	for _, i := range f.inputs {
		res += i.compile()
	}
	res += "</form>"
	return res
}

func (f *form) loadMySQL(db *sql.DB, tableName string, page int) {
	stmtGet, err := db.Prepare("SELECT * FROM " + tableName + " ")

	if err != nil {
		panic(err)
	}

	rows, _ := stmtGet.Query()

	columns, err := rows.Columns()

	for _, column := range columns {
		f.inputs = append(f.inputs, input{idescript: column, itype: "text"})
	}

	f.inputs = append(f.inputs, input{itype: "submit", ivalue: "Save"})
	f.method = "POST"
}
