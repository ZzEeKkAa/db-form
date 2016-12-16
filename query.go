package main

import (
	"database/sql"
	"strings"
)

type Compilable interface {
	compile() string
}

type Query struct {
	name        string
	description string
	sql         string
	data        map[string]interface{}
}

func (q *Query) BuildDescription(action, selected string) string {
	res := "<p>" + q.description
	for k, v := range q.data {
		var val string
		switch v.(type) {
		case input:
			inp := v.(input)
			val = inp.compile()
		case selectbox:
			sb := v.(selectbox)
			sb.iselected = selected
			val = sb.compile()
		}
		res = strings.Replace(res, "{{"+k+"}}", val, -1)
	}

	res += "</p><p>" + (&input{
		itype:  "submit",
		ivalue: "Показать",
	}).compile() + "</p>"
	return "<form action='" + action + "' method='POST'>" + res + "</form>"
}

func (q *Query) BuildSql() string {
	res := q.sql
	for k, v := range q.data {
		var val string
		switch v.(type) {
		case input:
			inp := v.(input)
			val = inp.ivalue
		case selectbox:
			sb := v.(selectbox)
			val = sb.iselected
		}
		res = strings.Replace(res, "{{"+k+"}}", val, -1)
	}
	return res
}

func initQueries(db *sql.DB) (queries []Query) {
	sbBanners := selectbox{
		iname:  "banner",
		ivalue: map[string]string{},
	}
	sbScript := selectbox{
		iname:  "script",
		ivalue: map[string]string{},
	}

	sbBanners.loadMySQL(db, "banners", "banner_id", "banner_id", "name")
	sbScript.loadMySQL(db, "scripts", "url", "name")

	queries = append(queries, Query{
		name:        "SSP банера",
		description: "Вывести SSP, которые показывают банер {{banner}}",
		sql:         "select distinct ssp from companies as c,companies_show_banners as csb,banners as b where c.url=csb.companies_url AND csb.banners_banner_id={{banner}}",
		data: map[string]interface{}{
			"banner": sbBanners,
		},
	}, Query{
		name:        "Банера с лицензией",
		description: "Определить названия баннеров, что используют только те лицензии, что и баннер {{banner}}",
		sql: `select b0.name from banners as b0 where b0.banner_id not in
		(select b.banner_id
		from
		banners as b,
		bsl
		where
		b.banner_id = bsl.banner_id

		AND bsl.licence_url not IN (select bsl2.licence_url
		from
		bsl as bsl2,
		licences as l2
		where
		bsl2.banner_id={{banner}} AND
		bsl2.licence_url = l2.url)
		);`,
		data: map[string]interface{}{
			"banner": sbBanners,
		},
	}, Query{
		name:        "Скрипты",
		description: "Вывести названия скриптов, которые используются количеством банеров больше, чем у скрипта {{script}}",
		sql:         "select s.name from scripts as s, bsl where bsl.script_url = s.url group by s.name having count(*)>(select count(*) from bsl as bsl2 where bsl2.script_url='{{script}}');",
		data: map[string]interface{}{
			"script": sbScript,
		},
	}, Query{
		name:        "Банера, что продаются компаниям",
		description: "Определить названия банеров, которые показываются теми и только теми компаниями, что и банер {{banner}}",
		sql: `select b0.name
	from
		banners as b0
	where
		not exists (
			select * from companies_show_banners as csb0 where csb0.banners_banner_id=b0.banner_id and
				csb0.companies_url not in (select csb.companies_url from companies_show_banners as csb where csb.banners_banner_id={{banner}}))
		and not exists(
			select * from companies_show_banners as csb where csb.banners_banner_id={{banner}} and csb.companies_url not in
            			(select csb0.companies_url from companies_show_banners as csb0 where csb0.banners_banner_id=b0.banner_id)
        );`,
		data: map[string]interface{}{
			"banner": sbBanners,
		},
	})

	return
}

type Zvit struct {
	name string
	sql  string
}

func initZvits() []Zvit {

	return []Zvit{
		Zvit{
			name: "Количество баннеров по компаниям",
			sql: "select c.url as `Сайт компании`, p.name as `Владелец`, count(*) as `Количество банеров`" + `
	from
		companies as c,
        banners as b,
        presidents as p
	where
		b.companie_owner=c.url AND
        c.president = p.passport_id
group by c.url, p.name;`,
		},
		Zvit{
			name: "Под сколькью разными лицензиями используется каждый скрипт",
			sql:  "select s.name as script, count(distinct bsl.licence_url) as licences from scripts as s, bsl where bsl.script_url=s.url group by script",
		},
	}
}
