package helper

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// SQLTime is a custom type to handle PostgreSQL's time type
type SQLTime time.Time

// MarshalJSON implements the json.Marshaler interface
func (st SQLTime) MarshalJSON() ([]byte, error) {
	t := time.Time(st)
	return json.Marshal(t.Format("15:04:05"))
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (st *SQLTime) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}

	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return err
	}
	*st = SQLTime(t)
	return nil
}

// Scan implements the sql.Scanner interface
func (st *SQLTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		t, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		*st = SQLTime(t)
	case time.Time:
		*st = SQLTime(v)
	default:
		return fmt.Errorf("unsupported Scan, storing %T into *SQLTime", value)
	}
	return nil
}

// Value implements the driver.Valuer interface
func (st SQLTime) Value() (driver.Value, error) {
	return time.Time(st).Format("15:04:05"), nil
}

// String returns the time in HH:MM:SS format
func (st SQLTime) String() string {
	return time.Time(st).Format("15:04:05")
}
