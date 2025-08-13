package block

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"testing"
	"time"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

type BlocksSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&BlocksSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *BlocksSuite) TestInit(c *C) {
	b1 := &Header{}
	b1.Init()
	c.Assert(b1.Type, Equals, "header")

	b2 := &PlainText{}
	b2.Init()
	c.Assert(b2.Type, Equals, "plain_text")

	b3 := &Markdown{}
	b3.Init()
	c.Assert(b3.Type, Equals, "markdown")

	b4 := &Divider{}
	b4.Init()
	c.Assert(b4.Type, Equals, "divider")

	b5 := &Input{}
	b5.Init()
	c.Assert(b5.Type, Equals, "input")

	b6 := &Select{}
	b6.Init()
	c.Assert(b6.Type, Equals, "select")

	b7 := &Radio{}
	b7.Init()
	c.Assert(b7.Type, Equals, "radio")

	b8 := &Checkbox{}
	b8.Init()
	c.Assert(b8.Type, Equals, "checkbox")

	b9 := &Date{}
	b9.Init()
	c.Assert(b9.Type, Equals, "date")

	b10 := &Time{}
	b10.Init()
	c.Assert(b10.Type, Equals, "time")

	b11 := &Files{}
	b11.Init()
	c.Assert(b11.Type, Equals, "file_input")
}

func (s *BlocksSuite) TestNil(c *C) {
	var bs *Select
	c.Assert(bs.AddOption("test", "Test1234", true), IsNil)

	var br *Radio
	c.Assert(br.AddOption("test", "Test1234", "Test description", true), IsNil)

	var bc *Checkbox
	c.Assert(bc.AddOption("test", "Test1234", "Test description", true), IsNil)

	var bd *Date
	c.Assert(bd.Set(2024, 2, 15), IsNil)
	c.Assert(bd.SetWithDate(time.Date(2025, 6, 20, 12, 30, 0, 0, time.Local)), IsNil)

	var bt *Time
	c.Assert(bt.Set(6, 45), IsNil)
	c.Assert(bt.SetWithDate(time.Date(2025, 6, 20, 12, 5, 0, 0, time.Local)), IsNil)
}

func (s *BlocksSuite) TestHelpers(c *C) {
	bs := &Select{}
	c.Assert(bs.AddOption("test", "Test1234", true), NotNil)
	c.Assert(bs.Options, HasLen, 1)
	c.Assert(bs.Options[0].Text, Equals, "test")
	c.Assert(bs.Options[0].Value, Equals, "Test1234")
	c.Assert(bs.Options[0].IsSelected, Equals, true)

	br := &Radio{}
	c.Assert(br.AddOption("test", "Test1234", "Test description", true), NotNil)
	c.Assert(br.Options, HasLen, 1)
	c.Assert(br.Options[0].Text, Equals, "test")
	c.Assert(br.Options[0].Value, Equals, "Test1234")
	c.Assert(br.Options[0].Description, Equals, "Test description")
	c.Assert(br.Options[0].IsSelected, Equals, true)

	bc := &Checkbox{}
	c.Assert(bc.AddOption("test", "Test1234", "Test description", true), NotNil)
	c.Assert(bc.Options, HasLen, 1)
	c.Assert(bc.Options[0].Text, Equals, "test")
	c.Assert(bc.Options[0].Value, Equals, "Test1234")
	c.Assert(bc.Options[0].Description, Equals, "Test description")
	c.Assert(bc.Options[0].IsSelected, Equals, true)

	bd := &Date{}
	c.Assert(bd.Set(2024, 2, 15), NotNil)
	c.Assert(bd.InitialValue, Equals, "2024-02-15")
	c.Assert(bd.SetWithDate(time.Date(2025, 6, 20, 12, 30, 0, 0, time.Local)), NotNil)
	c.Assert(bd.InitialValue, Equals, "2025-06-20")

	bt := &Time{}
	c.Assert(bt.Set(6, 45), NotNil)
	c.Assert(bt.InitialValue, Equals, "06:45")
	c.Assert(bt.SetWithDate(time.Date(2025, 6, 20, 12, 5, 0, 0, time.Local)), NotNil)
	c.Assert(bt.InitialValue, Equals, "12:05")
}
