package jmariadb

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/jigadhirasu/common/j"
	"gorm.io/gorm"
)

// 有紀錄
func Create(tags j.Tags, m Model) func(db *gorm.DB) j.Bytes {
	r := &Record{
		OpType:   tags.String("OpType"),
		OpID:     tags.String("OpID"),
		Method:   "create",
		Target:   m.TableName(),
		TargetID: m.UU(),
		OpBefore: []byte(`{}`),
		OpAfter:  j.JSON(m),
	}
	sh := &Pack{Doc: j.JSON(m)}
	return func(db *gorm.DB) j.Bytes {
		tx := db.Table(m.TableName()).Create(sh)
		if tx.Error != nil {
			return j.JSON(Result{Error: j.DBErr(tx.Error)})
		}
		db.Table(r.TableName()).Create(&Pack{Doc: j.JSON(r)})
		return j.JSON(Result{
			UUID:  m.UU(),
			Error: j.DBErr(tx.Error),
		})
	}
}

// 沒紀錄
func Createx(tags j.Tags, m Model) func(db *gorm.DB) j.Bytes {
	sh := &Pack{Doc: j.JSON(m)}
	return func(db *gorm.DB) j.Bytes {
		tx := db.Table(m.TableName()).Create(sh)
		if tx.Error != nil {
			return j.JSON(Result{Error: j.DBErr(tx.Error)})
		}
		return j.JSON(Result{
			UUID:  m.UU(),
			Error: j.DBErr(tx.Error),
		})
	}
}

func Delete(tags j.Tags, m Model) func(db *gorm.DB) j.Bytes {
	r := &Record{
		OpType:   tags.String("OpType"),
		OpID:     tags.String("OpID"),
		Method:   "delete",
		Target:   m.TableName(),
		TargetID: m.UU(),
		OpAfter:  []byte(`{}`),
	}
	return func(db *gorm.DB) j.Bytes {
		r.OpBefore = C{Table: m.TableName(), UUID: m.UU()}.Value(db)
		tx := db.Table(m.TableName()).Where("UUID = ?", m.UU()).Delete(m)
		db.Table(r.TableName()).Create(&Pack{Doc: j.JSON(r)})
		return j.JSON(Result{
			UUID:  m.UU(),
			Error: j.DBErr(tx.Error),
		})
	}
}

// return Result
func Update(tags j.Tags, m Model) func(db *gorm.DB) j.Bytes {
	r := &Record{
		OpType:   tags.String("OpType"),
		OpID:     tags.String("OpID"),
		Method:   "update",
		Target:   m.TableName(),
		TargetID: m.UU(),
	}

	oldPtr := reflect.New(reflect.Indirect(reflect.ValueOf(m)).Type())
	o := oldPtr.Interface()
	return func(db *gorm.DB) j.Bytes {
		// *** 檢查異動欄位 ***
		v := C{Table: m.TableName(), UUID: m.UU()}.Value(db)
		if len(v) < 1 {
			return j.JSON(Result{Error: j.NewErr(520, "data not found")})
		}
		v.Decode(o)
		SQL := fmt.Sprintf("UPDATE %s SET Doc = JSON_MERGE_PATCH(Doc, ?) WHERE UUID = ?", m.TableName())
		tx := db.Exec(SQL, j.JSON(m), m.UU())
		if err := diff(o, m); err != nil {
			return j.JSON(Result{Error: j.DBErr(err)})
		}
		r.OpBefore = j.JSON(o)
		r.OpAfter = j.JSON(m)
		// *** 僅記錄有異動的欄位 ***
		if tx.RowsAffected < 1 {
			return j.JSON(Result{
				Error: j.NewErr(520, "no changed"),
			})
		}
		db.Table(r.TableName()).Create(&Pack{Doc: j.JSON(r)})
		return j.JSON(Result{
			RowsAffected: tx.RowsAffected,
			Error:        j.DBErr(tx.Error),
		})
	}
}

// return Result
func Get(m Model) func(db *gorm.DB) j.Bytes {
	return func(db *gorm.DB) j.Bytes {
		b := C{Table: m.TableName(), UUID: m.UU()}.Value(db)
		result := Result{Data: b}
		if len(b) == 0 {
			result.Error = j.DBErr(errors.New("not found~"))
		}
		return j.JSON(result)
	}
}

// return Result
func Find(m IQuery) func(db *gorm.DB) j.Bytes {
	return func(db *gorm.DB) j.Bytes {
		tx := m.Where(db).Table(m.TableName())

		var count int64
		tx.Count(&count)

		bb := [][]byte{}
		tx = m.OrderBy(tx)
		tx = m.Limit(tx)
		tx = tx.Pluck(`JSON_MERGE_PATCH(Doc, JSON_OBJECT('CreatedAt',CreatedAt,'UpdatedAt',UpdatedAt)) as Doc`, &bb)

		b := bytes.Join(bb, []byte(","))
		return j.JSON(Result{
			Data:   bytes.Join([][]byte{[]byte("["), b, []byte("]")}, []byte("")),
			Length: count,
			Error:  j.DBErr(tx.Error),
		})
	}
}
