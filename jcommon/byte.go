package jcommon

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func JSON(v interface{}) Bytes {
	jb, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return Bytes(jb)
}

type Bytes []byte

// MarshalJSON returns m as the JSON encoding of m.
func (b Bytes) MarshalJSON() ([]byte, error) {
	if b == nil {
		return []byte("null"), nil
	}
	return b, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (b *Bytes) UnmarshalJSON(data []byte) error {
	if b == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*b = append((*b)[0:0], data...)
	return nil
}

func (b Bytes) Decode(scan interface{}) error {
	return json.Unmarshal(b, scan)
}

func (b Bytes) String() string {
	return string(b)
}
func (b Bytes) Int() int {
	i, _ := strconv.Atoi(b.String())
	return i
}

// Value return json value, implement driver.Valuer interface
func (b Bytes) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (b *Bytes) Scan(value interface{}) error {
	if value == nil {
		*b = nil
		return nil
	}
	var bytes []byte
	switch b := value.(type) {
	case []byte:
		bytes = b
	case string:
		bytes = []byte(b)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	scan := []byte{}
	err := json.Unmarshal(bytes, &scan)
	*b = Bytes(scan)
	return err
}

// GormDataType gorm common data type
func (Bytes) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (Bytes) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}
