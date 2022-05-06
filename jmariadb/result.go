package jmariadb

import "github.com/jigadhirasu/jgorm/jcommon"

type Result struct {
	UUID         string         `json:",omitempty"`
	RowsAffected int64          `json:",omitempty"`
	Error        *jcommon.Error `json:",omitempty"`
	Length       int64          `json:",omitempty"`
	Data         jcommon.Bytes  `json:",omitempty"`
}
