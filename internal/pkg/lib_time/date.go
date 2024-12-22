package lib_time

import (
	"database/sql/driver"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

func (id IntDate) GetTime() time.Time {
	return id.Time
}

func (id IntDate) GetDate() time.Time {
	return dateToTime(id)
}

func (t *IntDate) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return t.UnmarshalBinary(v)
	case string:
		return t.UnmarshalText(v)
	case time.Time:
		*t = IntDate{v}
	case nil:
		t = nil
	default:
		return fmt.Errorf("cannot sql.Scan() IntDate from: %#v", v)
	}
	if t != nil {
		y, _, _ := t.Date()
		if y == 1 {
			t = nil
		}
	}
	return nil
}

func (t IntDate) Value() (driver.Value, error) {
	return driver.Value(t.Time.Format(DateFormat)), nil
}

func (t *IntDate) UnmarshalJSON(bytes []byte) error {
	s := strings.Trim(string(bytes), `{`)
	s = strings.Trim(s, `}`)
	s = strings.ReplaceAll(s, "\"Time\":", "")
	s = strings.Trim(s, `"`)
	if s == "" || s == "null" || s == "{}" {
		return nil
	}

	nt, err := parseTime(s, true)
	if err != nil {
		return err
	}

	year, _, _ := nt.Date()
	if year != 1 {
		*t = IntDate{nt}
	}

	return nil
}

func (t *IntDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	s = strings.Trim(s, `{`)
	s = strings.Trim(s, `}`)
	s = strings.ReplaceAll(s, "\"Time\":", "")
	s = strings.Trim(s, `"`)
	if s == "" || s == "null" || s == "{}" {
		return nil
	}

	nt, err := parseTime(s, true)
	if err != nil {
		return err
	}

	year, _, _ := nt.Date()
	if year != 1 {
		*t = IntDate{nt}
	}

	return nil
}

func (t *IntDate) UnmarshalText(bytes string) error {
	s := strings.Trim(bytes, `{`)
	s = strings.Trim(s, `}`)
	s = strings.ReplaceAll(s, "\"Time\":", "")
	s = strings.Trim(s, `"`)
	if s == "" || s == "null" || s == "{}" {
		return nil
	}

	nt, err := parseTime(s, true)
	if err != nil {
		return err
	}

	year, _, _ := nt.Date()
	if year != 1 {
		*t = IntDate{nt}
	}
	return nil
}

func (t IntDate) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}

	b := make([]byte, 0, (len(DateFormat) + 2))

	b = append(b, '"')
	b = t.Time.AppendFormat(b, DateFormat)
	b = append(b, '"')

	return b, nil
}

func (t IntDate) MarshalText() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}

	b := make([]byte, 0, len(DateFormat))

	return t.Time.AppendFormat(b, DateFormat), nil
}

func dateToTime(t IntDate) time.Time {
	year, month, day := t.GetTime().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
