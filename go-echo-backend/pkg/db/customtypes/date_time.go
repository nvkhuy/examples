package customtypes

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const dateLayout = "2006-01-02"

type CustomDate time.Time

func (ct CustomDate) Format(layout string) string {
	return time.Time(ct).Format(layout)
}

func (ct CustomDate) Value() (driver.Value, error) {
	t, err := time.ParseInLocation(dateLayout, time.Time(ct).Format(dateLayout), time.Time(ct).Location())
	return t, err
}

func (ct *CustomDate) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return ct.UnmarshalText(string(v))
	case string:
		return ct.UnmarshalText(v)
	case time.Time:
		*ct = CustomDate(v)
	case nil:
		*ct = CustomDate{}
	default:
		return fmt.Errorf("cannot sql.Scan() CustomDate from: %#v", v)
	}

	return nil
}

func (ct *CustomDate) UnmarshalText(value string) error {
	dd, err := time.Parse(dateLayout, value)
	if err != nil {
		return err
	}
	*ct = CustomDate(dd)
	return nil
}

func (CustomDate) GormDataType() string {
	return "TIME"
}
