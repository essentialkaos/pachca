package pachca

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v13/errors"
	"github.com/essentialkaos/ek/v13/mathutil"
	"github.com/essentialkaos/ek/v13/path"
	"github.com/essentialkaos/ek/v13/req"
	"github.com/essentialkaos/ek/v13/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// API_URL is URL of Pachca API
const API_URL = "https://api.pachca.com/api/shared/v1"

// APP_URL is application URL used to generate links
const APP_URL = "https://app.pachca.com"

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	PROP_TYPE_DATE   PropertyType = "date"
	PROP_TYPE_LINK   PropertyType = "link"
	PROP_TYPE_NUMBER PropertyType = "number"
	PROP_TYPE_TEXT   PropertyType = "text"
)

const (
	INVITE_SENT      InviteStatus = "sent"
	INVITE_CONFIRMED InviteStatus = "confirmed"
)

const (
	ROLE_ADMIN       UserRole = "admin"
	ROLE_USER        UserRole = "user"
	ROLE_MULTI_GUEST UserRole = "multi_guest"
)

const (
	FILE_TYPE_FILE  FileType = "file"
	FILE_TYPE_IMAGE FileType = "image"
)

const (
	ENTITY_TYPE_DISCUSSION EntityType = "discussion"
	ENTITY_TYPE_THREAD     EntityType = "thread"
	ENTITY_TYPE_USER       EntityType = "user"
)

const (
	WEBHOOK_EVENT_NEW    WebhookEvent = "new"
	WEBHOOK_EVENT_UPDATE WebhookEvent = "update"
	WEBHOOK_EVENT_DELETE WebhookEvent = "delete"
)

const (
	WEBHOOK_TYPE_MESSAGE  WebhookType = "message"
	WEBHOOK_TYPE_REACTION WebhookType = "reaction"
	WEBHOOK_TYPE_BUTTON   WebhookType = "button"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Date is JSON date
type Date struct {
	time.Time
}

// ////////////////////////////////////////////////////////////////////////////////// //

// ID is default ID type
type ID uint

// EntityType is type of entity type
type EntityType string

// PropertyType is type for property type
type PropertyType string

// UserRole is type of user role
type UserRole string

// InviteStatus is type of invite status
type InviteStatus string

// FileType is type for file type
type FileType string

// WebhookEvent is type for webhook events
type WebhookEvent string

// WebhookType is type for webhook types
type WebhookType string

// ////////////////////////////////////////////////////////////////////////////////// //

// Chats is slice of chats
type Chats []*Chat

// Chat contains info about channel
type Chat struct {
	Members       []ID   `json:"member_ids"`
	GroupTagIDs   []ID   `json:"group_tag_ids"`
	ID            ID     `json:"id"`
	OwnerID       ID     `json:"owner_id"`
	Name          string `json:"name"`
	MeetRoomURL   string `json:"meet_room_url"`
	CreatedAt     Date   `json:"created_at"`
	LastMessageAt Date   `json:"last_message_at"`
	IsPublic      bool   `json:"public"`
	IsChannel     bool   `json:"channel"`
}

// Users is a slice of users
type Users []*User

// User contains info about user
type User struct {
	ID           ID           `json:"id"`
	CreatedAt    Date         `json:"created_at"`
	ImageURL     string       `json:"image_url"`
	Email        string       `json:"email"`
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	Nickname     string       `json:"nickname"`
	Role         UserRole     `json:"role"`
	PhoneNumber  string       `json:"phone_number"`
	TimeZone     string       `json:"time_zone"`
	Title        string       `json:"title"`
	InviteStatus InviteStatus `json:"invite_status"`
	Department   string       `json:"department"`
	Properties   Properties   `json:"custom_properties"`
	Tags         []string     `json:"list_tags"`
	Status       *Status      `json:"user_status"`
	IsBot        bool         `json:"bot"`
	IsSuspended  bool         `json:"suspended"`
}

// Status is user status
type Status struct {
	Emoji     string `json:"emoji"`
	Title     string `json:"title"`
	ExpiresAt Date   `json:"expires_at"`
}

// Properties is a slice of properties
type Properties []*Property

// Property is custom property
type Property struct {
	ID    ID           `json:"id"`
	Type  PropertyType `json:"data_type"`
	Name  string       `json:"name"`
	Value string       `json:"value"`
}

// Tag contains info about tag
type Tag struct {
	ID         ID     `json:"id"`
	Name       string `json:"name"`
	UsersCount int    `json:"users_count"`
}

// Tags is a slice of tags
type Tags []*Tag

// Reaction contains reaction info
type Reaction struct {
	UserID    ID     `json:"user_id"`
	CreatedAt Date   `json:"created_at"`
	Emoji     string `json:"code"`
}

// Reactions is a slice of reactions
type Reactions []*Reaction

// Thread contains info about thread
type Thread struct {
	ID            ID   `json:"id"`
	ChatID        ID   `json:"chat_id"`
	MessageID     ID   `json:"message_id"`
	MessageChatID ID   `json:"message_chat_id"`
	UpdatedAt     Date `json:"updated_at"`
}

// Message contains info about message
type Message struct {
	ID              ID          `json:"id"`
	EntityID        ID          `json:"entity_id"`
	ChatID          ID          `json:"chat_id"`
	ParentMessageID ID          `json:"parent_message_id"`
	UsedID          ID          `json:"user_id"`
	EntityType      EntityType  `json:"entity_type"`
	Content         string      `json:"content"`
	CreatedAt       Date        `json:"created_at"`
	Thread          *Thread     `json:"thread"`
	Files           Files       `json:"files"`
	Buttons         Buttons     `json:"buttons"`
	Forwarding      *Forwarding `json:"forwarding"`
}

// Forwarding contains info about message forwarding
type Forwarding struct {
	OriginalMessageID          ID   `json:"original_message_id"`
	OriginalChatID             ID   `json:"original_chat_id"`
	AuthorID                   ID   `json:"author_id"`
	OriginalThreadID           ID   `json:"original_thread_id"`
	OriginalThreadMessageID    ID   `json:"original_thread_message_id"`
	OriginalThreadParentChatID ID   `json:"original_thread_parent_chat_id"`
	OriginalCreatedAt          Date `json:"original_created_at"`
}

// File contains info about message attachment
type File struct {
	ID   ID       `json:"id,omitempty"`
	Key  string   `json:"key"`
	Name string   `json:"name"`
	Type FileType `json:"file_type"`
	URL  string   `json:"url,omitempty"`
	Size uint     `json:"size,omitempty"`
}

// Files is a slice of attachments
type Files []*File

// Button contains info about message button
type Button struct {
	Text string `json:"text"`
	URL  string `json:"url"`
	Data string `json:"data"`
}

// Buttons is a slice of buttons
type Buttons []*Button

// Upload contains upload info used for uploading files
type Upload struct {
	ContentDisposition string `json:"Content-Disposition"`
	ACL                string `json:"acl"`
	Policy             string `json:"policy"`
	Credential         string `json:"x-amz-credential"`
	Algorithm          string `json:"x-amz-algorithm"`
	Date               string `json:"x-amz-date"`
	Signature          string `json:"x-amz-signature"`
	Key                string `json:"key"`
	DirectURL          string `json:"direct_url"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// APIError contains API error info
type APIError struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Message    string `json:"message"`
	Code       string `json:"code"`
	StatusCode int
}

// ////////////////////////////////////////////////////////////////////////////////// //

// WebhookMessage is message webhook payload
type Webhook struct {
	Type            WebhookType  `json:"type"`
	ID              ID           `json:"id"`                // message
	Event           WebhookEvent `json:"event"`             // message, reaction
	EntityType      EntityType   `json:"entity_type"`       // message
	EntityID        ID           `json:"entity_id"`         // message
	Content         string       `json:"content"`           // message
	Emoji           string       `json:"code"`              // reaction
	Data            string       `json:"data"`              // button
	UserID          ID           `json:"user_id"`           // message, reaction
	CreatedAt       Date         `json:"created_at"`        // message, reaction, button
	ChatID          ID           `json:"chat_id"`           // message
	MessageID       ID           `json:"message_id"`        // reaction, button
	ParentMessageID ID           `json:"parent_message_id"` // message
	Thread          *Thread      `json:"thread"`            // message
}

// WebhookThread contains info about message thread
type WebhookThread struct {
	MessageID     ID `json:"message_id"`
	MessageChatID ID `json:"message_chat_id"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// uploadInfo contains info about uploaded file
type uploadInfo struct {
	Key         string        // Uploading key
	Name        string        // File name
	Size        uint          // File size
	ContentType string        // Content type
	Buffer      *bytes.Buffer // Buffer with data
}

// ////////////////////////////////////////////////////////////////////////////////// //

// ChatFilter is configuration for filtering chats
type ChatFilter struct {
	LastMessageAfter  time.Time
	LastMessageBefore time.Time
	Public            bool
}

// UserRequest is a struct with information needed to create or modify a user
type UserRequest struct {
	Email           string     `json:"email"`
	FirstName       string     `json:"first_name,omitempty"`
	LastName        string     `json:"last_name,omitempty"`
	Nickname        string     `json:"nickname,omitempty"`
	Role            UserRole   `json:"role,omitempty"`
	PhoneNumber     string     `json:"phone_number,omitempty"`
	Title           string     `json:"title,omitempty"`
	Department      string     `json:"department,omitempty"`
	Properties      Properties `json:"custom_properties,omitempty"`
	Tags            []string   `json:"list_tags,omitempty"`
	IsSuspended     bool       `json:"suspended,omitempty"`
	SkipEmailNotify bool       `json:"skip_email_notify,omitempty"`
}

// ChatRequest is a struct with information needed to create or modify a chat
type ChatRequest struct {
	Name        string `json:"name"`
	MemberIDs   []ID   `json:"member_ids,omitempty"`
	GroupTagIDs []ID   `json:"group_tag_ids,omitempty"`
	IsChannel   bool   `json:"channel,omitempty"`
	IsPublic    bool   `json:"public,omitempty"`
}

// MessageRequest is a struct with information needed to create or modify a message
type MessageRequest struct {
	EntityType         EntityType `json:"entity_type,omitempty"`
	EntityID           ID         `json:"entity_id"`
	Content            string     `json:"content"`
	Files              Files      `json:"files"`
	Buttons            Buttons    `json:"buttons,omitempty"`
	ParentMessageID    Buttons    `json:"parent_message_id,omitempty"`
	SkipInviteMentions bool       `json:"skip_invite_mentions,omitempty"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// UnmarshalJSON parses JSON date
func (d *Date) UnmarshalJSON(b []byte) error {
	data := string(b)

	if data == "null" {
		d.Time = time.Time{}
		return nil
	}

	date, err := time.Parse(`"2006-01-02T15:04:05.999Z"`, data)

	if err != nil {
		return err
	}

	d.Time = date

	return nil
}

// Error returns error text
func (e APIError) Error() string {
	return fmt.Sprintf("(%s) %s [%s:%s]", e.Code, e.Message, e.Key, e.Value)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// tokenValidationRegex is regex pattern for token validation
var tokenValidationRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_]{43}$`)

// s3ErrorExtractRegex is regex pattern for extracting text from S3 error message
var s3ErrorExtractRegex = regexp.MustCompile(`\<Message\>(.*)\<\/Message\>`)

var (
	ErrNilClient         = errors.New("Client is nil")
	ErrNilUserRequest    = errors.New("User requests is nil")
	ErrNilChatRequest    = errors.New("Chat requests is nil")
	ErrNilMessageRequest = errors.New("Message requests is nil")
	ErrNilProperty       = errors.New("Property requests is nil")
	ErrEmptyToken        = errors.New("Token is empty")
	ErrEmptyTag          = errors.New("Group tag is empty")
	ErrEmptyMessage      = errors.New("Message text is empty")
	ErrEmptyUserEmail    = errors.New("User email is required for creating user account")
	ErrEmptyChatName     = errors.New("Name is required for creating new chat")
	ErrEmptyUsersIDS     = errors.New("Users IDs are empty")
	ErrEmptyTagsIDS      = errors.New("Tags IDs are empty")
	ErrEmptyFilePath     = errors.New("Path to file is empty")
	ErrInvalidToken      = errors.New("Token is has wrong format")
	ErrInvalidMessageID  = errors.New("Message ID must be greater than 0")
	ErrInvalidChatID     = errors.New("Chat ID must be greater than 0")
	ErrInvalidUserID     = errors.New("User ID must be greater than 0")
	ErrInvalidThreadID   = errors.New("Thread ID must be greater than 0")
	ErrInvalidTagID      = errors.New("Group tag ID must be greater than 0")
	ErrInvalidEntityID   = errors.New("Entity ID must be greater than 0")
	ErrBlankEmoji        = errors.New("Non-blank emoji is required")
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Client is Pachca API client
type Client struct {
	BatchSize   int   // BatchSize is a number of items for paginated requests
	MaxFileSize int64 // Maximum file size to upload

	engine *req.Engine
	token  string
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewClient creates new client with given token
func NewClient(token string) (*Client, error) {
	err := ValidateToken(token)

	if err != nil {
		return nil, err
	}

	return &Client{
		BatchSize:   50,
		MaxFileSize: 10 * 1024 * 1024, // 10 MB

		token:  token,
		engine: &req.Engine{},
	}, nil
}

// ValidateToken validates API access token
func ValidateToken(token string) error {
	switch {
	case token == "":
		return ErrEmptyToken
	case !tokenValidationRegex.MatchString(token):
		return ErrInvalidToken
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// SetUserAgent sets user-agent info
func (c *Client) SetUserAgent(app, ver string) {
	if c == nil || c.engine == nil {
		return
	}

	c.engine.SetUserAgent(app, ver, "EK-Pachca.go/1")
}

// Engine returns pointer to request engine used for all HTTP requests to API
func (c *Client) Engine() *req.Engine {
	if c == nil || c.engine == nil {
		return nil
	}

	return c.engine
}

// CUSTOM PROPERTIES //////////////////////////////////////////////////////////////// //

// GetProperties returns custom properties
//
// https://crm.pachca.com/dev/common/fields/
func (c *Client) GetProperties() (Properties, error) {
	if c == nil || c.engine == nil {
		return nil, ErrNilClient
	}

	query := req.Query{"entity_type": "User"}

	resp := &struct {
		Data Properties `json:"data"`
	}{}

	err := c.sendRequest(
		req.GET, getURL("/custom_properties"),
		query, nil, resp,
	)

	if err != nil {
		return nil, fmt.Errorf("Can't fetch custom properties: %w", err)
	}

	return resp.Data, nil
}

// REACTIONS //////////////////////////////////////////////////////////////////////// //

// GetReactions returns slice with reactions added to given message
//
// https://crm.pachca.com/dev/reactions/list/
func (c *Client) GetReactions(messageID ID) (Reactions, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case messageID == 0:
		return nil, ErrInvalidMessageID
	}

	var result Reactions

	query := req.Query{"per": c.getBatchSize()}

	for i := 1; i < 100; i++ {
		query["page"] = i

		resp := &struct {
			Data Reactions `json:"data"`
		}{}

		err := c.sendRequest(
			req.GET, getURL("/messages/%d/reactions", messageID),
			query, nil, resp,
		)

		if err != nil {
			return nil, fmt.Errorf("Can't fetch reactions for message %d: %w", messageID, err)
		}

		result = append(result, resp.Data...)

		if len(resp.Data) != c.getBatchSize() {
			break
		}
	}

	return result, nil
}

// AddReaction adds given emoji reaction to the message
//
// https://crm.pachca.com/dev/reactions/new/
func (c *Client) AddReaction(messageID ID, emoji string) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case messageID == 0:
		return ErrInvalidMessageID
	case emoji == "":
		return ErrBlankEmoji
	}

	err := c.sendRequest(
		req.POST, getURL("/messages/%d/reactions", messageID),
		req.Query{"code": emoji}, nil, nil,
	)

	if err != nil {
		return fmt.Errorf("Can't add reaction %q to message %d: %w", emoji, messageID, err)
	}

	return nil
}

// DeleteReaction removes given emoji reaction from the message
//
// https://crm.pachca.com/dev/reactions/delete/
func (c *Client) DeleteReaction(messageID ID, emoji string) error {
	switch {
	case c == nil:
		return ErrNilClient
	case messageID == 0:
		return ErrInvalidMessageID
	case emoji == "":
		return ErrBlankEmoji
	}

	err := c.sendRequest(
		req.DELETE, getURL("/messages/%d/reactions", messageID),
		req.Query{"code": emoji}, nil, nil,
	)

	if err != nil {
		return fmt.Errorf("Can't remove reaction %q from message %d: %w", emoji, messageID, err)
	}

	return nil
}

// USERS //////////////////////////////////////////////////////////////////////////// //

// GetUser returns info about user
//
// https://crm.pachca.com/dev/users/get/
func (c *Client) GetUser(userID ID) (*User, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case userID == 0:
		return nil, ErrInvalidUserID
	}

	resp := &struct {
		Data *User `json:"data"`
	}{}

	err := c.sendRequest(
		req.GET, getURL("/users/%d", userID),
		nil, nil, resp,
	)

	if err != nil {
		return nil, fmt.Errorf("Can't fetch user info: %w", err)
	}

	return resp.Data, nil
}

// GetUsers returns info about all users
//
// https://crm.pachca.com/dev/users/list/
func (c *Client) GetUsers(searchQuery ...string) (Users, error) {
	if c == nil || c.engine == nil {
		return nil, ErrNilClient
	}

	var result Users

	query := req.Query{"per": c.getBatchSize()}

	if len(searchQuery) != 0 {
		query["query"] = searchQuery[0]
	}

	for i := 1; i < 100; i++ {
		query["page"] = i

		resp := &struct {
			Data Users `json:"data"`
		}{}

		err := c.sendRequest(req.GET, getURL("/users"), query, nil, resp)

		if err != nil {
			return nil, fmt.Errorf("Can't fetch users: %w", err)
		}

		result = append(result, resp.Data...)

		if len(resp.Data) != c.getBatchSize() {
			break
		}
	}

	return result, nil
}

// AddUser creates a new user
//
// https://crm.pachca.com/dev/users/new/
func (c *Client) AddUser(user *UserRequest) (*User, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case user == nil:
		return nil, ErrNilUserRequest
	case user.Email == "":
		return nil, ErrEmptyUserEmail
	}

	payload := &struct {
		User *UserRequest `json:"user"`
	}{
		User: user,
	}

	resp := &struct {
		Data *User `json:"data"`
	}{}

	err := c.sendRequest(req.POST, getURL("/users"), nil, payload, resp)

	if err != nil {
		return nil, fmt.Errorf("Can't create a new user: %w", err)
	}

	return resp.Data, nil
}

// EditUser modifies an existing user
//
// https://crm.pachca.com/dev/users/update/
func (c *Client) EditUser(userID ID, user *UserRequest) (*User, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case userID == 0:
		return nil, ErrInvalidUserID
	case user == nil:
		return nil, ErrNilUserRequest
	}

	payload := &struct {
		User *UserRequest `json:"user"`
	}{
		User: user,
	}

	resp := &struct {
		Data *User `json:"data"`
	}{}

	err := c.sendRequest(req.PUT, getURL("/users/%d", userID), nil, payload, resp)

	if err != nil {
		return nil, fmt.Errorf("Can't edit user %d: %w", userID, err)
	}

	return resp.Data, nil
}

// DeleteUser deletes an existing user
//
// https://crm.pachca.com/dev/users/delete/
func (c *Client) DeleteUser(userID ID) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case userID == 0:
		return ErrInvalidUserID
	}

	err := c.sendRequest(req.DELETE, getURL("/users/%d", userID), nil, nil, nil)

	if err != nil {
		return fmt.Errorf("Can't delete user %d: %w", userID, err)
	}

	return nil
}

// GROUP TAGS /////////////////////////////////////////////////////////////////////// //

// GetTags returns all group tags
//
// https://crm.pachca.com/dev/group_tags/list/
func (c *Client) GetTags() (Tags, error) {
	if c == nil || c.engine == nil {
		return nil, ErrNilClient
	}

	var result Tags

	query := req.Query{"per": c.getBatchSize()}

	for i := 1; i < 100; i++ {
		query["page"] = i

		resp := &struct {
			Data Tags `json:"data"`
		}{}

		err := c.sendRequest(req.GET, getURL("/group_tags"), query, nil, resp)

		if err != nil {
			return nil, fmt.Errorf("Can't fetch group tags: %w", err)
		}

		result = append(result, resp.Data...)

		if len(resp.Data) != c.getBatchSize() {
			break
		}
	}

	return result, nil
}

// GetTag returns info about group tag with given ID
//
// https://crm.pachca.com/dev/group_tags/get/
func (c *Client) GetTag(groupTagID ID) (*Tag, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case groupTagID == 0:
		return nil, ErrInvalidTagID
	}

	resp := &struct {
		Data *Tag `json:"data"`
	}{}

	err := c.sendRequest(
		req.GET, getURL("/group_tags/%d", groupTagID),
		nil, nil, resp,
	)

	if err != nil {
		return nil, fmt.Errorf("Can't fetch group tag: %w", err)
	}

	return resp.Data, nil
}

// GetTagUsers returns slice with users with given tag
//
// https://crm.pachca.com/dev/group_tags/users/
func (c *Client) GetTagUsers(groupTagID ID) (Users, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case groupTagID == 0:
		return nil, ErrInvalidTagID
	}

	var result Users

	query := req.Query{"per": c.getBatchSize()}

	for i := 1; i < 100; i++ {
		query["page"] = i

		resp := &struct {
			Data Users `json:"data"`
		}{}

		err := c.sendRequest(
			req.GET, getURL("/group_tags/%d/users", groupTagID),
			query, nil, resp,
		)

		if err != nil {
			return nil, fmt.Errorf("Can't fetch group tag users: %w", err)
		}

		result = append(result, resp.Data...)

		if len(resp.Data) != c.getBatchSize() {
			break
		}
	}

	return result, nil
}

// AddTag creates new group tag
//
// https://crm.pachca.com/dev/group_tags/new/
func (c *Client) AddTag(groupTagName string) (*Tag, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case groupTagName == "":
		return nil, ErrEmptyTag
	}

	payload := &struct {
		Name string `json:"name"`
	}{
		Name: groupTagName,
	}

	resp := &struct {
		Data *Tag `json:"data"`
	}{}

	err := c.sendRequest(req.POST, getURL("/group_tags"), nil, payload, resp)

	if err != nil {
		return nil, fmt.Errorf("Can't create new group tag %q: %w", groupTagName, err)
	}

	return resp.Data, nil
}

// EditTag changes name of given group tag
//
// https://crm.pachca.com/dev/group_tags/update/
func (c *Client) EditTag(groupTagID ID, groupTagName string) (*Tag, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case groupTagID == 0:
		return nil, ErrInvalidTagID
	case groupTagName == "":
		return nil, ErrEmptyTag
	}

	payload := &struct {
		Name string `json:"name"`
	}{
		Name: groupTagName,
	}

	resp := &struct {
		Data *Tag `json:"data"`
	}{}

	err := c.sendRequest(
		req.PUT, getURL("/group_tags/%d", groupTagID),
		nil, payload, resp,
	)

	if err != nil {
		return nil, fmt.Errorf("Can't edit group tag %d: %w", groupTagID, err)
	}

	return resp.Data, nil
}

// DeleteTag removes group tag
//
// https://crm.pachca.com/dev/group_tags/delete/
func (c *Client) DeleteTag(groupTagID ID) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case groupTagID == 0:
		return ErrInvalidTagID
	}

	err := c.sendRequest(
		req.DELETE, getURL("/group_tags/%d", groupTagID),
		nil, nil, nil,
	)

	if err != nil {
		return fmt.Errorf("Can't delete group tag %d: %w", groupTagID, err)
	}

	return nil
}

// CHATS //////////////////////////////////////////////////////////////////////////// //

// GetChats returns all chats and conversations
//
// https://crm.pachca.com/dev/chats/list/
func (c *Client) GetChats(filter ...ChatFilter) (Chats, error) {
	if c == nil || c.engine == nil {
		return nil, ErrNilClient
	}

	var result Chats
	var query req.Query

	if len(filter) == 0 {
		query = req.Query{"per": c.getBatchSize()}
	} else {
		query = filter[0].ToQuery()
		query["per"] = c.getBatchSize()
	}

	for i := 1; i < 100; i++ {
		query["page"] = i

		resp := &struct {
			Data Chats `json:"data"`
		}{}

		err := c.sendRequest(req.GET, getURL("/chats"), query, nil, resp)

		if err != nil {
			return nil, fmt.Errorf("Can't fetch chats: %w", err)
		}

		result = append(result, resp.Data...)

		if len(resp.Data) != c.getBatchSize() {
			break
		}
	}

	return result, nil
}

// GetChats returns info about specific channel
//
// https://crm.pachca.com/dev/chats/get/
func (c *Client) GetChat(chatID ID) (*Chat, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case chatID == 0:
		return nil, ErrInvalidChatID
	}

	resp := &struct {
		Data *Chat `json:"data"`
	}{}

	err := c.sendRequest(req.GET, getURL("/chats/%d", chatID), nil, nil, resp)

	if err != nil {
		return nil, fmt.Errorf("Can't fetch user info: %w", err)
	}

	return resp.Data, nil
}

// AddChat creates new chat
//
// https://crm.pachca.com/dev/chats/new/
func (c *Client) AddChat(chat *ChatRequest) (*Chat, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case chat == nil:
		return nil, ErrNilChatRequest
	case chat.Name == "":
		return nil, ErrEmptyChatName
	}

	payload := &struct {
		Chat *ChatRequest `json:"chat"`
	}{
		Chat: chat,
	}

	resp := &struct {
		Data *Chat `json:"data"`
	}{}

	err := c.sendRequest(req.POST, getURL("/chats"), nil, payload, resp)

	if err != nil {
		return nil, fmt.Errorf("Can't create a new chat %q: %w", chat.Name, err)
	}

	return resp.Data, nil
}

// EditChat modifies chat
//
// https://crm.pachca.com/dev/chats/new/
func (c *Client) EditChat(chatID ID, chat *ChatRequest) (*Chat, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case chatID == 0:
		return nil, ErrInvalidChatID
	case chat == nil:
		return nil, ErrNilChatRequest
	}

	payload := &struct {
		Chat *ChatRequest `json:"chat"`
	}{
		Chat: chat,
	}

	resp := &struct {
		Data *Chat `json:"data"`
	}{}

	err := c.sendRequest(req.PUT, getURL("/chats/%d", chatID), nil, payload, resp)

	if err != nil {
		return nil, fmt.Errorf("Can't modify chat %d: %w", chatID, err)
	}

	return resp.Data, nil
}

// AddChatUsers adds users with given IDs to the chat
//
// https://crm.pachca.com/dev/members/users/new/
func (c *Client) AddChatUsers(chatID ID, membersIDs []ID, silent bool) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case chatID == 0:
		return ErrInvalidChatID
	case len(membersIDs) == 0:
		return ErrEmptyUsersIDS
	}

	var query req.Query

	if silent {
		query = req.Query{"silent": silent}
	}

	payload := &struct {
		IDs []ID `json:"member_ids"`
	}{
		IDs: membersIDs,
	}

	err := c.sendRequest(
		req.PUT, getURL("/chats/%d/members", chatID),
		query, payload, nil,
	)

	if err != nil {
		return fmt.Errorf("Can't add users to chat %d: %w", chatID, err)
	}

	return nil
}

// AddChatTags adds group tags to the chat
//
// https://crm.pachca.com/dev/members/tags/new/
func (c *Client) AddChatTags(chatID ID, tagIDs []ID) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case chatID == 0:
		return ErrInvalidChatID
	case len(tagIDs) == 0:
		return ErrEmptyTagsIDS
	}

	payload := &struct {
		IDs []ID `json:"group_tag_ids"`
	}{
		IDs: tagIDs,
	}

	err := c.sendRequest(
		req.PUT, getURL("/chats/%d/group_tags", chatID),
		nil, payload, nil,
	)

	if err != nil {
		return fmt.Errorf("Can't add group tags to chat %d: %w", chatID, err)
	}

	return nil
}

// ExcludeChatUser excludes the user from the chat
//
// https://crm.pachca.com/dev/members/users/delete/
func (c *Client) ExcludeChatUser(chatID, userID ID) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case chatID == 0:
		return ErrInvalidChatID
	case userID == 0:
		return ErrInvalidUserID
	}

	err := c.sendRequest(
		req.DELETE, getURL("/chats/%d/members/%d", chatID, userID),
		nil, nil, nil,
	)

	if err != nil {
		return fmt.Errorf(
			"Can't exclude user %d from chat %d: %w",
			userID, chatID, err,
		)
	}

	return nil
}

// ExcludeChatTag excludes the group tag from the chat
//
// https://crm.pachca.com/dev/members/tags/delete/
func (c *Client) ExcludeChatTag(chatID, tagID ID) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case chatID == 0:
		return ErrInvalidChatID
	case tagID == 0:
		return ErrInvalidTagID
	}

	err := c.sendRequest(
		req.DELETE, getURL("/chats/%d/group_tags/%d", chatID, tagID),
		nil, nil, nil,
	)

	if err != nil {
		return fmt.Errorf(
			"Can't exclude group tag %d from chat %d: %w",
			tagID, chatID, err,
		)
	}

	return nil
}

// MESSAGES ///////////////////////////////////////////////////////////////////////// //

// GetMessage returns info about message
//
// https://crm.pachca.com/dev/messages/get/
func (c *Client) GetMessage(messageID ID) (*Message, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case messageID == 0:
		return nil, ErrInvalidMessageID
	}

	resp := &struct {
		Data *Message `json:"data"`
	}{}

	err := c.sendRequest(req.GET, getURL("/messages/%d", messageID), nil, nil, resp)

	if err != nil {
		return nil, fmt.Errorf("Can't fetch thread info: %w", err)
	}

	return resp.Data, nil
}

// AddMessage creates new message to user or chat
//
// https://crm.pachca.com/dev/messages/new/
func (c *Client) AddMessage(message *MessageRequest) (*Message, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case message == nil:
		return nil, ErrNilMessageRequest
	case message.EntityID == 0:
		return nil, ErrInvalidEntityID
	}

	payload := &struct {
		Message *MessageRequest `json:"message"`
	}{
		Message: message,
	}

	resp := &struct {
		Data *Message `json:"data"`
	}{}

	err := c.sendRequest(req.POST, getURL("/messages"), nil, payload, resp)

	if err != nil {
		return nil, fmt.Errorf("Can't create a new message: %w", err)
	}

	return resp.Data, nil
}

// EditMessage modifies message
//
// https://crm.pachca.com/dev/messages/update/
func (c *Client) EditMessage(messageID ID, message *MessageRequest) (*Message, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case messageID == 0:
		return nil, ErrInvalidMessageID
	case message == nil:
		return nil, ErrNilMessageRequest
	}

	payload := &struct {
		Message *MessageRequest `json:"message"`
	}{
		Message: message,
	}

	resp := &struct {
		Data *Message `json:"data"`
	}{}

	err := c.sendRequest(req.PUT, getURL("/messages/%d", messageID), nil, payload, resp)

	if err != nil {
		return nil, fmt.Errorf("Can't modify message %d: %w", messageID, err)
	}

	return resp.Data, nil
}

// DeleteMessage deletes message with given ID
//
// https://crm.pachca.com/dev/messages/delete/
func (c *Client) DeleteMessage(messageID ID) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case messageID == 0:
		return ErrInvalidMessageID
	}

	err := c.sendRequest(req.DELETE, getURL("/messages/%d", messageID), nil, nil, nil)

	if err != nil {
		return fmt.Errorf("Can't delete message %d: %w", messageID, err)
	}

	return nil
}

// PinMessage pins message to chat
//
// https://crm.pachca.com/dev/messages/pin/new/
func (c *Client) PinMessage(messageID ID) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case messageID == 0:
		return ErrInvalidMessageID
	}

	err := c.sendRequest(req.POST, getURL("/messages/%d/pin", messageID), nil, nil, nil)

	if err != nil {
		return fmt.Errorf("Can't pin message %d: %w", messageID, err)
	}

	return nil
}

// UnpinMessage unpins message from chat
//
// https://crm.pachca.com/dev/messages/pin/new/
func (c *Client) UnpinMessage(messageID ID) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case messageID == 0:
		return ErrInvalidMessageID
	}

	err := c.sendRequest(
		req.DELETE, getURL("/messages/%d/pin", messageID),
		nil, nil, nil,
	)

	if err != nil {
		return fmt.Errorf("Can't unpin message %d: %w", messageID, err)
	}

	return nil
}

// SendMessageToUser helper to send message to user with given ID
func (c *Client) SendMessageToUser(userID ID, text string) (*Message, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case userID == 0:
		return nil, ErrInvalidUserID
	case text == "":
		return nil, ErrEmptyMessage
	}

	return c.AddMessage(&MessageRequest{
		EntityType: ENTITY_TYPE_USER,
		EntityID:   userID,
		Content:    text,
	})
}

// SendMessageToChat helper to send message to chat with given ID
func (c *Client) SendMessageToChat(chatID ID, text string) (*Message, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case chatID == 0:
		return nil, ErrInvalidChatID
	case text == "":
		return nil, ErrEmptyMessage
	}

	return c.AddMessage(&MessageRequest{
		EntityType: ENTITY_TYPE_DISCUSSION,
		EntityID:   chatID,
		Content:    text,
	})
}

// SendMessageToThread helper to send message to thread with given ID
func (c *Client) SendMessageToThread(threadID ID, text string) (*Message, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case threadID == 0:
		return nil, ErrInvalidThreadID
	case text == "":
		return nil, ErrEmptyMessage
	}

	return c.AddMessage(&MessageRequest{
		EntityType: ENTITY_TYPE_THREAD,
		EntityID:   threadID,
		Content:    text,
	})
}

// ChangeMessageText helper to change message text
func (c *Client) ChangeMessageText(messageID ID, text string) (*Message, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case messageID == 0:
		return nil, ErrInvalidMessageID
	case text == "":
		return nil, ErrEmptyMessage
	}

	msg, err := c.GetMessage(messageID)

	if err != nil {
		return nil, err
	}

	msgReq := &MessageRequest{Content: text, Files: Files{}}

	if len(msg.Files) > 0 {
		msgReq.Files = msg.Files
	}

	return c.EditMessage(messageID, msgReq)
}

// THREADS ////////////////////////////////////////////////////////////////////////// //

// GetThread returns info about thread
//
// https://crm.pachca.com/dev/threads/get/
func (c *Client) GetThread(threadID ID) (*Thread, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case threadID == 0:
		return nil, ErrInvalidThreadID
	}

	resp := &struct {
		Data *Thread `json:"data"`
	}{}

	err := c.sendRequest(req.GET, getURL("/threads/%d", threadID), nil, nil, resp)

	if err != nil {
		return nil, fmt.Errorf("Can't fetch thread info: %w", err)
	}

	return resp.Data, nil
}

// NewThread creates a new tread
//
// https://crm.pachca.com/dev/threads/new/
func (c *Client) NewThread(messageID ID) (*Thread, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case messageID == 0:
		return nil, ErrInvalidMessageID
	}

	resp := &struct {
		Data *Thread `json:"data"`
	}{}

	err := c.sendRequest(
		req.POST, getURL("/messages/%d/thread", messageID),
		nil, nil, resp,
	)

	if err != nil {
		return nil, fmt.Errorf("Can't create thread for message %d: %w", messageID, err)
	}

	return resp.Data, nil
}

// AddThreadMessage helper to create thread and add new message to it
func (c *Client) AddThreadMessage(messageID ID, message *MessageRequest) (*Message, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case messageID == 0:
		return nil, ErrInvalidMessageID
	case message == nil:
		return nil, ErrNilMessageRequest
	}

	thread, err := c.NewThread(messageID)

	if err != nil {
		return nil, err
	}

	message.EntityID = thread.ID
	message.EntityType = ENTITY_TYPE_THREAD

	return c.AddMessage(message)
}

// FILES //////////////////////////////////////////////////////////////////////////// //

// UploadFile uploads new file and returns key of it
//
// https://crm.pachca.com/dev/common/files/
func (c *Client) UploadFile(file string) (*File, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case file == "":
		return nil, ErrEmptyFilePath
	}

	upload := &Upload{}
	err := c.sendRequest(req.POST, getURL("/uploads"), nil, nil, upload)

	if err != nil {
		return nil, fmt.Errorf("Can't create upload for file %q: %w", file, err)
	}

	fmt.Printf("%#v\n", upload)

	info, err := createMultipartData(file, upload, c.MaxFileSize)

	if err != nil {
		return nil, err
	}

	resp, err := c.engine.Post(
		req.Request{
			URL:         upload.DirectURL,
			ContentType: info.ContentType,
			Body:        info.Buffer,
		},
	)

	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("Can't upload file %q data: %w", file, err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf(
			"Can't upload file %q data: %w",
			file, extractS3Error(resp.String()),
		)
	}

	return &File{
		Key:  info.Key,
		Name: info.Name,
		Size: info.Size,
		Type: guessFileType(info.Name),
	}, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Get returns custom property with given name
func (p Properties) Get(name string) *Property {
	for _, pp := range p {
		if pp.Name == name {
			return pp
		}
	}

	return nil
}

// GetAny returns first found property with one of given names
func (p Properties) GetAny(name ...string) *Property {
	for _, pp := range p {
		if slices.Contains(name, pp.Name) {
			return pp
		}
	}

	return nil
}

// Names returns slice with properties names
func (p Properties) Names() []string {
	var result []string

	for _, pp := range p {
		result = append(result, pp.Name)
	}

	return result
}

// IsText returns true if property has text type
func (p *Property) IsText() bool {
	return p != nil && p.Type == PROP_TYPE_TEXT
}

// IsLink returns true if property has URL type
func (p *Property) IsLink() bool {
	return p != nil && p.Type == PROP_TYPE_LINK
}

// IsDate returns true if property has date type
func (p *Property) IsDate() bool {
	return p != nil && p.Type == PROP_TYPE_DATE
}

// IsNumber returns true if property has number type
func (p *Property) IsNumber() bool {
	return p != nil && p.Type == PROP_TYPE_NUMBER
}

// String returns property value
func (p *Property) String() string {
	if p == nil {
		return ""
	}

	return p.Value
}

// ToDate tries to convert property value to date
func (p *Property) ToDate() (time.Time, error) {
	switch {
	case p == nil:
		return time.Time{}, ErrNilProperty
	case p.Value == "":
		return time.Time{}, nil
	case p.Type != PROP_TYPE_DATE:
		return time.Time{}, fmt.Errorf("Invalid property type for date (%s)", p.Type)
	}

	return parseDate(p.Value)
}

// Date returns property value as date
func (p *Property) Date() time.Time {
	d, _ := p.ToDate()
	return d
}

// ToInt tries to convert property value to int
func (p *Property) ToInt() (int, error) {
	switch {
	case p == nil:
		return 0, ErrNilProperty
	case p.Value == "":
		return 0, nil
	case p.Type != PROP_TYPE_NUMBER:
		return 0, fmt.Errorf("Invalid property type for date (%s)", p.Type)
	}

	return strconv.Atoi(p.Value)
}

// Int returns property value as int
func (p *Property) Int() int {
	i, _ := p.ToInt()
	return i
}

// FullName returns user full name
func (u *User) FullName() string {
	if u == nil {
		return ""
	}

	switch {
	case u == nil:
		return ""
	case u.FirstName == "" && u.LastName != "":
		return u.LastName
	case u.FirstName != "" && u.LastName == "":
		return u.FirstName
	}

	return u.FirstName + " " + u.LastName
}

// Active returns slice with active users
func (u Users) Active() Users {
	var result Users

	for _, uu := range u {
		if !uu.IsSuspended {
			result = append(result, uu)
		}
	}

	return result
}

// Suspended returns slice with inactive users
func (u Users) Suspended() Users {
	var result Users

	for _, uu := range u {
		if uu.IsSuspended {
			result = append(result, uu)
		}
	}

	return result
}

// Invited returns all invited users
func (u Users) Invited() Users {
	var result Users

	for _, uu := range u {
		if uu.InviteStatus == INVITE_SENT {
			result = append(result, uu)
		}
	}

	return result
}

// Bots returns slice with bots
func (u Users) Bots() Users {
	var result Users

	for _, uu := range u {
		if uu.IsBot {
			result = append(result, uu)
		}
	}

	return result
}

// Admins returns slice with admins
func (u Users) Admins() Users {
	var result Users

	for _, uu := range u {
		if uu.Role == ROLE_ADMIN {
			result = append(result, uu)
		}
	}

	return result
}

// Regular returns slice with regular users
func (u Users) Regular() Users {
	var result Users

	for _, uu := range u {
		if uu.Role == ROLE_USER {
			result = append(result, uu)
		}
	}

	return result
}

// Guests returns slice with guests
func (u Users) Guests() Users {
	var result Users

	for _, uu := range u {
		if uu.Role == ROLE_MULTI_GUEST {
			result = append(result, uu)
		}
	}

	return result
}

// Get returns chat with given name
func (c Chats) Get(name string) *Chat {
	for _, cc := range c {
		if cc.Name == name {
			return cc
		}
	}

	return nil
}

// Public returns slice with public chats
func (c Chats) Public() Chats {
	var result Chats

	for _, cc := range c {
		if cc.IsPublic {
			result = append(result, cc)
		}
	}

	return result
}

// Channels returns slice with channels
func (c Chats) Channels() Chats {
	var result Chats

	for _, cc := range c {
		if cc.IsChannel {
			result = append(result, cc)
		}
	}

	return result
}

// URL returns chat URL
func (c *Chat) URL() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("%s/chats/%d", APP_URL, c.ID)
}

// URL returns URL of user profile
func (u *User) URL() string {
	if u == nil {
		return ""
	}

	return fmt.Sprintf("%s/chats?user_id=%d", APP_URL, u.ID)
}

// URL returns message URL
func (m *Message) URL() string {
	if m == nil {
		return ""
	}

	return fmt.Sprintf(
		"%s/chats/%d?message=%d",
		APP_URL, m.ChatID, m.ID,
	)
}

// URL returns thread URL
func (t *Thread) URL() string {
	if t == nil {
		return ""
	}

	return fmt.Sprintf(
		"%s/chats?thread_message_id=%d&sidebar_message=%d",
		APP_URL, t.MessageID, t.ID,
	)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// IsReaction returns true if webhook contains data for message event
func (w *Webhook) IsMessage() bool {
	return w != nil && w.Type == WEBHOOK_TYPE_MESSAGE
}

// IsReaction returns true if webhook contains data for reaction event
func (w *Webhook) IsReaction() bool {
	return w != nil && w.Type == WEBHOOK_TYPE_REACTION
}

// IsReaction returns true if webhook contains data for button event
func (w *Webhook) IsButton() bool {
	return w != nil && w.Type == WEBHOOK_TYPE_BUTTON
}

// IsNew returns true if there is a webhook event for new message
func (w *Webhook) IsNew() bool {
	return w != nil && w.Event == WEBHOOK_EVENT_NEW
}

// IsUpdate returns true if there is a webhook event for updated message
func (w *Webhook) IsUpdate() bool {
	return w != nil && w.Event == WEBHOOK_EVENT_UPDATE
}

// IsDelete returns true if there is a webhook event for deleted message
func (w *Webhook) IsDelete() bool {
	return w != nil && w.Event == WEBHOOK_EVENT_DELETE
}

// Command returns slash command name from the beginning of the message
func (w *Webhook) Command() string {
	if w == nil || w.Content == "" {
		return ""
	}

	return strings.TrimLeft(strutil.ReadField(w.Content, 0, false, ' '), "/")
}

// ////////////////////////////////////////////////////////////////////////////////// //

// ToQuery converts filter struct to request query
func (f ChatFilter) ToQuery() req.Query {
	query := req.Query{}

	if f.Public {
		query["availability"] = "public"
	}

	if !f.LastMessageBefore.IsZero() {
		query["last_message_at_before"] = formatDate(f.LastMessageBefore)
	}

	if !f.LastMessageAfter.IsZero() {
		query["last_message_at_after"] = formatDate(f.LastMessageAfter)
	}

	return query
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getBatchSize returns batch size for paginated responses
func (c *Client) getBatchSize() int {
	return mathutil.Between(c.BatchSize, 5, 50)
}

// sendRequest sends request to Pachca API
func (c *Client) sendRequest(method, url string, query req.Query, payload any, response any) error {
	r := req.Request{
		Method:     method,
		URL:        url,
		Query:      query,
		Accept:     req.CONTENT_TYPE_JSON,
		BearerAuth: c.token,
	}

	if payload != nil {
		r.ContentType = req.CONTENT_TYPE_JSON
		r.Body = payload
	}

	resp, err := c.engine.Do(r)

	if err != nil {
		return fmt.Errorf("Can't send request to API: %w", err)
	}

	defer resp.Discard()

	if resp.StatusCode >= 400 {
		errResp := &struct {
			Errors []APIError `json:"errors"`
		}{}

		err = resp.JSON(errResp)

		if err != nil || len(errResp.Errors) == 0 {
			return fmt.Errorf("API returned non-ok status code %d", resp.StatusCode)
		}

		errResp.Errors[0].StatusCode = resp.StatusCode

		return errResp.Errors[0]
	}

	if response != nil {
		err = resp.JSON(response)

		if err != nil {
			return fmt.Errorf("Can't decode API response: %w", err)
		}
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getURL returns full URL of API endpoint
func getURL(endpoint string, args ...any) string {
	if len(args) == 0 {
		return API_URL + endpoint
	}

	return API_URL + fmt.Sprintf(endpoint, args...)
}

// createMultipartData creates multipart data with file data
func createMultipartData(file string, upload *Upload, maxFileSize int64) (*uploadInfo, error) {
	fd, err := os.Open(file)

	if err != nil {
		return nil, fmt.Errorf("Can't open file %q to upload: %w", file, err)
	}

	stat, err := fd.Stat()

	if err != nil {
		return nil, fmt.Errorf("Can't get file %q info: %w", file, err)
	}

	if stat.Size() >= maxFileSize {
		return nil, fmt.Errorf("File size exceeds the limit (%d): %w", maxFileSize, err)
	}

	fileName := path.Base(fd.Name())

	info := &uploadInfo{
		Key:  strings.ReplaceAll(upload.Key, "${filename}", fileName),
		Name: fileName,
		Size: uint(stat.Size()),
	}

	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)

	var errs errors.Bundle

	errs.Add(
		mw.WriteField("Content-Disposition", upload.ContentDisposition),
		mw.WriteField("acl", upload.ACL),
		mw.WriteField("policy", upload.Policy),
		mw.WriteField("x-amz-credential", upload.Credential),
		mw.WriteField("x-amz-algorithm", upload.Algorithm),
		mw.WriteField("x-amz-date", upload.Date),
		mw.WriteField("x-amz-signature", upload.Signature),
		mw.WriteField("key", upload.Key),
	)

	if !errs.IsEmpty() {
		return nil, fmt.Errorf("Can't create multipart upload: %w", errs.First())
	}

	fw, err := mw.CreateFormFile("file", fd.Name())

	if err != nil {
		return nil, fmt.Errorf("Can't write file %q part: %w", file, err)
	}

	_, err = io.Copy(fw, fd)

	if err != nil {
		return nil, fmt.Errorf("Can't write file %q part: %w", file, err)
	}

	errs.Reset()

	errs.Add(
		mw.Close(),
		fd.Close(),
	)

	if !errs.IsEmpty() {
		return nil, errs.First()
	}

	info.Buffer = buf
	info.ContentType = mw.FormDataContentType()

	return info, nil
}

// extractS3Error extracts error text from S3 error message
func extractS3Error(errorMessage string) error {
	found := s3ErrorExtractRegex.FindStringSubmatch(errorMessage)

	if len(found) == 2 {
		return errors.New(found[1])
	}

	return errors.New("Unknown error")
}

// formatDate converts date from time.Time to string (ISO-8601)
func formatDate(d time.Time) string {
	return d.Format("2006-01-02T15:04:05.999Z")
}

// parseDate converts date from string (ISO-8601) to time.Time
func parseDate(d string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.999Z", d)
}

// guessFileType tries to guess file type by it extension
func guessFileType(name string) FileType {
	switch strings.ToLower(path.Ext(name)) {
	case ".jpg", ".jpeg", ".png", ".bmp", ".gif":
		return FILE_TYPE_IMAGE
	}

	return FILE_TYPE_FILE
}
