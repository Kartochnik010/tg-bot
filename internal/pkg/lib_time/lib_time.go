package lib_time

import (
	"errors"
	"strings"
	"time"
)

type (
	IntTime struct {
		time.Time
	}
	IntTimeIntf interface {
		IntTime
		GetTime() time.Time
	}

	IntDate struct {
		time.Time
	}
	IntDateIntf interface {
		IntDate
		GetTime() time.Time
	}

	IntTimeUTC struct {
		time.Time
	}
	IntTimeUTCIntf interface {
		IntDate
		GetTime() time.Time
	}

	IntDateTimeIntf interface {
		IntTime | IntDate | IntTimeUTC
		GetTime() time.Time
	}
)

func Now[T IntDateTimeIntf]() T {
	return T{time.Now()}
}

func As[F IntDateTimeIntf, T IntDateTimeIntf](dt T) F {
	return F(dt)
}

func Convert[T IntDateTimeIntf](dt T) time.Time {
	return dt.GetTime()
}

func ToTime[T IntDateTimeIntf](dt time.Time) T {
	return T{dt}
}

// # StdFormat helper for date time format
//
// # Year format
//
//	YYYY -> 2006 (Four-digit year)
//	YY -> 22 (Two-digit year)
//
// # Month format
//
//	MMMM -> January (Full month name)
//	MMM - > Jan (Three-letter abbreviation of the month)
//	MM -> 01 (Two-digit month (with a leading 0 if necessary))
//	M -> 1 (At most two-digit month (without a leading 0))
//
// # Day format
//
//	DDDD -> Monday (Full weekday name)
//	DDD -> Mon (Three-letter abbreviation of the weekday)
//	DD -> 02 (Two-digit month day (with a leading 0 if necessary))
//	D -> 2 (At most two-digit month day (without a leading 0))
//
// # Hour format
//
//	hh -> 15 (Two-digit 24h format hour)
//	h -> 3 (At most two-digit 12h format hour (without a leading 0))
//	ap/pm -> PM (AM/PM mark (uppercase))
//
// # Minute format
//
//	mm -> 04 (Two-digit minute (with a leading 0 if necessary))
//	m -> 4 (At most two-digit minute (without a leading 0))
//
// # Second format
//
//	ss -> 05 (Two-digit second (with a leading 0 if necessary))
//	s -> 5 (At most two-digit second (without a leading 0))
//	.s -> .99 (A fractional second (trailing zeros included))
//
// # Time zone format
//
//	zzzz -> Z07:00:00 (Numeric time zone offset with hours, minutes, and seconds separated by colons)
//	zzz -> Z0700 (Numeric time zone offset with hours and minutes)
//	zz -> Z07:00 (Numeric time zone offset with hours and minutes separated by colons)
//	z -> Z07 (Numeric time zone offset with hours)
func StdFormat(dt any, layout string) string {
	switch tp := dt.(type) {
	case time.Time:
		return tp.Format(rp.Replace(layout))
	case IntTime:
		return tp.GetTime().Format(rp.Replace(layout))
	case IntDate:
		return tp.GetTime().Format(rp.Replace(layout))
	default:
		return ""
	}
}

var rp = strings.NewReplacer(
	"YYYY", "2006",
	"YY", "06",

	"MMMM", "January",
	"MMM", "Jan",
	"MM", "01",
	"M", "1",

	"DDDD", "Monday",
	"DDD", "Mon",
	"DD", "02",
	"D", "2",

	"hh", "15",
	"h", "3",
	"am/pm", "PM",

	"mm", "04",
	"m", "4",

	"ss", "05",
	"s", "5",
	".s", ".99",

	"zzzz", "Z07:00:00",
	"zzz", "Z0700",
	"zz", "Z07:00",
	"z", "Z07",
)

func parseTime(s string, onlyDate bool) (nt time.Time, err error) {
	for i := range layouts {
		nt, err = time.Parse(layouts[i], s)
		if err == nil {
			switch {
			case onlyDate:
				return toDate(nt), nil
			default:
				return nt, nil
			}
		}
	}

	return time.Time{}, errors.New("error while parsing time: " + s)
}

const (
	DateTimeFormatUtc = time.RFC3339
	DateTimeFormat    = "2006-01-02T15:04:05.000Z07:00"
	DateFormat        = "2006-01-02"
)

var layouts = []string{
	time.Layout,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampNano,
	time.StampMicro,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,
	"2006-01-02T",
	"2006-01-02Z0700",
	"2006-01-02Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05Z0700",
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05.000",
	"2006-01-02T15:04:05.000Z0700",
	"2006-01-02T15:04:05.000Z07:00",
	"2006-01-02T15:04:05.000000Z0700",
	"2006-01-02T15:04:05.000000Z07:00",
	"2006-01-02T15:04:05.000000000Z0700",
	"2006-01-02T15:04:05.000000000Z07:00",
	"20060102150405",
	"02.01.2006",
	"02.01.2006 15:04:05",
	"02.01.2006T",
	"02.01.2006T15:04:05",
	"02.01.2006T15:04:05Z0700",
	"02.01.2006T15:04:05Z07:00",
	"02.01.2006T15:04:05.000000Z0700",
	"02.01.2006T15:04:05.000000Z07:00",
	"02.01.2006T15:04:05.000000000Z0700",
	"02.01.2006T15:04:05.000000000Z07:00",
	"Jan 2 2006 3:04PM",
	"02-01-2006",
	"2006.01.02 15:04:05",
}

func toDate(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
