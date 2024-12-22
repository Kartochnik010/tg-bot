package lib_time

import (
	"database/sql/driver"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

func (id IntTimeUTC) GetTime() time.Time {
	return id.Time
}

func (id IntTimeUTC) GetDate() time.Time {
	return utcTimeToDate(id)
}

func (t *IntTimeUTC) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return t.UnmarshalBinary(v)
	case string:
		return t.UnmarshalText(v)
	case time.Time:
		*t = IntTimeUTC{v}
	case nil:
		t = nil
	default:
		return fmt.Errorf("cannot sql.Scan() IntTimeUTC from: %#v", v)
	}
	if t != nil {
		y, _, _ := t.Date()
		if y == 1 {
			t = nil
		}
	}
	return nil
}

func (t IntTimeUTC) Value() (driver.Value, error) {
	return driver.Value(t.Time.Format(DateTimeFormatUtc)), nil
}

func (t *IntTimeUTC) UnmarshalJSON(bytes []byte) error {
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
		*t = IntTimeUTC{nt}
	}

	return nil
}

func (t *IntTimeUTC) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
		*t = IntTimeUTC{nt}
	}

	return nil
}

func (t *IntTimeUTC) UnmarshalText(bytes string) error {
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
		*t = IntTimeUTC{nt}
	}
	return nil
}

func (t IntTimeUTC) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		b := make([]byte, 0)
		b = append(b, 'n')
		b = append(b, 'u')
		b = append(b, 'l')
		b = append(b, 'l')
		return b, nil
	}
	b := make([]byte, 0, len(DateTimeFormatUtc)+2)
	b = append(b, '"')
	b = t.Time.UTC().AppendFormat(b, DateTimeFormatUtc)
	b = append(b, '"')
	return b, nil
}

func (t IntTimeUTC) MarshalText() ([]byte, error) {
	if t.IsZero() {
		b := make([]byte, 0)
		b = append(b, 'n')
		b = append(b, 'u')
		b = append(b, 'l')
		b = append(b, 'l')
		return b, nil
	}
	b := make([]byte, 0, len(DateTimeFormatUtc))
	return t.Time.UTC().AppendFormat(b, DateTimeFormatUtc), nil
}

func utcTimeToDate(t IntTimeUTC) time.Time {
	year, month, day := t.GetTime().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
