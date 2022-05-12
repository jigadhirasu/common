package jmariadb

import "github.com/jigadhirasu/common/j"

type Result struct {
	UUID         string   `json:",omitempty"`
	RowsAffected int64    `json:",omitempty"`
	Error        *j.Error `json:",omitempty"`
	Length       int64    `json:",omitempty"`
	Data         j.Bytes  `json:",omitempty"`
}
