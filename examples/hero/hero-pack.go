package hero

import (
	"github.com/jigadhirasu/common/jmariadb"
	"gorm.io/gorm"
)

type HeroPack struct {
	UUID string `gorm:"column:UUID; uniqueIndex; type:varchar(40) AS (JSON_VALUE(doc, '$.UUID'))"`
	Name string `gorm:"column:Name; type:varchar(40) AS (JSON_VALUE(doc, '$.Name'))"`
	jmariadb.Pack
	jmariadb.Query
}

func (HeroPack) TableName() string {
	return Hero{}.TableName()
}

func (m HeroPack) Where(tx *gorm.DB) *gorm.DB {
	if m.UUID != "" {
		tx = tx.Where("UUID = ?", m.UUID)
	}
	if m.Name != "" {
		tx = tx.Where("Name = ?", m.Name)
	}
	return tx
}
