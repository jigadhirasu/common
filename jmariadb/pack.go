package jmariadb

import (
	"time"

	"github.com/jigadhirasu/common/j"
)

type Pack struct {
	Doc       j.Bytes   `gorm:"column:Doc;"`
	ID        int64     `gorm:"column:ID;"`
	CreatedAt time.Time `gorm:"column:CreatedAt; default:CURRENT_TIMESTAMP;"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt; default:CURRENT_TIMESTAMP;"`
}

type Model interface {
	TableName() string
	UU(uuid ...string) string
}
