package helper

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// SQLTime is a custom type to handle PostgreSQL's time type
type SQLTime time.Time

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
