package lib_time

import (
	"database/sql/driver"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

func (id IntTime) GetTime() time.Time {
	return id.Time
}

func (id IntTime) GetDate() time.Time {
	return timeToDate(id)
}

func (t *IntTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return t.UnmarshalBinary(v)
	case string:
		return t.UnmarshalText(v)
	case time.Time:
		*t = IntTime{v}
	case nil:
		t = nil
	default:
		return fmt.Errorf("cannot sql.Scan() IntTime from: %#v", v)
	}
	if t != nil {
		y, _, _ := t.Date()
		if y == 1 {
			t = nil
		}
	}
	return nil
}

func (t IntTime) Value() (driver.Value, error) {
	return driver.Value(t.Time.Format(DateTimeFormat)), nil
}

func (t *IntTime) UnmarshalJSON(bytes []byte) error {
	s := strings.Trim(string(bytes), `{`)
	s = strings.Trim(s, `}`)
	s = strings.ReplaceAll(s, "\"Time\":", "")
	s = strings.Trim(s, `"`)
	if s == "" || s == "null" || s == "{}" {
		return nil
	}

	nt, err := parseTime(s, false)
	if err != nil {
		return err
	}

	year, _, _ := nt.Date()
	if year != 1 {
		*t = IntTime{nt}
	}

	return nil
}

func (t *IntTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

	nt, err := parseTime(s, false)
	if err != nil {
		return err
	}

	year, _, _ := nt.Date()
	if year != 1 {
		*t = IntTime{nt}
	}

	return nil
}

func (t *IntTime) UnmarshalText(bytes string) error {
	s := strings.Trim(bytes, `{`)
	s = strings.Trim(s, `}`)
	s = strings.ReplaceAll(s, "\"Time\":", "")
	s = strings.Trim(s, `"`)
	if s == "" || s == "null" || s == "{}" {
		return nil
	}

	nt, err := parseTime(s, false)
	if err != nil {
		return err
	}

	year, _, _ := nt.Date()
	if year != 1 {
		*t = IntTime{nt}
	}
	return nil
}

func (t IntTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		b := make([]byte, 0)
		b = append(b, 'n')
		b = append(b, 'u')
		b = append(b, 'l')
		b = append(b, 'l')
		return b, nil
	}
	b := make([]byte, 0, len(DateTimeFormat)+2)
	b = append(b, '"')
	b = t.Time.AppendFormat(b, DateTimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t IntTime) MarshalText() ([]byte, error) {
	if t.IsZero() {
		b := make([]byte, 0)
		b = append(b, 'n')
		b = append(b, 'u')
		b = append(b, 'l')
		b = append(b, 'l')
		return b, nil
	}
	b := make([]byte, 0, len(DateTimeFormat))
	return t.Time.AppendFormat(b, DateTimeFormat), nil
}

func timeToDate(t IntTime) time.Time {
	year, month, day := t.GetTime().Local().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
