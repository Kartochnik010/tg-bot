package lib_time

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDisableDuplicateLayouts(t *testing.T) {
	dedup := map[string]bool{}
	for i := range layouts {
		dedup[layouts[i]] = true
	}

	/*for k := range dedup {
		t.Log(`"` + k + `",`)
	}*/
}

// var raw = `
// {
// 	"time1": "02.01.2006T15:04:05.000000000+06:00",
// 	"time2": "02.01.2006",
// 	"time3": "20060102150405",
// 	"time4": "22.06.2023 00:40:00",
// 	"time5": "2023-10-26T06:40:00.000+06:00",
// 	"time6": "24.08.2023T01:40:00.000000000+06:00",
// 	"time7": "25.09.2023T00:40:00",
// 	"time8": "26.10.2023 00:40:00"
// }`

// type strcutTC struct {
// 	Time1 IntTime    `json:"time1"`
// 	Time2 IntTime    `json:"time2"`
// 	Time3 IntDate    `json:"time3"`
// 	Time4 IntDate    `json:"time4"`
// 	Time5 IntTimeUTC `json:"time5"`
// 	Time6 IntDate    `json:"time6"`
// 	Time7 IntDate    `json:"time7"`
// 	Time8 IntDate    `json:"time8"`
// }

func TestXMLUnmarshal(t *testing.T) {
	type timeVars struct {
		XML     string      `json:"-"`
		StdTime *time.Time  `xml:"time.Time" json:"time.Time"`
		Time    *IntTime    `xml:"intTime" json:"intTime"`
		TimeUtc *IntTimeUTC `xml:"intTimeUtc" json:"intTimeUtc"`
		Date    *IntDate    `xml:"intDate" json:"intDate"`
	}

	tcs := []timeVars{
		{
			XML: `<?xml version="1.0" encoding="UTF-8"?>
			<root>
			   <intDate>2023-10-26T06:40:00.000+06:00</intDate>
			   <intTime>2023-10-26T06:40:00.000+06:00</intTime>
			   <intTimeUtc>2023-10-26T06:40:00.000+06:00</intTimeUtc>
			   <time.Time>2023-10-26T06:40:00.000+06:00</time.Time>
			</root>`,
		},
		{
			XML: `<?xml version="1.0" encoding="UTF-8"?>
			<root>
			   <intDate>2023-10-26</intDate>
			   <intTime>2023-10-26T06:40:00.000+06:00</intTime>
			   <intTimeUtc>2023-10-26T00:40:00Z</intTimeUtc>
			   <time.Time>2023-10-26T06:40:00+06:00</time.Time>
			</root>`,
		},
		{
			XML: `<?xml version="1.0" encoding="UTF-8"?>
			<root>
				<intDate>2023-10-26T06:40:00</intDate>
				<intTime>2023-10-26T06:40:00</intTime>
				<intTimeUtc>2023-10-26T06:40:00</intTimeUtc>
				<time.Time>2023-10-26T06:40:00.000+06:00</time.Time>
			</root>`,
		},
	}

	for i := range tcs {
		buf := bytes.NewBufferString(tcs[i].XML)
		err := xml.NewDecoder(buf).Decode(&tcs[i])
		// err := json.Unmarshal([]byte(tcs[i].Json), &tcs[i])
		if err != nil {
			t.Fatalf("%v, %v", tcs[i], err)
		}

		//		_ = json.NewDecoder().Decode(&tcs[i])

		//buf := &bytes.NewBuffer([]byte(tcs[i].Json)

		zn, zo := tcs[i].StdTime.Zone()
		t.Logf("%s - %s (%s, %d)\n", "StdTime", tcs[i].StdTime, zn, zo)
		zn, zo = tcs[i].Time.Zone()
		t.Logf("%s - %s (%s, %d)\n", "Time   ", tcs[i].Time, zn, zo)
		zn, zo = tcs[i].TimeUtc.Zone()
		t.Logf("%s - %s (%s, %d)\n", "TimeUtc", tcs[i].TimeUtc, zn, zo)
		zn, zo = tcs[i].Date.Zone()
		t.Logf("%s - %s (%s, %d)\n", "Date   ", tcs[i].Date, zn, zo)

		b, _ := xml.MarshalIndent(&tcs[i], " ", "    ")
		t.Logf("\n%s", b)
	}
}

func TestJsonUnmarshal(t *testing.T) {
	type timeVars struct {
		Json    string      `json:"-"`
		StdTime *time.Time  `json:"time.Time"`
		Time    *IntTime    `json:"intTime"`
		TimeUtc *IntTimeUTC `json:"intTimeUtc"`
		Date    *IntDate    `json:"intDate"`
	}

	tcs := []timeVars{
		{
			Json: `{"time.Time":"2023-10-26T06:40:00.000+06:00", "intTime":"2023-10-26T06:40:00.000+06:00", "intTimeUtc":"2023-10-26T06:40:00.000+06:00", "intDate":"2023-10-26T06:40:00.000+06:00"}`,
		},
		{
			Json: `{
				"time.Time": "2023-10-26T06:40:00+06:00",
				"intTime": "2023-10-26T06:40:00.000+06:00",
				"intTimeUtc": "2023-10-26T00:40:00Z",
				"intDate": "2023-10-26"
			}`,
		},
		{
			Json: `{"time.Time":"2023-10-26T06:40:00.000+06:00", "intTime":"2023-10-26T06:40:00", "intTimeUtc":"2023-10-26T06:40:00", "intDate":"2023-10-26T06:40:00"}`,
		},
	}

	for i := range tcs {
		buf := bytes.NewBufferString(tcs[i].Json)
		err := json.NewDecoder(buf).Decode(&tcs[i])
		// err := json.Unmarshal([]byte(tcs[i].Json), &tcs[i])
		if err != nil {
			t.Fatalf("%v, %v", tcs[i], err)
		}

		//		_ = json.NewDecoder().Decode(&tcs[i])

		//buf := &bytes.NewBuffer([]byte(tcs[i].Json)

		zn, zo := tcs[i].StdTime.Zone()
		t.Logf("%s - %s (%s, %d)\n", "StdTime", tcs[i].StdTime, zn, zo)
		zn, zo = tcs[i].Time.Zone()
		t.Logf("%s - %s (%s, %d)\n", "Time   ", tcs[i].Time, zn, zo)
		zn, zo = tcs[i].TimeUtc.Zone()
		t.Logf("%s - %s (%s, %d)\n", "TimeUtc", tcs[i].TimeUtc, zn, zo)
		zn, zo = tcs[i].Date.Zone()
		t.Logf("%s - %s (%s, %d)\n", "Date   ", tcs[i].Date, zn, zo)

		b, _ := json.MarshalIndent(&tcs[i], " ", "    ")
		t.Logf("\n%s", b)
	}
}

func TestToTime(t *testing.T) {
	ids, _ := time.Parse("2006-01-02", "1990-06-21")
	oops := As[IntTime](ToTime[IntTimeUTC](ids))
	assert.NotNil(t, oops)
}

func TestInitDate(t *testing.T) {
	ids, _ := time.Parse("2006-01-02", "1990-06-21")
	hf := StdFormat(ids, "DD.MM.YYYY")

	assert.Equal(t, hf, "21.06.1990")
	t.Log(hf)
}

func TestInitTime(t *testing.T) {
	ids, _ := time.Parse("2006-01-02 15:04:05", "1990-06-21 11:33:55")
	hf := StdFormat(ids, "DD.MM.YYYY hh:mm:ss")

	assert.Equal(t, hf, "21.06.1990 11:33:55")
	t.Log(hf)
}

func TestMarshal(t *testing.T) {
	type testCase struct {
		TimeUtc IntTimeUTC `json:"timeUtc"`
		Time    IntTime    `json:"time"`
		Date    IntDate    `json:"date"`
	}

	_t, _d, _tu := Now[IntTime](), Now[IntDate](), Now[IntTimeUTC]()

	tc := testCase{
		Time:    _t,
		Date:    _d,
		TimeUtc: _tu,
	}

	bts, err := json.Marshal(&tc)
	if assert.Nil(t, err) {
		t.Logf("%s", bts)
	}
}

func TestInitTimeUTC(t *testing.T) {
	ids, _ := time.Parse(DateTimeFormatUtc, "2023-11-30T16:02:06+06:00")
	hf := StdFormat(ids, "DD.MM.YYYY hh:mm:ss")

	assert.Equal(t, hf, "30.11.2023 16:02:06")
	t.Log(hf)
}

func TestNow(t *testing.T) {
	ids, err := time.Parse("2006-01-02T15:04:05Z07:00", "2023-03-09T01:08:23+06:00")
	if err != nil {
		t.Error(err)
	}

	type testCase struct {
		placeholder string
		validData   string
	}

	layouts := []testCase{
		{
			placeholder: "YY",
			validData:   "23",
		}, {
			placeholder: "YYYY",
			validData:   "2023",
		}, {
			placeholder: "M",
			validData:   "3",
		}, {
			placeholder: "MM",
			validData:   "03",
		}, {
			placeholder: "MMM",
			validData:   "Mar",
		}, {
			placeholder: "MMMM",
			validData:   "March",
		}, {
			placeholder: "D",
			validData:   "9",
		}, {
			placeholder: "DD",
			validData:   "09",
		}, {
			placeholder: "DDD",
			validData:   "Thu",
		}, {
			placeholder: "DDDD",
			validData:   "Thursday",
		}, {
			placeholder: "hh",
			validData:   "01",
		}, {
			placeholder: "h",
			validData:   "1",
		}, {
			placeholder: "mm",
			validData:   "08",
		}, {
			placeholder: "m",
			validData:   "8",
		}, {
			placeholder: "s",
			validData:   "23",
		}, {
			placeholder: "ss",
			validData:   "23",
		}, {
			placeholder: ".s",
			validData:   "",
		}, {
			placeholder: "z",
			validData:   "+06",
		}, {
			placeholder: "zz",
			validData:   "+06:00",
		}, {
			placeholder: "zzz",
			validData:   "+0600",
		}, {
			placeholder: "zzzz",
			validData:   "+06:00:00",
		}, {
			placeholder: "2006-01-02 15:04:05",
			validData:   "2023-03-09 01:08:23",
		},
	}

	for i := range layouts {
		std := StdFormat(ids, layouts[i].placeholder)
		t.Logf("%s -> %s", layouts[i], std)
		assert.Equal(t, std, layouts[i].validData)
	}
}

func TestLibNow(t *testing.T) {
	now := Now[IntTime]()

	b, _ := json.Marshal(&now)
	t.Logf("%v - %s", now, b)

	buf := &bytes.Buffer{}

	_ = json.NewEncoder(buf).Encode(now)
	t.Logf("%v - %s", now, buf.Bytes())

	var now2 *IntTime = &now

	b, _ = json.Marshal(now2)
	t.Logf("%v - %s", now, b)

	buf = &bytes.Buffer{}

	_ = json.NewEncoder(buf).Encode(now)
	t.Logf("%v - %s", now, buf.Bytes())
}
