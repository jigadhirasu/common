package hero

type Hero struct {
	UUID string // 編號
	Name string // 資料庫名稱
	Rank int    // 稀有階級
}

func (m *Hero) UU(uuid ...string) string {
	if len(uuid) > 0 {
		m.UUID = uuid[0]
	}
	return m.UUID
}
func (m Hero) TableName() string {
	return "heros"
}
