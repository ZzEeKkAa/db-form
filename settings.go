package main

import "strings"

type Settings struct {
	tableKey      map[string]([]string)
	connectionKey map[string]string
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
	set.connectionColumns = make(map[string][]string,0)
	set.connectionKey = make(map[string]string,0)
	set.SetTableKey("adgroups", "adgroup_id")
	set.SetTableKey("banners", "banner_id")
	set.SetTableKey("companies", "url")

	set.SetConnection("companies", "president", "presidents", "passprot_id", "name")
}
