package data

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type Payload struct {
	DateStart   Date    `json:"date_start"`
	DateEnd     Date    `json:"date_end"`
	RequestDoc  Files   `json:"request_doc"`
	Info        string  `json:"info"`
	Newsletters Options `json:"newsletters"`
	Time        Time    `json:"time"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

type DataSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&DataSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *DataSuite) TestUnmarshal(c *C) {
	data := `{
    "date_start": "2025-07-01",
    "date_end": "2025-07-14",
    "request_doc": [{
      "name": "request.png",
      "size": 19153,
      "url": "https://domain.com/request.png"
    }], 
    "info": "Поеду в сибирь на свадьбу лучшего друга",
    "newsletters": ["new_tasks", "project_updates"],
    "time": "16:38"
  }`

	v := &Payload{}
	err := json.Unmarshal([]byte(data), v)

	c.Assert(err, IsNil)
	c.Assert(v.DateStart.In(time.UTC).Unix(), Equals, int64(1751328000))
	c.Assert(v.DateEnd.In(time.UTC).Unix(), Equals, int64(1752451200))
	c.Assert(v.RequestDoc, HasLen, 1)
	c.Assert(v.RequestDoc[0].Name, Equals, "request.png")
	c.Assert(v.RequestDoc[0].URL, Equals, "https://domain.com/request.png")
	c.Assert(v.RequestDoc[0].Size, Equals, int64(19153))
	c.Assert(v.Info, Equals, "Поеду в сибирь на свадьбу лучшего друга")
	c.Assert([]string(v.Newsletters), DeepEquals, []string{"new_tasks", "project_updates"})
	c.Assert(v.Time.Hour, Equals, 16)
	c.Assert(v.Time.Minute, Equals, 38)
}

func (s *DataSuite) TestUnmarshalErrors(c *C) {
	d := &Date{}

	err := d.UnmarshalJSON([]byte(`null`))
	c.Assert(err, IsNil)

	err = d.UnmarshalJSON([]byte(`"ABCD-01-02"`))
	c.Assert(err, NotNil)

	t := &Time{}

	err = t.UnmarshalJSON([]byte(`null`))
	c.Assert(err, IsNil)
	err = t.UnmarshalJSON([]byte(`test`))
	c.Assert(err, NotNil)
	err = t.UnmarshalJSON([]byte(`AA:BB`))
	c.Assert(err, NotNil)
	err = t.UnmarshalJSON([]byte(`11:BB`))
	c.Assert(err, NotNil)
	err = t.UnmarshalJSON([]byte(`34:10`))
	c.Assert(err, NotNil)
	err = t.UnmarshalJSON([]byte(`12:70`))
	c.Assert(err, NotNil)
}

func (s *DataSuite) TestHelpers(c *C) {
	var t *Time
	c.Assert(t.String(), Equals, "00:00")

	t = &Time{6, 6}
	c.Assert(t.String(), Equals, "06:06")

	var o Options
	c.Assert(o.Has("test"), Equals, false)

	o = Options{"test1", "test2"}
	c.Assert(o.Has("test1"), Equals, true)

	var f Files
	c.Assert(f.Get(0), IsNil)

	f = Files{&File{"test.jpg", "https://domain.com/test.jpg", 11733}}
	c.Assert(f, HasLen, 1)
	c.Assert(f.Get(0), NotNil)
	c.Assert(f.Get(1), IsNil)
}
