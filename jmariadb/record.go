package jmariadb

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/jigadhirasu/jgorm/jcommon"
	"gorm.io/gorm"
)

type Record struct {
	OpType   string
	OpID     string
	Method   string
	Target   string
	TargetID string
	OpBefore jcommon.Bytes
	OpAfter  jcommon.Bytes
}

func (Record) TableName() string {
	return "records"
}

func (a *Record) UU(uuid ...string) string {
	return ""
}

type RecordPack struct {
	OpID     string `gorm:"column:OpID; index; type:varchar(40) AS (JSON_VALUE(Doc, '$.OpID'))"`
	TargetID string `gorm:"column:TargetID; index(target); type:varchar(40) AS (JSON_VALUE(Doc, '$.TargetID'))"`
	Target   string `gorm:"column:Target; index(target); type:varchar(40) AS (JSON_VALUE(Doc, '$.Target'))"`
	Method   string `gorm:"column:Method; type:varchar(40) AS (JSON_VALUE(Doc, '$.Method'))"`
	Pack
	Query
	Targets []string `gorm:"-"`
}

func (RecordPack) TableName() string {
	return "records"
}

func (m RecordPack) Where(tx *gorm.DB) *gorm.DB {
	if m.OpID != "" {
		tx = tx.Where("OpID = ?", m.OpID)
	}
	if m.TargetID != "" {
		tx = tx.Where("TargetID = ?", m.TargetID)
	}
	if m.Target != "" {
		tx = tx.Where("Target = ?", m.Target)
	}
	if m.Method != "" {
		tx = tx.Where("Method = ?", m.Method)
	}
	if len(m.Targets) > 0 {
		tx = tx.Where("Target IN ?", m.Targets)
	}
	return tx
}

func diff(old, new interface{}) error {

	o := reflect.ValueOf(old)
	n := reflect.ValueOf(new)
	if o.Kind() != reflect.Ptr {
		return jcommon.NewErr(484, "diff target is not ptr")
	}
	if o.Kind() != n.Kind() {
		return jcommon.NewErr(484, "diff target is not ptr")
	}
	diffStruct(o, n)
	return nil
}

func diffStruct(old, new reflect.Value) {
	o := reflect.Indirect(old)
	n := reflect.Indirect(new)

	for i := 0; i < n.NumField(); i++ {
		name := n.Type().Field(i).Name
		nf := n.Field(i)
		of := o.FieldByName(name)

		if !of.IsValid() {
			of.Set(reflect.Zero(nf.Type()))
		}

		if nf.Kind() == reflect.Struct {
			diffStruct(of, nf)
			continue
		}
		if nf.Kind() == reflect.Map {
			diffMap(of, nf)
			continue
		}

		if of.CanSet() {
			IsCompare := func() bool {
				switch nf.Interface().(type) {
				case []byte, json.RawMessage, jcommon.Bytes:
					return bytes.Equal(nf.Interface().(jcommon.Bytes), of.Interface().(jcommon.Bytes))
				case int, int8, int64, string, float64:
					return nf.Interface() == of.Interface()
				default:
					return bytes.Equal(jcommon.JSON(nf.Interface()), jcommon.JSON(of.Interface()))
				}
			}()
			if IsCompare {
				of.Set(reflect.Zero(of.Type()))
				nf.Set(reflect.Zero(nf.Type()))
			}
		}
	}
}

func diffMap(old, new reflect.Value) {
	o := reflect.Indirect(old)
	n := reflect.Indirect(new)
	// set := reflect.MakeMap(n.Type())
	iter := n.MapRange()
	for iter.Next() {
		key, nf := iter.Key(), iter.Value()
		of := o.MapIndex(key)

		if !of.IsValid() {
			of.Set(reflect.Zero(nf.Type()))
		}

		if nf.Kind() == reflect.Struct {
			diffStruct(of, nf)
			continue
		}
		if nf.Kind() == reflect.Map {
			diffMap(of, nf)
			continue
		}

		if of.CanSet() {
			IsCompare := func() bool {
				switch nf.Interface().(type) {
				case []byte, json.RawMessage, jcommon.Bytes:
					return bytes.Equal(nf.Interface().(jcommon.Bytes), of.Interface().(jcommon.Bytes))
				case []string:
					return bytes.Equal(jcommon.JSON(nf.Interface()), jcommon.JSON(of.Interface()))
				default:
					return nf.Interface() == of.Interface()
				}
			}()
			if IsCompare {
				of.Set(reflect.Zero(of.Type()))
				nf.Set(reflect.Zero(nf.Type()))
			}
		}
	}

}
