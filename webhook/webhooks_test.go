package webhook

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type ViewData struct {
	Info string `json:"info"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

type WebhookSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&WebhookSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *WebhookSuite) SetUpSuite(c *C) {
	MaxAge = 100000 * time.Hour
}

func (s *WebhookSuite) TestRead(c *C) {
	r, _ := http.NewRequest(
		"GET", "https://webhook.com",
		bytes.NewBufferString(`{"event":"new","type":"reaction","webhook_timestamp":1755117405}`),
	)

	w, err := Read(r)

	c.Assert(err, IsNil)
	c.Assert(w, NotNil)
}

func (s *WebhookSuite) TestReadSigned(c *C) {
	r, _ := http.NewRequest(
		"GET", "https://webhook.com",
		bytes.NewBufferString(`{"event":"new","type":"reaction","webhook_timestamp":1755117405}`),
	)

	r.Header.Set("Pachca-Signature", "158ea77bf7072f0d709df018429975c5df09f0485f0fc742c40b10fae3a6768f")

	w, err := ReadSigned(r, `pachca_wh_WDZVF7cqjjJNHfhkWRtP5ZtK6QNVJFFSalnBqKmwlvXaIs8tAujAiWua0ChJgDTt`)

	c.Assert(err, IsNil)
	c.Assert(w, NotNil)
}

func (s *WebhookSuite) TestDecode(c *C) {
	w, err := Decode([]byte(`{"event":"new","type":"message","webhook_timestamp":1755117405}`))
	c.Assert(err, IsNil)
	c.Assert(w, NotNil)
	c.Assert(w, FitsTypeOf, &Message{})
	c.Assert(w.Age(), Not(Equals), time.Duration(0))
	c.Assert(w.Is(TYPE_MESSAGE), Equals, true)
	c.Assert(w.GetType(), Equals, TYPE_MESSAGE)
	c.Assert(fmt.Sprint(w), Equals, string(TYPE_MESSAGE))

	w, err = Decode([]byte(`{"event":"new","type":"reaction","webhook_timestamp":1755117405}`))
	c.Assert(err, IsNil)
	c.Assert(w, NotNil)
	c.Assert(w, FitsTypeOf, &Reaction{})

	w, err = Decode([]byte(`{"type":"button","webhook_timestamp":1755117405}`))
	c.Assert(err, IsNil)
	c.Assert(w, NotNil)
	c.Assert(w, FitsTypeOf, &Button{})

	w, err = Decode([]byte(`{"event":"add","type":"chat_member","webhook_timestamp":1755117405}`))
	c.Assert(err, IsNil)
	c.Assert(w, NotNil)
	c.Assert(w, FitsTypeOf, &ChatMember{})

	w, err = Decode([]byte(`{"event":"invite","type":"company_member","webhook_timestamp":1755117405}`))
	c.Assert(err, IsNil)
	c.Assert(w, NotNil)
	c.Assert(w, FitsTypeOf, &OrgMember{})

	w, err = Decode([]byte(`{"event":"submit","type":"view","webhook_timestamp":1755117405}`))
	c.Assert(err, IsNil)
	c.Assert(w, NotNil)
	c.Assert(w, FitsTypeOf, &View{})
}

func (s *WebhookSuite) TestView(c *C) {
	w, err := Decode([]byte(`{"event":"submit","type":"view","webhook_timestamp":1755117405,"data":{"info":"Test info"}}`))

	c.Assert(err, IsNil)
	c.Assert(w, NotNil)

	d := &ViewData{}
	err = w.(*View).UnmarshalData(d)

	c.Assert(err, IsNil)
	c.Assert(d.Info, Equals, "Test info")
}

func (s *WebhookSuite) TestMessageCommand(c *C) {
	var wh1 *Message
	wh2 := &Message{Content: ""}
	wh3 := &Message{Content: "test"}
	wh4 := &Message{Content: "/test"}
	wh5 := &Message{Content: `/test "John Doe" 123`}

	cn, ca := wh1.Command()
	c.Assert(cn, Equals, "")
	c.Assert(ca, IsNil)

	cn, ca = wh2.Command()
	c.Assert(cn, Equals, "")
	c.Assert(ca, IsNil)

	cn, ca = wh3.Command()
	c.Assert(cn, Equals, "")
	c.Assert(ca, IsNil)

	cn, ca = wh4.Command()
	c.Assert(cn, Equals, "/test")
	c.Assert(ca, IsNil)

	cn, ca = wh5.Command()
	c.Assert(cn, Equals, "/test")
	c.Assert(ca, DeepEquals, []string{"John Doe", "123"})
}

func (s *WebhookSuite) TestErrors(c *C) {
	r, _ := http.NewRequest(
		"GET", "https://webhook.com",
		bytes.NewBufferString(`{"event":"new","type":"reaction","webhook_timestamp":1755117405}`),
	)

	_, err := Read(nil)
	c.Assert(err, Equals, ErrNilRequest)

	_, err = ReadSigned(nil, `test`)
	c.Assert(err, Equals, ErrNilRequest)
	_, err = ReadSigned(r, `test`)
	c.Assert(err, Equals, ErrNoSignature)

	r.Header.Set("Pachca-Signature", "0000")

	_, err = ReadSigned(r, `xxxx`)
	c.Assert(err, Equals, ErrInvalidSig)

	_, err = Decode([]byte("{"))
	c.Assert(err, ErrorMatches, `Can't parse webhook JSON: unexpected end of JSON input`)
	_, err = Decode([]byte(`{"event":"new","type":"unknown","webhook_timestamp":1755117405}`))
	c.Assert(err, ErrorMatches, `Unsupported webhook type "unknown"`)

	MaxAge = time.Second
	_, err = Decode([]byte(`{"event":"new","type":"message","webhook_timestamp":1755117405}`))
	c.Assert(err, ErrorMatches, `Webhook is too old .*`)
	MaxAge = 100000 * time.Hour

	var v *View
	d := &ViewData{}

	c.Assert(v.UnmarshalData(d), Equals, ErrNilWebhook)

	w, _ := Decode([]byte(`{"event":"submit","type":"view","webhook_timestamp":1755117405}`))
	c.Assert(w.(*View).UnmarshalData(d), Equals, ErrEmptyData)

	var ww *Basic
	c.Assert(ww.Is(TYPE_MESSAGE), Equals, false)
	c.Assert(ww.Age(), Equals, time.Duration(0))
	c.Assert(ww.GetType(), Equals, WebhookType(""))
	c.Assert(ww.String(), Equals, "")
}
