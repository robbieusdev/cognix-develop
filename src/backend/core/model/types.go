package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	buf, err := json.Marshal(j)
	return string(buf), err
}

func (j *JSONMap) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}
	err := json.Unmarshal(source, j)
	if err != nil {
		return err
	}
	if *j == nil {
		err = json.Unmarshal([]byte("{}"), j)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *JSONMap) ToStruct(dest interface{}) error {
	buf, err := json.Marshal(j)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, dest)
}

func (j *JSONMap) FromStruct(src interface{}) error {
	buf, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, j)
}

// StringSlice is a type that represents string slice.
type StringSlice []string

// Value implements the driver Valuer interface.
func (t StringSlice) Value() (driver.Value, error) {
	return pq.StringArray(t).Value()
}

// Scan implements the Scanner interface.
func (t *StringSlice) Scan(value interface{}) error {
	items := pq.StringArray{}
	if err := items.Scan(value); err != nil {
		return err
	}
	*t = []string(items)
	return nil
}

type JSON json.RawMessage

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	*j = bytes
	return nil
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return j, nil
}
