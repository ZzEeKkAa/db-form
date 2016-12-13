package main

import "strings"

type Settings struct {
	tableKey      map[string]([]string)
	connectionKey map[string]string
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

func (s *Settings) SetConnection(srcTable, srcKey, distTable, distKey string) {
	if s.tableKey == nil {
		s.tableKey = make(map[string]([]string), 0)
	}
	s.connectionKey[srcTable+"::"+srcKey] = distTable + "::" + distKey
}

func (s *Settings) GetConnection(table, key string) (tab, col string, ok bool) {
	distString, ok := s.connectionKey[table+"::"+key]
	if ok {
		dist := strings.Split(distString, "::")
		tab, col = dist[0], dist[1]
	}
	return
}

func initSettings() {
	set.SetTableKey("adgroups", "adgroup_id")
	set.SetTableKey("banners", "banner_id")
}
