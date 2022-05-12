package hero

import (
	"github.com/jigadhirasu/common/j"
	"github.com/jigadhirasu/common/jmariadb"
)

func HeroPatch(tags j.Tags, h *Hero) j.Bytes {
	result := &jmariadb.Result{}
	db := jmariadb.MariaDB("asgame")
	jmariadb.Find(HeroPack{UUID: h.UUID})(db).Decode(result)

	if result.Length < 1 {
		return jmariadb.Create(tags, h)(db)
	}

	return jmariadb.Update(tags, h)(db)
}

// 列表
func HeroList(rq *HeroPack) j.Bytes {
	db := jmariadb.MariaDB("asgame")
	return jmariadb.Find(rq)(db)
}

// 取得一個
func HeroGet(h *Hero) j.Bytes {
	db := jmariadb.MariaDB("asgame")
	return jmariadb.Get(h)(db)
}
