package jmariadb

import (
	"fmt"
	"strings"

	"github.com/jigadhirasu/jgorm/jcommon"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func column(field ...string) string {
	cols := []string{}
	for _, f := range field {
		if !strings.Contains(f, "$") {
			cols = append(cols, fmt.Sprintf(`'$.%s'`, f))
			continue
		}
		cols = append(cols, fmt.Sprintf(`'%s'`, f))
	}

	if len(cols) < 1 {
		return ""
	}
	return fmt.Sprintf(`JSON_EXTRACT(Doc, %s)`, strings.Join(cols, ", "))
}

type C struct {
	Table     string   // 資料表明
	UUID      string   // 唯一編號
	Field     []string // 不用提取的屬性
	Attr      []string // 需要提取的屬性
	Index     KV       // 有索引的條件
	Condition KV       // 沒索引條件
}

func (c C) Select() string {
	col := "Doc"
	switch len(c.Field) + len(c.Attr) {
	case 0:
	case 1:
		col = strings.Join(c.Field, ",") + column(c.Attr...)
	default:
		ss := c.Field
		for _, f := range c.Attr {
			ss = append(ss, fmt.Sprintf("'%s'", f), column(f))
		}
		col = fmt.Sprintf(`JSON_OBJECT(%s)`, strings.Join(ss, ", "))
	}
	return col
}

func (c C) Where(tx *gorm.DB) *gorm.DB {
	tx = tx.Table(c.Table)
	if c.UUID != "" {
		tx = tx.Where("UUID = ?", c.UUID)
	}
	// 僅索引可以搜尋
	for k, v := range c.Index {
		tx = tx.Where(k+" = ?", v)
	}
	// 僅索引可以搜尋
	for k, v := range c.Condition {
		tx = tx.Where(datatypes.JSONQuery("Doc").Equals(v, k))
	}
	return tx
}

func (c C) Count(tx *gorm.DB) int64 {
	tx = c.Where(tx)
	var count int64
	tx.Count(&count)
	return count
}

func (c C) Value(tx *gorm.DB) jcommon.Bytes {
	tx = c.Where(tx)
	bb := [][]byte{}
	tx.Limit(1).Pluck(c.Select(), &bb)
	if len(bb) > 0 {
		return bb[0]
	}
	return jcommon.Bytes{}
}

func (c C) Values(tx *gorm.DB) jcommon.Bytes {
	tx = c.Where(tx)
	bb := [][]byte{}
	tx.Limit(1).Pluck(fmt.Sprintf(`JSON_ARRAYAGG(%s)`, c.Select())+" as Doc", &bb)
	if len(bb) > 0 {
		return bb[0]
	}
	return jcommon.Bytes{}
}

type KV map[string]interface{}
