package jmariadb

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/jigadhirasu/common/jcommon"
	"gorm.io/gorm"
)

// 有紀錄
func Create(tags jcommon.Tags, m Model) func(db *gorm.DB) jcommon.Bytes {
	r := &Record{
		OpType:   tags.String("OpType"),
		OpID:     tags.String("OpID"),
		Method:   "create",
		Target:   m.TableName(),
		TargetID: m.UU(),
		OpBefore: []byte(`{}`),
		OpAfter:  jcommon.JSON(m),
	}
	sh := &Pack{Doc: jcommon.JSON(m)}
	return func(db *gorm.DB) jcommon.Bytes {
		tx := db.Table(m.TableName()).Create(sh)
		if tx.Error != nil {
			return jcommon.JSON(Result{Error: jcommon.DBErr(tx.Error)})
		}
		db.Table(r.TableName()).Create(&Pack{Doc: jcommon.JSON(r)})
		return jcommon.JSON(Result{
			UUID:  m.UU(),
			Error: jcommon.DBErr(tx.Error),
		})
	}
}

// 沒紀錄
func Createx(tags jcommon.Tags, m Model) func(db *gorm.DB) jcommon.Bytes {
	sh := &Pack{Doc: jcommon.JSON(m)}
	return func(db *gorm.DB) jcommon.Bytes {
		tx := db.Table(m.TableName()).Create(sh)
		if tx.Error != nil {
			return jcommon.JSON(Result{Error: jcommon.DBErr(tx.Error)})
		}
		return jcommon.JSON(Result{
			UUID:  m.UU(),
			Error: jcommon.DBErr(tx.Error),
		})
	}
}

func Delete(tags jcommon.Tags, m Model) func(db *gorm.DB) jcommon.Bytes {
	r := &Record{
		OpType:   tags.String("OpType"),
		OpID:     tags.String("OpID"),
		Method:   "delete",
		Target:   m.TableName(),
		TargetID: m.UU(),
		OpAfter:  []byte(`{}`),
	}
	return func(db *gorm.DB) jcommon.Bytes {
		r.OpBefore = C{Table: m.TableName(), UUID: m.UU()}.Value(db)
		tx := db.Table(m.TableName()).Where("UUID = ?", m.UU()).Delete(m)
		db.Table(r.TableName()).Create(&Pack{Doc: jcommon.JSON(r)})
		return jcommon.JSON(Result{
			UUID:  m.UU(),
			Error: jcommon.DBErr(tx.Error),
		})
	}
}

// return Result
func Update(tags jcommon.Tags, m Model) func(db *gorm.DB) jcommon.Bytes {
	r := &Record{
		OpType:   tags.String("OpType"),
		OpID:     tags.String("OpID"),
		Method:   "update",
		Target:   m.TableName(),
		TargetID: m.UU(),
	}

	oldPtr := reflect.New(reflect.Indirect(reflect.ValueOf(m)).Type())
	o := oldPtr.Interface()
	return func(db *gorm.DB) jcommon.Bytes {
		// *** 檢查異動欄位 ***
		v := C{Table: m.TableName(), UUID: m.UU()}.Value(db)
		if len(v) < 1 {
			return jcommon.JSON(Result{Error: jcommon.NewErr(520, "data not found")})
		}
		v.Decode(o)
		SQL := fmt.Sprintf("UPDATE %s SET Doc = JSON_MERGE_PATCH(Doc, ?) WHERE UUID = ?", m.TableName())
		tx := db.Exec(SQL, jcommon.JSON(m), m.UU())
		if err := diff(o, m); err != nil {
			return jcommon.JSON(Result{Error: jcommon.DBErr(err)})
		}
		r.OpBefore = jcommon.JSON(o)
		r.OpAfter = jcommon.JSON(m)
		// *** 僅記錄有異動的欄位 ***
		if tx.RowsAffected < 1 {
			return jcommon.JSON(Result{
				Error: jcommon.NewErr(520, "no changed"),
			})
		}
		db.Table(r.TableName()).Create(&Pack{Doc: jcommon.JSON(r)})
		return jcommon.JSON(Result{
			RowsAffected: tx.RowsAffected,
			Error:        jcommon.DBErr(tx.Error),
		})
	}
}

// return Result
func Get(m Model) func(db *gorm.DB) jcommon.Bytes {
	return func(db *gorm.DB) jcommon.Bytes {
		b := C{Table: m.TableName(), UUID: m.UU()}.Value(db)
		result := Result{Data: b}
		if len(b) == 0 {
			result.Error = jcommon.DBErr(errors.New("not found~"))
		}
		return jcommon.JSON(result)
	}
}

// return Result
func Find(m IQuery) func(db *gorm.DB) jcommon.Bytes {
	return func(db *gorm.DB) jcommon.Bytes {
		tx := m.Where(db).Table(m.TableName())

		var count int64
		tx.Count(&count)

		bb := [][]byte{}
		tx = m.OrderBy(tx)
		tx = m.Limit(tx)
		tx = tx.Pluck(`JSON_MERGE_PATCH(Doc, JSON_OBJECT('CreatedAt',CreatedAt,'UpdatedAt',UpdatedAt)) as Doc`, &bb)

		b := bytes.Join(bb, []byte(","))
		return jcommon.JSON(Result{
			Data:   bytes.Join([][]byte{[]byte("["), b, []byte("]")}, []byte("")),
			Length: count,
			Error:  jcommon.DBErr(tx.Error),
		})
	}
}
