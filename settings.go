package main

import "strings"

type Settings struct {
	tableKey          map[string]([]string)
	connectionKey     map[string]string
	connectionColumns map[string][]string
}

func (s *Settings) SetTableKey(table string, key ...string) {
	if s.tableKey == nil {
		s.tableKey = make(map[string]([]string), 0)
	}
	s.tableKey[table] = key
}

func (s *Settings) GetTableKey(table string) []string {
	return s.tableKey[table]
}

func (s *Settings) SetConnection(srcTable, srcKey, distTable, distKey string, columnsToShow ...string) {
	if s.tableKey == nil {
		s.tableKey = make(map[string]([]string), 0)
	}
	s.connectionKey[srcTable+"::"+srcKey] = distTable + "::" + distKey
	s.connectionColumns[srcTable+"::"+srcKey] = columnsToShow
}

func (s *Settings) GetConnection(table, key string) (tab, col string, columnsToShow []string, ok bool) {
	distString, ok := s.connectionKey[table+"::"+key]
	if ok {
		dist := strings.Split(distString, "::")
		tab, col = dist[0], dist[1]
		columnsToShow = s.connectionColumns[table+"::"+key]
	}
	return
}

func initSettings() {
	set.connectionColumns = make(map[string][]string, 0)
	set.connectionKey = make(map[string]string, 0)
	set.SetTableKey("presidents", "passport_id")
	set.SetTableKey("companies", "url")
	set.SetTableKey("companies_show_banners", "banners_banner_id", "companies_url")
	set.SetTableKey("banners", "banner_id")
	set.SetTableKey("static_banners", "banner_id")
	set.SetTableKey("video_banners", "banner_id")
	set.SetTableKey("interactive_banners", "banner_id")
	set.SetTableKey("scripts", "url")
	set.SetTableKey("licences", "url")
	set.SetTableKey("bsl", "licence_url", "banner_id", "script_url")

	set.SetTableKey("adgroups", "adgroup_id")
	set.SetConnection("companies", "president", "presidents", "passport_id", "name", "born_year")
	set.SetConnection("banners", "companie_owner", "companies", "url", "url")
	set.SetConnection("companies_show_banners", "banners_banner_id", "banners", "banner_id", "banner_id", "name")
	set.SetConnection("companies_show_banners", "companies_url", "companies", "url", "url")
	set.SetConnection("static_banners", "banner_id", "banners", "banner_id", "banner_id", "name")
	set.SetConnection("video_banners", "banner_id", "banners", "banner_id", "banner_id", "name")
	set.SetConnection("interactive_banners", "banner_id", "banners", "banner_id", "banner_id", "name")

	set.SetConnection("bsl", "banner_id", "interactive_banners", "banner_id", "banner_id", "name")
	set.SetConnection("bsl", "licence_url", "licences", "url", "name")
	set.SetConnection("bsl", "script_url", "scripts", "url", "name", "year")

	set.SetConnection("scripts", "based_script_url", "scripts", "url", "name", "year")

}
