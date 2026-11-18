package webhook

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v14/hashutil"
	"github.com/essentialkaos/ek/v14/strutil"

	"github.com/essentialkaos/pachca"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	// EVENT_NEW identifies a creation event
	EVENT_NEW WebhookEvent = "new"

	// EVENT_UPDATE identifies an update event
	EVENT_UPDATE WebhookEvent = "update"

	// EVENT_DELETE identifies a deletion event
	EVENT_DELETE WebhookEvent = "delete"

	// EVENT_ADD identifies an addition event
	EVENT_ADD WebhookEvent = "add"

	// EVENT_REMOVE identifies a removal event
	EVENT_REMOVE WebhookEvent = "remove"

	// EVENT_INVITE identifies a workspace invitation
	EVENT_INVITE WebhookEvent = "invite"

	// EVENT_CONFIRM identifies a workspace invitation confirmation
	EVENT_CONFIRM WebhookEvent = "confirm"

	// EVENT_SUSPEND identifies a workspace member suspension
	EVENT_SUSPEND WebhookEvent = "suspend"

	// EVENT_ACTIVATE identifies a workspace member activation
	EVENT_ACTIVATE WebhookEvent = "activate"

	// EVENT_LINK_SHARED identifies a link-sharing event
	EVENT_LINK_SHARED WebhookEvent = "link_shared"

	// EVENT_SUBMIT identifies a form submission
	EVENT_SUBMIT WebhookEvent = "submit"
)

const (
	// TYPE_MESSAGE identifies a message webhook
	TYPE_MESSAGE WebhookType = "message"

	// TYPE_REACTION identifies a reaction webhook
	TYPE_REACTION WebhookType = "reaction"

	// TYPE_BUTTON identifies a button click webhook
	TYPE_BUTTON WebhookType = "button"

	// TYPE_CHAT_MEMBER identifies a chat membership webhook
	TYPE_CHAT_MEMBER WebhookType = "chat_member"

	// TYPE_ORG_MEMBER identifies a workspace membership webhook
	TYPE_ORG_MEMBER WebhookType = "company_member"

	// TYPE_VIEW identifies a form submission webhook
	TYPE_VIEW WebhookType = "view"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Webhook is the common interface for decoded Pachca webhooks
type Webhook interface {
	// Is reports whether the webhook has the given type
	Is(typ WebhookType) bool

	// Age returns the time elapsed since the webhook timestamp
	Age() time.Duration

	// GetType returns the webhook type
	GetType() WebhookType
}

// ////////////////////////////////////////////////////////////////////////////////// //

// WebhookEvent identifies the action reported by a webhook
type WebhookEvent string

// WebhookType identifies the kind of webhook payload
type WebhookType string

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic holds the webhook type and Unix timestamp shared by all payloads
type Basic struct {
	Type      WebhookType `json:"type"`
	Timestamp int64       `json:"webhook_timestamp"`
}

// Message holds the payload of a message webhook
//
// https://dev.pachca.com/guides/webhook#novye-soobscheniya
type Message struct {
	Basic
	MessageID       uint              `json:"id"`
	Event           WebhookEvent      `json:"event"`
	EntityType      pachca.EntityType `json:"entity_type"`
	EntityID        uint              `json:"entity_id"`
	Content         string            `json:"content"`
	UserID          uint              `json:"user_id"`
	ChatID          uint              `json:"chat_id"`
	ParentMessageID uint              `json:"parent_message_id"`
	CreatedAt       pachca.Date       `json:"created_at"`
	Thread          *Thread           `json:"thread"`
	Links           []*UnfurlLink     `json:"links"`
}

// Reaction holds the payload of a message reaction webhook
//
// https://dev.pachca.com/guides/webhook#reaktsii
type Reaction struct {
	Basic
	Event     WebhookEvent `json:"event"`
	UserID    uint         `json:"user_id"`
	ChatID    uint         `json:"chat_id"`
	MessageID uint         `json:"message_id"`
	CreatedAt pachca.Date  `json:"created_at"`
}

// Button holds the payload of a message button click
//
// https://dev.pachca.com/guides/webhook#nazhatie-knopok
type Button struct {
	Basic
	Data      string `json:"data"`
	UserID    uint   `json:"user_id"`
	ChatID    uint   `json:"chat_id"`
	MessageID uint   `json:"message_id"`
	TriggerID string `json:"trigger_id"`
}

// ChatMember holds the payload of a chat or thread membership change
//
// https://dev.pachca.com/guides/webhook#izmenenie-sostava-uchastnikov-chatov
type ChatMember struct {
	Basic
	Event     WebhookEvent `json:"event"`
	ChatID    uint         `json:"chat_id"`
	ThreadID  uint         `json:"thread_id"`
	CreatedAt pachca.Date  `json:"created_at"`
	UserIDs   []uint       `json:"user_ids"`
}

// OrgMember holds the payload of a workspace membership change
//
// https://dev.pachca.com/guides/webhook#izmenenie-sostava-uchastnikov-prostranstva
type OrgMember struct {
	Basic
	Event     WebhookEvent `json:"event"`
	UserIDs   []uint       `json:"user_ids"`
	CreatedAt pachca.Date  `json:"created_at"`
}

// View holds the payload of a form submission, including its raw JSON data
//
// https://dev.pachca.com/guides/webhook#zapolnenie-formy
type View struct {
	Basic
	Event      WebhookEvent    `json:"event"`
	Metadata   string          `json:"private_metadata"`
	CallbackID string          `json:"callback_id"`
	ChatID     uint            `json:"chat_id"`
	UserID     uint            `json:"user_id"`
	Data       json.RawMessage `json:"data"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Thread identifies a thread's parent message and chat
type Thread struct {
	MessageID     uint `json:"message_id"`
	MessageChatID uint `json:"message_chat_id"`
}

// UnfurlLink identifies a shared URL and its domain for link previews
type UnfurlLink struct {
	URL    string `json:"url"`
	Domain string `json:"domain"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	// ErrNilRequest indicates that the HTTP request or its body is nil
	ErrNilRequest = errors.New("request is nil")

	// ErrNilWebhook indicates that the webhook receiver is nil
	ErrNilWebhook = errors.New("webhook is nil")

	// ErrEmptyData indicates that the webhook contains no form data
	ErrEmptyData = errors.New("webhook has no data")

	// ErrInvalidSig indicates that the webhook signature did not match
	ErrInvalidSig = errors.New("webhook has invalid signature")

	// ErrNoSignature indicates that the request has no Pachca-Signature header
	ErrNoSignature = errors.New("webhook has no signature")
)

// MaxSize limits the number of request body bytes Read and ReadSigned consume
var MaxSize int64 = 1024 * 1024

// MaxAge is the maximum webhook age accepted by DecodeBytes, Read, and ReadSigned
var MaxAge = time.Minute

// ////////////////////////////////////////////////////////////////////////////////// //

// Read reads and decodes a webhook from an HTTP request without verifying its signature
func Read(r *http.Request) (Webhook, error) {
	if r == nil || r.Body == nil {
		return nil, ErrNilRequest
	}

	rr := io.LimitReader(r.Body, MaxSize)
	data, err := io.ReadAll(rr)

	if err != nil {
		return nil, fmt.Errorf("can't read webhook data: %w", err)
	}

	return DecodeBytes(data)
}

// ReadSigned verifies the request's Pachca-Signature header using HMAC-SHA256 before
// decoding the webhook
func ReadSigned(r *http.Request, secret string) (Webhook, error) {
	switch {
	case r == nil || r.Body == nil:
		return nil, ErrNilRequest
	case r.Header.Get("Pachca-Signature") == "":
		return nil, ErrNoSignature
	}

	rr := io.LimitReader(r.Body, MaxSize)
	data, err := io.ReadAll(rr)

	if err != nil {
		return nil, fmt.Errorf("can't read webhook data: %w", err)
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, err = mac.Write(data)

	if err != nil {
		return nil, fmt.Errorf("can't calculate webhook HMAC hash: %w", err)
	}

	if !hashutil.Sum(mac).EqualString(r.Header.Get("Pachca-Signature")) {
		return nil, ErrInvalidSig
	}

	return DecodeBytes(data)
}

// DecodeBytes decodes a JSON webhook and rejects it if it is older than MaxAge
func DecodeBytes(data []byte) (Webhook, error) {
	w, err := decodeWebhookData(data)

	if err != nil {
		return nil, err
	}

	if w.Age() > MaxAge {
		return nil, fmt.Errorf("webhook is too old (%s > %s)", w.Age(), MaxAge)
	}

	return w, nil
}

// DecodeJSON decodes a raw JSON webhook without checking its age
func DecodeJSON(data json.RawMessage) (Webhook, error) {
	return decodeWebhookData([]byte(data))
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Is reports whether the webhook has the given type
func (w *Basic) Is(typ WebhookType) bool {
	return w != nil && w.Type == typ
}

// Age returns the time elapsed since the webhook timestamp
func (w *Basic) Age() time.Duration {
	if w == nil {
		return 0
	}

	return time.Since(time.Unix(w.Timestamp, 0))
}

// GetType returns the webhook type
func (w *Basic) GetType() WebhookType {
	if w == nil {
		return ""
	}

	return w.Type
}

// String returns the webhook type as a string
func (w *Basic) String() string {
	return string(w.GetType())
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Command returns the slash command and its arguments, or an empty command and
// nil arguments if none is present
func (wh *Message) Command() (string, []string) {
	if wh == nil || wh.Content == "" {
		return "", nil
	}

	data := strings.Trim(wh.Content, " ")

	if data == "" || data[0] != '/' {
		return "", nil
	}

	cmd := strutil.Fields(data)

	if len(cmd) == 1 {
		return cmd[0], nil
	}

	return cmd[0], cmd[1:]
}

// ////////////////////////////////////////////////////////////////////////////////// //

// UnmarshalData decodes the form submission's JSON data into v
func (w *View) UnmarshalData(v any) error {
	switch {
	case w == nil:
		return ErrNilWebhook
	case len(w.Data) == 0:
		return ErrEmptyData
	}

	return json.Unmarshal(w.Data, v)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// decodeWebhookData decodes webhook data
func decodeWebhookData(data []byte) (Webhook, error) {
	w := &Basic{}
	err := json.Unmarshal(data, w)

	if err != nil {
		return nil, fmt.Errorf("can't parse webhook JSON: %w", err)
	}

	var ww Webhook

	switch w.Type {
	case TYPE_MESSAGE:
		ww = &Message{}
	case TYPE_REACTION:
		ww = &Reaction{}
	case TYPE_BUTTON:
		ww = &Button{}
	case TYPE_CHAT_MEMBER:
		ww = &ChatMember{}
	case TYPE_ORG_MEMBER:
		ww = &OrgMember{}
	case TYPE_VIEW:
		ww = &View{}
	default:
		return nil, fmt.Errorf("unsupported webhook type %q", w.Type)
	}

	json.Unmarshal(data, ww)

	return ww, nil
}
