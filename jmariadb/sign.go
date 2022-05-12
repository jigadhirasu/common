package jmariadb

import (
	"github.com/jigadhirasu/common/j"
	"github.com/jigadhirasu/common/jtype"
	"gorm.io/gorm"
)

type Sign struct {
	UserID    string
	UserType  string // iadmin iagent islot
	UserAgent string
	Address   string
	Region    string      `json:",omitempty"` // 區域
	City      string      `json:",omitempty"` // 城市
	Method    string      `json:",omitempty"` // in password otp
	Reason    *j.Error    `json:",omitempty"` // 錯誤訊息
	SignAt    *jtype.Time `json:",omitempty"`
	SignOut   *jtype.Time `json:",omitempty"`
}

func (a *Sign) UU(uuid ...string) string {
	return ""
}

func (Sign) TableName() string {
	return "signs"
}

type SignPack struct {
	UserID   string `gorm:"column:UserID; index; type:varchar(40) AS (COALESCE(JSON_VALUE(Doc, '$.UserID'), ''))"`
	UserType string `gorm:"column:UserType; type:varchar(40) AS (COALESCE(JSON_VALUE(Doc, '$.UserType'), ''))"`
	Method   string `gorm:"column:Method; type:varchar(40) AS (COALESCE(JSON_VALUE(Doc, '$.Method'), ''))"`
	Address  string `gorm:"column:Address; type:varchar(40) AS (COALESCE(JSON_VALUE(Doc, '$.Address'), ''))"`
	Region   string `gorm:"column:Region; type:varchar(40) AS (COALESCE(JSON_VALUE(Doc, '$.Region'), ''))"`
	City     string `gorm:"column:City; type:varchar(40) AS (COALESCE(JSON_VALUE(Doc, '$.City'), ''))"`
	Pack
	Query
}

func (SignPack) TableName() string {
	return Sign{}.TableName()
}

func (m SignPack) Where(tx *gorm.DB) *gorm.DB {
	if m.UserID != "" {
		tx = tx.Where("UserID = ?", m.UserID)
	}
	if m.UserType != "" {
		tx = tx.Where("UserType = ?", m.UserType)
	}
	return tx
}
