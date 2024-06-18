package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Map map[string]interface{}

func (m Map) Value() (driver.Value, error) {
	j, err := json.Marshal(m)
	return string(j), err
}

func (m *Map) Scan(val interface{}) error {
	value, ok := val.([]byte)
	if !ok {
		return fmt.Errorf("Type assertion .([]byte) failed.")
	}
	return json.Unmarshal(value, m)
}

func (m *Map) ToStruct(s interface{}) error {
	buf, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, s)
}
