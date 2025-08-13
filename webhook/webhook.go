package webhook

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
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
	"time"

	"github.com/essentialkaos/ek/v13/hashutil"

	"github.com/essentialkaos/pachca"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// MAX_WEBHOOK_SIZE is the maximum size of the webhook payload
const MAX_WEBHOOK_SIZE = 1024 * 1024

const (
	EVENT_NEW    WebhookEvent = "new"
	EVENT_UPDATE WebhookEvent = "update"
	EVENT_DELETE WebhookEvent = "delete"

	EVENT_ADD    WebhookEvent = "add"
	EVENT_REMOVE WebhookEvent = "remove"

	EVENT_INVITE   WebhookEvent = "invite"
	EVENT_CONFIRM  WebhookEvent = "confirm"
	EVENT_SUSPEND  WebhookEvent = "suspend"
	EVENT_ACTIVATE WebhookEvent = "activate"

	EVENT_LINK_SHARED WebhookEvent = "link_shared"

	EVENT_SUBMIT WebhookEvent = "submit"
)

const (
	TYPE_MESSAGE        WebhookType = "message"
	TYPE_REACTION       WebhookType = "reaction"
	TYPE_BUTTON         WebhookType = "button"
	TYPE_CHAT_MEMBER    WebhookType = "chat_member"
	TYPE_COMPANY_MEMBER WebhookType = "company_member"
	TYPE_VIEW           WebhookType = "view"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// WebhookEvent is type for webhook events
type WebhookEvent string

// WebhookType is type for webhook types
type WebhookType string

// ////////////////////////////////////////////////////////////////////////////////// //

// Webhook is basic webhook object
type Webhook struct {
	Type      WebhookType `json:"type"`
	Timestamp int64       `json:"webhook_timestamp"`
}

// WebhookMessage contains payload of new message webhook
//
// https://crm.pachca.com/dev/getting-started/webhooks/#new-message
type Message struct {
	Webhook
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
}

// Reaction contains payload of reaction webhook
//
// https://crm.pachca.com/dev/getting-started/webhooks/#reaction
type Reaction struct {
	Webhook
	Event     WebhookEvent `json:"event"`
	UserID    uint         `json:"user_id"`
	MessageID uint         `json:"message_id"`
	CreatedAt pachca.Date  `json:"created_at"`
}

// Button contains payload of button webhook
//
// https://crm.pachca.com/dev/getting-started/webhooks/#button
type Button struct {
	Webhook
	Data      string `json:"data"`
	UserID    uint   `json:"user_id"`
	MessageID uint   `json:"message_id"`
	TriggerID string `json:"trigger_id"`
}

// Chat contains payload of chat members changes webhook
//
// https://crm.pachca.com/dev/getting-started/webhooks/#chat-member
type Chat struct {
	Webhook
	Event     WebhookEvent `json:"event"`
	ChatID    uint         `json:"chat_id"`
	ThreadID  uint         `json:"thread_id"`
	CreatedAt pachca.Date  `json:"created_at"`
	UserIDs   []uint       `json:"user_ids"`
}

// Organization contains payload of organization members changes webhook
//
// https://crm.pachca.com/dev/getting-started/webhooks/#company-member
type Organization struct {
	Webhook
	Event     WebhookEvent `json:"event"`
	UserIDs   []uint       `json:"user_ids"`
	CreatedAt pachca.Date  `json:"created_at"`
}

// Unfurl contains payload for generating URL preview
//
// https://crm.pachca.com/dev/messages/link_previews/
type Unfurl struct {
	Webhook
	Event     WebhookEvent  `json:"event"`
	ChatID    uint          `json:"chat_id"`
	MessageID uint          `json:"message_id"`
	Links     []*UnfurlLink `json:"links"`
}

// View contains payload from view form
type View struct {
	Event      WebhookEvent    `json:"event"`
	Metadata   string          `json:"private_metadata"`
	CallbackID string          `json:"callback_id"`
	UserID     uint            `json:"user_id"`
	Data       json.RawMessage `json:"data"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// WebhookThread contains info about message thread
type Thread struct {
	MessageID     uint `json:"message_id"`
	MessageChatID uint `json:"message_chat_id"`
}

// UnfurlLink contains info about link in message to unfurl
type UnfurlLink struct {
	URL    string `json:"url"`
	Domain string `json:"domain"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrNilWebhook  = errors.New("Webhook is nil")
	ErrEmptyData   = errors.New("Webhook has no data")
	ErrInvalidSig  = errors.New("Webhook has invalid signature")
	ErrNoSignature = errors.New("Webhook has no signature")
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Read reads webhook data from HTTP request
func Read(r *http.Request) (any, error) {
	rr := io.LimitReader(r.Body, MAX_WEBHOOK_SIZE)
	data, err := io.ReadAll(rr)

	if err != nil {
		return nil, fmt.Errorf("Can't read webhook data: %w", err)
	}

	return Decode(data)
}

// ReadSigned reads webhook data from HTTP request and validates signature
func ReadSigned(r *http.Request, secret string) (any, error) {
	if r.Header.Get("Pachca-Signature") == "" {
		return nil, ErrNoSignature
	}

	rr := io.LimitReader(r.Body, MAX_WEBHOOK_SIZE)
	data, err := io.ReadAll(rr)

	if err != nil {
		return nil, fmt.Errorf("Can't read webhook data: %w", err)
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, err = mac.Write(data)

	if err != nil {
		return nil, fmt.Errorf("Can't calculate webhook HMAC hash: %w", err)
	}

	if !hashutil.Sum(mac).EqualString(r.Header.Get("Pachca-Signature")) {
		return nil, ErrInvalidSig
	}

	return Decode(data)
}

// Decode unmarshals webhook JSON data
func Decode(data []byte) (any, error) {
	w := &Webhook{}
	err := json.Unmarshal(data, w)

	if err != nil {
		return nil, fmt.Errorf("Can't parse webhook JSON: %w", err)
	}

	var ww any

	switch w.Type {
	case TYPE_MESSAGE:
		ww = &Message{}
	case TYPE_REACTION:
		ww = &Reaction{}
	case TYPE_BUTTON:
		ww = &Button{}
	case TYPE_CHAT_MEMBER:
		ww = &Chat{}
	case TYPE_COMPANY_MEMBER:
		ww = &Organization{}
	case TYPE_VIEW:
		ww = &View{}
	default:
		return nil, fmt.Errorf("Unsupported webhook type %q", w.Type)
	}

	err = json.Unmarshal(data, ww)

	if err != nil {
		return nil, fmt.Errorf("Can't parse webhook JSON: %w", err)
	}

	return ww, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Is returns true if webhook has given type
func (w *Webhook) Is(typ WebhookType) bool {
	return w != nil && w.Type == typ
}

// Age returns age of webhook
func (w *Webhook) Age() time.Duration {
	if w == nil {
		return 0
	}

	return time.Since(time.Unix(w.Timestamp, 0))
}

// ////////////////////////////////////////////////////////////////////////////////// //

// UnmarshalData unmarshals webhook data
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
