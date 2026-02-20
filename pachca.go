package pachca

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"encoding/json"
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
	"github.com/essentialkaos/ek/v13/sliceutil"
	"github.com/essentialkaos/ek/v13/strutil"

	"github.com/essentialkaos/pachca/block"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// API_URL is URL of Pachca API
const API_URL = "https://api.pachca.com/api/shared/v1"

// APP_URL is application URL used to generate links
const APP_URL = "https://app.pachca.com"

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	SORT_FIELD_ID           = "id"
	SORT_FIELD_LAST_MESSAGE = "last_message_at"
)

const (
	SORT_ORDER_ASC  = "asc"
	SORT_ORDER_DESC = "desc"
)

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
	ROLE_REGULAR     UserRole = "user"
	ROLE_MULTI_GUEST UserRole = "multi_guest"
	ROLE_GUEST       UserRole = "guest"
)

const (
	CHAT_ROLE_ADMIN  ChatRole = "admin"
	CHAT_ROLE_OWNER  ChatRole = "owner"
	CHAT_ROLE_EDITOR ChatRole = "editor"
	CHAT_ROLE_MEMBER ChatRole = "member"

	CHAT_ROLE_ANY = "all"
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
	VIEW_MODAL ViewType = "modal"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// MAX_PAGES is the maximum number of pages using for listing items
const MAX_PAGES = 100_000

// MAX_PER_PAGE is the maximum number of entities per page
const MAX_PER_PAGE = 50

// ////////////////////////////////////////////////////////////////////////////////// //

// Date is JSON date
type Date struct {
	time.Time
}

// ////////////////////////////////////////////////////////////////////////////////// //

// EntityType is type of entity type
type EntityType string

// PropertyType is type for property type
type PropertyType string

// UserRole is type of user role
type UserRole string

// ChatRole is type of user in chat
type ChatRole string

// InviteStatus is type of invite status
type InviteStatus string

// FileType is type for file type
type FileType string

// ViewType is type for view
type ViewType string

// ////////////////////////////////////////////////////////////////////////////////// //

// Chats is slice of chats
type Chats []*Chat

// Chat contains info about channel
type Chat struct {
	Members       []uint `json:"member_ids"`
	GroupTags     []uint `json:"group_tag_ids"`
	ID            uint   `json:"id"`
	OwnerID       uint   `json:"owner_id"`
	Name          string `json:"name"`
	MeetRoomURL   string `json:"meet_room_url"`
	CreatedAt     Date   `json:"created_at"`
	LastMessageAt Date   `json:"last_message_at"`
	IsPublic      bool   `json:"public"`
	IsChannel     bool   `json:"channel"`
	IsPersonal    bool   `json:"personal"`
}

// Users is a slice of users
type Users []*User

// User contains info about user
type User struct {
	ID             uint         `json:"id"`
	CreatedAt      Date         `json:"created_at"`
	LastActivityAt Date         `json:"last_activity_at"`
	ImageURL       string       `json:"image_url"`
	Email          string       `json:"email"`
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	Nickname       string       `json:"nickname"`
	Role           UserRole     `json:"role"`
	PhoneNumber    string       `json:"phone_number"`
	TimeZone       string       `json:"time_zone"`
	Title          string       `json:"title"`
	InviteStatus   InviteStatus `json:"invite_status"`
	Department     string       `json:"department"`
	Properties     Properties   `json:"custom_properties"`
	Tags           []string     `json:"list_tags"`
	Status         *Status      `json:"user_status"`
	IsBot          bool         `json:"bot"`
	IsSuspended    bool         `json:"suspended"`
	IsSSO          bool         `json:"sso"`
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
	ID    uint         `json:"id"`
	Type  PropertyType `json:"data_type"`
	Name  string       `json:"name"`
	Value string       `json:"value"`
}

// Tag contains info about tag
type Tag struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	UsersCount int    `json:"users_count"`
}

// Tags is a slice of tags
type Tags []*Tag

// Reaction contains reaction info
type Reaction struct {
	UserID    uint   `json:"user_id"`
	CreatedAt Date   `json:"created_at"`
	Emoji     string `json:"code"`
	Name      string `json:"name"`
}

// Reactions is a slice of reactions
type Reactions []*Reaction

// Thread contains info about thread
type Thread struct {
	ID            uint `json:"id"`
	ChatID        uint `json:"chat_id"`
	MessageID     uint `json:"message_id"`
	MessageChatID uint `json:"message_chat_id"`
	UpdatedAt     Date `json:"updated_at"`
}

// Message contains info about message
type Message struct {
	ID              uint        `json:"id"`
	EntityID        uint        `json:"entity_id"`
	ChatID          uint        `json:"chat_id"`
	ParentMessageID uint        `json:"parent_message_id"`
	UsedID          uint        `json:"user_id"`
	EntityType      EntityType  `json:"entity_type"`
	Content         string      `json:"content"`
	CreatedAt       Date        `json:"created_at"`
	Thread          *Thread     `json:"thread"`
	Files           Files       `json:"files"`
	Buttons         Buttons     `json:"buttons"`
	Forwarding      *Forwarding `json:"forwarding"`
}

// Messages is a slice of messages
type Messages []*Message

// Forwarding contains info about message forwarding
type Forwarding struct {
	OriginalMessageID          uint `json:"original_message_id"`
	OriginalChatID             uint `json:"original_chat_id"`
	AuthorID                   uint `json:"author_id"`
	OriginalThreadID           uint `json:"original_thread_id"`
	OriginalThreadMessageID    uint `json:"original_thread_message_id"`
	OriginalThreadParentChatID uint `json:"original_thread_parent_chat_id"`
	OriginalCreatedAt          Date `json:"original_created_at"`
}

// File contains info about message attachment
type File struct {
	ID     uint     `json:"id,omitempty"`
	Key    string   `json:"key"`
	Name   string   `json:"name"`
	Type   FileType `json:"file_type,omitempty"`
	URL    string   `json:"url,omitempty"`
	Size   int64    `json:"size"`
	Width  int      `json:"width,omitzero"`
	Height int      `json:"height,omitzero"`
}

// Files is a slice of attachments
type Files []*File

// Button contains info about message button
type Button struct {
	Text string `json:"text"`
	URL  string `json:"url,omitempty"`
	Data string `json:"data,omitempty"`
}

// Buttons is a slice of buttons
type Buttons []ButtonLine

// ButtonCollection
type ButtonLine []*Button

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

// View contains form view
type View struct {
	Title      string        `json:"title"`
	CloseText  string        `json:"close_text,omitempty"`
	SubmitText string        `json:"submit_text,omitempty"`
	Blocks     []block.Block `json:"blocks"`
}

// WebhookEvent contains webhook event data
type WebhookEvent struct {
	ID        string          `json:"id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
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

// Metadata is listing metadata
type Metadata struct {
	Paginate *Paginate `json:"paginate"`
}

// Paginate contains cursor to the next page
type Paginate struct {
	NextPage string `json:"next_page"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// uploadInfo contains info about uploaded file
type uploadInfo struct {
	Key         string    // Uploading key
	Name        string    // File name
	Size        int64     // File size
	ContentType string    // Content type
	Reader      io.Reader // Data reader
}

// ////////////////////////////////////////////////////////////////////////////////// //

// ChatFilter is configuration for filtering chats
type ChatFilter struct {
	Sort              map[string]string
	LastMessageAfter  time.Time
	LastMessageBefore time.Time
	Public            bool
}

// UserRequest is a struct with information needed to create or modify a user
type UserRequest struct {
	Email           string           `json:"email,omitempty"`
	FirstName       string           `json:"first_name,omitempty"`
	LastName        string           `json:"last_name,omitempty"`
	Nickname        string           `json:"nickname,omitempty"`
	Role            UserRole         `json:"role,omitempty"`
	PhoneNumber     string           `json:"phone_number,omitempty"`
	Title           string           `json:"title,omitempty"`
	Department      string           `json:"department,omitempty"`
	Properties      PropertyRequests `json:"custom_properties,omitempty"`
	Tags            []string         `json:"list_tags,omitempty"`
	IsSuspended     bool             `json:"suspended,omitempty"`
	SkipEmailNotify bool             `json:"skip_email_notify,omitempty"`
}

// PropertyRequest is a struct with property info
type PropertyRequest struct {
	ID    uint   `json:"id"`
	Value string `json:"value"`
}

// PropertyRequests is a slice with properties requests
type PropertyRequests []*PropertyRequest

// ChatRequest is a struct with information needed to create or modify a chat
type ChatRequest struct {
	Name       string `json:"name,omitempty"`
	Members    []uint `json:"member_ids,omitempty"`
	Groups     []uint `json:"group_tag_ids,omitempty"`
	IsChannel  bool   `json:"channel,omitempty"`
	IsPublic   bool   `json:"public,omitempty"`
	IsPersonal bool   `json:"personal,omitempty"`
}

// MessageRequest is a struct with information needed to create or modify a message
type MessageRequest struct {
	EntityType         EntityType `json:"entity_type,omitempty"`
	EntityID           uint       `json:"entity_id,omitempty"`
	ParentMessageID    uint       `json:"parent_message_id,omitempty"`
	Content            string     `json:"content"`
	DisplayAvatarURL   string     `json:"display_avatar_url,omitempty"`
	DisplayName        string     `json:"display_name,omitempty"`
	Files              Files      `json:"files,omitzero"`
	Buttons            Buttons    `json:"buttons,omitzero"`
	SkipInviteMentions bool       `json:"skip_invite_mentions,omitempty"`
}

// ReactionRequest is a payload for message reaction
type ReactionRequest struct {
	Code string `json:"code"`
	Name string `json:"name,omitempty"`
}

// ViewRequest is a payload for open a view
type ViewRequest struct {
	Type       ViewType `json:"type"`
	TriggerID  string   `json:"trigger_id"`
	Metadata   string   `json:"private_metadata,omitempty"`
	CallbackID string   `json:"callback_id,omitempty"`
	View       *View    `json:"view"`
}

// LinkPreview contains link preview data
type LinkPreview struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url,omitempty"`
	Image       *File  `json:"image,omitempty"`
}

// LinkPreviews is map (url → preview data) with link previews
type LinkPreviews map[string]*LinkPreview

// ////////////////////////////////////////////////////////////////////////////////// //

// ViewErrors is a map with view errors
type ViewErrors map[string]string

// S3Error represents S3 error
type S3Error struct {
	Message string
	Full    string
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
	return fmt.Sprintf(
		"(%s) %s [%s:%s]",
		e.Code, e.Message, e.Key, strutil.Q(e.Value, "-"),
	)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// tokenValidationRegex is regex pattern for token validation
var tokenValidationRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_]{43}$`)

// s3ErrorExtractRegex is regex pattern for extracting text from S3 error message
var s3ErrorExtractRegex = regexp.MustCompile(`\<Message\>(.*)\<\/Message\>`)

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrNilClient          = errors.New("client is nil")
	ErrNilUserRequest     = errors.New("user request is nil")
	ErrNilChatRequest     = errors.New("chat request is nil")
	ErrNilMessageRequest  = errors.New("message request is nil")
	ErrNilPropertyRequest = errors.New("property request is nil")
	ErrNilViewRequest     = errors.New("view request is nil")
	ErrNilView            = errors.New("view data is nil")
	ErrEmptyToken         = errors.New("token is empty")
	ErrEmptyTag           = errors.New("group tag is empty")
	ErrEmptyMessage       = errors.New("message text is empty")
	ErrEmptyUserEmail     = errors.New("user email is required for creating user account")
	ErrEmptyChatName      = errors.New("name is required for creating new chat")
	ErrEmptyUsersIDS      = errors.New("users IDs are empty")
	ErrEmptyTagsIDS       = errors.New("tags IDs are empty")
	ErrEmptyFilePath      = errors.New("path to file is empty")
	ErrInvalidToken       = errors.New("token has wrong format")
	ErrInvalidMessageID   = errors.New("message ID must be greater than 0")
	ErrInvalidChatID      = errors.New("chat ID must be greater than 0")
	ErrInvalidUserID      = errors.New("user ID must be greater than 0")
	ErrInvalidThreadID    = errors.New("thread ID must be greater than 0")
	ErrInvalidTagID       = errors.New("group tag ID must be greater than 0")
	ErrInvalidEntityID    = errors.New("entity ID must be greater than 0")
	ErrInvalidBotID       = errors.New("bot ID must be greater than 0")
	ErrInvalidEventID     = errors.New("invalid event ID")
	ErrBlankReaction      = errors.New("non-blank emoji is required")
	ErrEmptyPreviews      = errors.New("previews map has no data")
	ErrInvalidPageNum     = errors.New("page number must be greater than 0")
	ErrInvalidMessageNum  = errors.New("number of messages must be greater than 0")
	ErrInvalidPerPageNum  = errors.New("per page number must be between 1 and 50")
	ErrViewHasNoBlocks    = errors.New("view has no blocks")
	ErrEmptyTriggerID     = errors.New("view has empty trigger ID")
	ErrInvalidMaxPages    = errors.New("minimum number of result pages must be greater than 0")
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

	e := &req.Engine{}
	e.SetUserAgent("EK|Pachca.go", "1")

	return &Client{
		BatchSize:   MAX_PER_PAGE,
		MaxFileSize: 10 * 1024 * 1024, // 10 MB

		token:  token,
		engine: e,
	}, nil
}

// NewPropertyRequest creates new custom property
func NewPropertyRequest(id uint, value any) *PropertyRequest {
	var v string

	switch t := value.(type) {
	case time.Time:
		v = formatDate(t.UTC())

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		v = fmt.Sprintf("%d", value)

	case float32:
		v = fmt.Sprintf("%d", int64(t))

	case float64:
		v = fmt.Sprintf("%d", int64(t))

	default:
		v = fmt.Sprintf("%v", value)
	}

	return &PropertyRequest{ID: id, Value: v}
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

	c.engine.SetUserAgent(app, ver, "EK|Pachca.go/0")
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
// https://dev.pachca.com/common/list
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
		return nil, fmt.Errorf("can't fetch custom properties: %w", err)
	}

	return resp.Data, nil
}

// REACTIONS //////////////////////////////////////////////////////////////////////// //

// GetReactions returns slice with reactions added to given message
//
// https://dev.pachca.com/reactions/list-reactions
func (c *Client) GetReactions(messageID uint) (Reactions, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case messageID == 0:
		return nil, ErrInvalidMessageID
	}

	var result Reactions

	query := req.Query{"per": c.getBatchSize()}

	for i := 1; i < MAX_PAGES; i++ {
		query["page"] = i

		resp := &struct {
			Data Reactions `json:"data"`
		}{}

		err := c.sendRequest(
			req.GET, getURL("/messages/%d/reactions", messageID),
			query, nil, resp,
		)

		if err != nil {
			return nil, fmt.Errorf("can't fetch reactions for message %d: %w", messageID, err)
		}

		result = append(result, resp.Data...)

		if len(resp.Data) != c.getBatchSize() {
			break
		}
	}

	return result, nil
}

// AddReaction adds given emoji reaction to the message. To add custom reaction
// add it name after amoji using ":" as separator. For example "😲:omg".
//
// https://dev.pachca.com/reactions/add-reactions
func (c *Client) AddReaction(messageID uint, reaction string) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case messageID == 0:
		return ErrInvalidMessageID
	case reaction == "":
		return ErrBlankReaction
	}

	emoji, name, _ := strings.Cut(reaction, ":")
	payload := &ReactionRequest{Code: emoji, Name: name}

	err := c.sendRequest(
		req.POST, getURL("/messages/%d/reactions", messageID),
		nil, payload, nil,
	)

	if err != nil {
		return fmt.Errorf("can't add reaction %q to message %d: %w", reaction, messageID, err)
	}

	return nil
}

// DeleteReaction removes given emoji reaction from the message
//
// https://dev.pachca.com/reactions/remove-reactions
func (c *Client) DeleteReaction(messageID uint, reaction string) error {
	switch {
	case c == nil:
		return ErrNilClient
	case messageID == 0:
		return ErrInvalidMessageID
	case reaction == "":
		return ErrBlankReaction
	}

	emoji, name, _ := strings.Cut(reaction, ":")
	payload := &ReactionRequest{Code: emoji, Name: name}

	err := c.sendRequest(
		req.DELETE, getURL("/messages/%d/reactions", messageID),
		nil, payload, nil,
	)

	if err != nil {
		return fmt.Errorf("can't remove reaction %q from message %d: %w", reaction, messageID, err)
	}

	return nil
}

// USERS //////////////////////////////////////////////////////////////////////////// //

// CurrentUser returns info about current user
//
// https://dev.pachca.com/profile/profile
func (c *Client) CurrentUser() (*User, error) {
	if c == nil || c.engine == nil {
		return nil, ErrNilClient
	}

	resp := &struct {
		Data *User `json:"data"`
	}{}

	err := c.sendRequest(req.GET, getURL("/profile"), nil, nil, resp)

	if err != nil {
		return nil, fmt.Errorf("can't fetch current user info: %w", err)
	}

	return resp.Data, nil
}

// GetUser returns info about specific user
//
// https://dev.pachca.com/users/get
func (c *Client) GetUser(userID uint) (*User, error) {
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
		return nil, fmt.Errorf("can't fetch user info: %w", err)
	}

	return resp.Data, nil
}

// GetUsers returns info about all users
//
// https://dev.pachca.com/users/list
func (c *Client) GetUsers(searchQuery ...string) (Users, error) {
	if c == nil || c.engine == nil {
		return nil, ErrNilClient
	}

	var result Users

	query := req.Query{"per": c.getBatchSize()}

	if len(searchQuery) != 0 {
		query["query"] = searchQuery[0]
	}

	for i := 1; i < MAX_PAGES; i++ {
		query["page"] = i

		resp := &struct {
			Data Users `json:"data"`
		}{}

		err := c.sendRequest(req.GET, getURL("/users"), query, nil, resp)

		if err != nil {
			return nil, fmt.Errorf("can't fetch users: %w", err)
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
// https://dev.pachca.com/users/create
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
		return nil, fmt.Errorf("can't create a new user: %w", err)
	}

	return resp.Data, nil
}

// EditUser modifies an existing user
//
// https://dev.pachca.com/users/update
func (c *Client) EditUser(userID uint, user *UserRequest) (*User, error) {
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
		return nil, fmt.Errorf("can't edit user %d: %w", userID, err)
	}

	return resp.Data, nil
}

// DeleteUser deletes an existing user
//
// https://dev.pachca.com/users/delete
func (c *Client) DeleteUser(userID uint) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case userID == 0:
		return ErrInvalidUserID
	}

	err := c.sendRequest(req.DELETE, getURL("/users/%d", userID), nil, nil, nil)

	if err != nil {
		return fmt.Errorf("can't delete user %d: %w", userID, err)
	}

	return nil
}

// GROUP TAGS /////////////////////////////////////////////////////////////////////// //

// GetTags returns all group tags
//
// https://dev.pachca.com/group-tags/list
func (c *Client) GetTags(names ...string) (Tags, error) {
	if c == nil || c.engine == nil {
		return nil, ErrNilClient
	}

	var result Tags

	query := req.Query{"per": c.getBatchSize()}
	query.SetIf(len(names) > 0, "names[]", names)

	for i := 1; i < MAX_PAGES; i++ {
		query["page"] = i

		resp := &struct {
			Data Tags `json:"data"`
		}{}

		err := c.sendRequest(req.GET, getURL("/group_tags"), query, nil, resp)

		if err != nil {
			return nil, fmt.Errorf("can't fetch group tags: %w", err)
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
// https://dev.pachca.com/group-tags/get
func (c *Client) GetTag(groupTagID uint) (*Tag, error) {
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
		return nil, fmt.Errorf("can't fetch group tag: %w", err)
	}

	return resp.Data, nil
}

// GetTagUsers returns slice with users with given tag
//
// https://dev.pachca.com/group-tags/list-users
func (c *Client) GetTagUsers(groupTagID uint) (Users, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case groupTagID == 0:
		return nil, ErrInvalidTagID
	}

	var result Users

	query := req.Query{"per": c.getBatchSize()}

	for i := 1; i < MAX_PAGES; i++ {
		query["page"] = i

		resp := &struct {
			Data Users `json:"data"`
		}{}

		err := c.sendRequest(
			req.GET, getURL("/group_tags/%d/users", groupTagID),
			query, nil, resp,
		)

		if err != nil {
			return nil, fmt.Errorf("can't fetch group tag users: %w", err)
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
// https://dev.pachca.com/group-tags/create
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
		return nil, fmt.Errorf("can't create new group tag %q: %w", groupTagName, err)
	}

	return resp.Data, nil
}

// EditTag changes name of given group tag
//
// http://dev.pachca.com/group-tags/update
func (c *Client) EditTag(groupTagID uint, groupTagName string) (*Tag, error) {
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
		return nil, fmt.Errorf("can't edit group tag %d: %w", groupTagID, err)
	}

	return resp.Data, nil
}

// DeleteTag removes group tag
//
// https://dev.pachca.com/group-tags/delete
func (c *Client) DeleteTag(groupTagID uint) error {
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
		return fmt.Errorf("can't delete group tag %d: %w", groupTagID, err)
	}

	return nil
}

// CHATS //////////////////////////////////////////////////////////////////////////// //

// GetChats returns all chats and conversations
//
// https://dev.pachca.com/chats/list
func (c *Client) GetChats(filter ...ChatFilter) (Chats, error) {
	if c == nil || c.engine == nil {
		return nil, ErrNilClient
	}

	var result Chats
	var query req.Query

	if len(filter) == 0 {
		query = req.Query{"limit": c.getBatchSize()}
	} else {
		query = filter[0].ToQuery()
		query["limit"] = c.getBatchSize()
	}

	for i := 0; i < MAX_PAGES; i++ {
		resp := &struct {
			Data Chats     `json:"data"`
			Meta *Metadata `json:"meta"`
		}{}

		err := c.sendRequest(req.GET, getURL("/chats"), query, nil, resp)

		if err != nil {
			return nil, fmt.Errorf("can't fetch chats: %w", err)
		}

		result = append(result, resp.Data...)

		if len(resp.Data) != c.getBatchSize() {
			break
		}

		query.SetIf(
			resp.Meta != nil && resp.Meta.Paginate != nil,
			"cursor", resp.Meta.Paginate.NextPage,
		)
	}

	return result, nil
}

// GetChat returns info about specific channel
//
// https://dev.pachca.com/chats/get
func (c *Client) GetChat(chatID uint) (*Chat, error) {
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
		return nil, fmt.Errorf("can't fetch chat info: %w", err)
	}

	return resp.Data, nil
}

// AddChat creates new chat
//
// https://dev.pachca.com/chats/create
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
		return nil, fmt.Errorf("can't create a new chat %q: %w", chat.Name, err)
	}

	return resp.Data, nil
}

// EditChat modifies chat
//
// https://dev.pachca.com/chats/update
func (c *Client) EditChat(chatID uint, chat *ChatRequest) (*Chat, error) {
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
		return nil, fmt.Errorf("can't modify chat %d: %w", chatID, err)
	}

	return resp.Data, nil
}

// GetChatUsers returns all users of given chat
//
// https://dev.pachca.com/members/list-members
func (c *Client) GetChatUsers(chatID uint, memberRole ChatRole) (Users, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case chatID == 0:
		return nil, ErrInvalidChatID
	}

	switch memberRole {
	case "":
		memberRole = CHAT_ROLE_ANY
	case CHAT_ROLE_ANY, CHAT_ROLE_ADMIN, CHAT_ROLE_OWNER,
		CHAT_ROLE_EDITOR, CHAT_ROLE_MEMBER:
		// okay
	default:
		return nil, fmt.Errorf("unknown chat users role %q", memberRole)
	}

	query := req.Query{
		"role":  memberRole,
		"limit": c.getBatchSize(),
	}

	var users Users

	for i := 0; i < MAX_PAGES; i++ {
		resp := &struct {
			Data Users     `json:"data"`
			Meta *Metadata `json:"meta"`
		}{}

		err := c.sendRequest(
			req.GET, getURL("/chats/%d/members", chatID),
			query, nil, resp,
		)

		if err != nil {
			return nil, fmt.Errorf("can't fetch chat users info: %w", err)
		}

		users = append(users, resp.Data...)

		if len(resp.Data) != c.getBatchSize() {
			break
		}

		query.SetIf(
			resp.Meta != nil && resp.Meta.Paginate != nil,
			"cursor", resp.Meta.Paginate.NextPage,
		)
	}

	return users, nil
}

// AddChatUsers adds users with given IDs to the chat, channel or thread
//
// https://dev.pachca.com/members/add-members
func (c *Client) AddChatUsers(chatID uint, membersIDs []uint, silent bool) error {
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
		IDs []uint `json:"member_ids"`
	}{
		IDs: membersIDs,
	}

	err := c.sendRequest(
		req.POST, getURL("/chats/%d/members", chatID),
		query, payload, nil,
	)

	if err != nil {
		return fmt.Errorf("can't add users to chat %d: %w", chatID, err)
	}

	return nil
}

// AddChatTags adds group tags to the chat
//
// https://dev.pachca.com/members/add-group-tags
func (c *Client) AddChatTags(chatID uint, tagIDs []uint) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case chatID == 0:
		return ErrInvalidChatID
	case len(tagIDs) == 0:
		return ErrEmptyTagsIDS
	}

	payload := &struct {
		IDs []uint `json:"group_tag_ids"`
	}{
		IDs: tagIDs,
	}

	err := c.sendRequest(
		req.PUT, getURL("/chats/%d/group_tags", chatID),
		nil, payload, nil,
	)

	if err != nil {
		return fmt.Errorf("can't add group tags to chat %d: %w", chatID, err)
	}

	return nil
}

// SetChatUserRole sets user role in given chat
//
// https://dev.pachca.com/members/update-members
func (c *Client) SetChatUserRole(chatID, userID uint, role ChatRole) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case chatID == 0:
		return ErrInvalidChatID
	case userID == 0:
		return ErrInvalidUserID
	}

	switch role {
	case CHAT_ROLE_ADMIN, CHAT_ROLE_EDITOR, CHAT_ROLE_MEMBER:
		// okay
	default:
		return fmt.Errorf(
			"invalid chat role %q (must be %s, %s or %s)",
			role, CHAT_ROLE_ADMIN, CHAT_ROLE_EDITOR, CHAT_ROLE_MEMBER,
		)
	}

	err := c.sendRequest(
		req.PUT, getURL("/chats/%d/members/%d", chatID, userID),
		req.Query{"role": role}, nil, nil,
	)

	if err != nil {
		return fmt.Errorf(
			"can't set role to %q for user with ID %d in chat %d: %w",
			role, userID, chatID, err,
		)
	}

	return nil
}

// ExcludeChatUser excludes the user from the chat
//
// https://dev.pachca.com/members/remove-member
func (c *Client) ExcludeChatUser(chatID, userID uint) error {
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
			"can't exclude user %d from chat %d: %w",
			userID, chatID, err,
		)
	}

	return nil
}

// ExcludeChatTag excludes the group tag from the chat
//
// https://dev.pachca.com/members/remove-group-tag
func (c *Client) ExcludeChatTag(chatID, tagID uint) error {
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
			"can't exclude group tag %d from chat %d: %w",
			tagID, chatID, err,
		)
	}

	return nil
}

// ArchiveChat sends chat to archive
//
// https://dev.pachca.com/chats/update-archive
func (c *Client) ArchiveChat(chatID uint) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case chatID == 0:
		return ErrInvalidChatID
	}

	err := c.sendRequest(
		req.PUT, getURL("/chats/%d/archive", chatID),
		nil, nil, nil,
	)

	if err != nil {
		return fmt.Errorf("can't archive chat %d: %w", chatID, err)
	}

	return nil
}

// UnarchiveChat restores chat from archive
//
// https://dev.pachca.com/chats/update-unarchive
func (c *Client) UnarchiveChat(chatID uint) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case chatID == 0:
		return ErrInvalidChatID
	}

	err := c.sendRequest(
		req.PUT, getURL("/chats/%d/unarchive", chatID),
		nil, nil, nil,
	)

	if err != nil {
		return fmt.Errorf("can't unarchive chat %d: %w", chatID, err)
	}

	return nil
}

// MESSAGES ///////////////////////////////////////////////////////////////////////// //

// GetMessages returns messages from given chat
//
// https://dev.pachca.com/messages/list
func (c *Client) GetMessages(chatID uint, page, perPage int) (Messages, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case chatID == 0:
		return nil, ErrInvalidChatID
	case page < 1:
		return nil, ErrInvalidPageNum
	case perPage < 1 || perPage > MAX_PER_PAGE:
		return nil, ErrInvalidPerPageNum
	}

	resp := &struct {
		Data Messages `json:"data"`
	}{}

	query := req.Query{"chat_id": chatID, "page": page, "per": perPage}
	err := c.sendRequest(req.GET, getURL("/messages"), query, nil, resp)

	if err != nil {
		return nil, fmt.Errorf("can't get messages of chat with ID %d: %w", chatID, err)
	}

	return resp.Data, nil
}

// GetLatestMessages returns specified number of the latest messages from the chat
func (c *Client) GetLatestMessages(chatID uint, numMessages int) (Messages, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case chatID == 0:
		return nil, ErrInvalidChatID
	case numMessages < 1:
		return nil, ErrInvalidMessageNum
	}

	result := make(Messages, 0, numMessages)

	var perPage int

	for page := 1; page < MAX_PAGES; page++ {
		perPage = min(MAX_PER_PAGE, numMessages)

		messages, err := c.GetMessages(chatID, page, perPage)

		if err != nil {
			return nil, err
		}

		result = append(result, messages...)
		numMessages -= len(messages)

		if perPage < MAX_PER_PAGE || len(messages) < MAX_PER_PAGE {
			break
		}
	}

	return result, nil
}

// GetMessage returns info about message
//
// https://dev.pachca.com/messages/get
func (c *Client) GetMessage(messageID uint) (*Message, error) {
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
		return nil, fmt.Errorf("can't fetch message info: %w", err)
	}

	return resp.Data, nil
}

// GetMessageReads returns a slice with IDs of users who have read the message
//
// https://dev.pachca.com/read-member/list-read-member-ids
func (c *Client) GetMessageReads(messageID uint) ([]uint, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case messageID == 0:
		return nil, ErrInvalidMessageID
	}

	var result []uint

	resp := &struct {
		Data []uint `json:"data"`
	}{}

	query := req.Query{"per": 300}

	for i := 1; i < MAX_PAGES; i++ {
		query["page"] = i

		err := c.sendRequest(
			req.GET, getURL("/messages/%d/read_member_ids", messageID),
			query, nil, resp,
		)

		if err != nil {
			return nil, fmt.Errorf("can't fetch message reads info: %w", err)
		}

		result = append(result, resp.Data...)

		if len(resp.Data) != 300 {
			break
		}
	}

	return result, nil
}

// AddMessage creates new message to user or chat
//
// https://dev.pachca.com/messages/create
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
		return nil, fmt.Errorf("can't create a new message: %w", err)
	}

	return resp.Data, nil
}

// EditMessage modifies message
//
// https://dev.pachca.com/messages/update
func (c *Client) EditMessage(messageID uint, message *MessageRequest) (*Message, error) {
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
		return nil, fmt.Errorf("can't modify message %d: %w", messageID, err)
	}

	return resp.Data, nil
}

// DeleteMessage deletes message with given ID
//
// https://dev.pachca.com/messages/delete
func (c *Client) DeleteMessage(messageID uint) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case messageID == 0:
		return ErrInvalidMessageID
	}

	err := c.sendRequest(req.DELETE, getURL("/messages/%d", messageID), nil, nil, nil)

	if err != nil {
		return fmt.Errorf("can't delete message %d: %w", messageID, err)
	}

	return nil
}

// PinMessage pins message to chat
//
// https://dev.pachca.com/messages/pin
func (c *Client) PinMessage(messageID uint) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case messageID == 0:
		return ErrInvalidMessageID
	}

	err := c.sendRequest(req.POST, getURL("/messages/%d/pin", messageID), nil, nil, nil)

	if err != nil {
		return fmt.Errorf("can't pin message %d: %w", messageID, err)
	}

	return nil
}

// UnpinMessage unpins message from chat
//
// https://dev.pachca.com/messages/unpin
func (c *Client) UnpinMessage(messageID uint) error {
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
		return fmt.Errorf("can't unpin message %d: %w", messageID, err)
	}

	return nil
}

// AddLinkPreview adds link previews to message with given ID
//
// https://dev.pachca.com/link-previews/add-link-previews
func (c *Client) AddLinkPreview(messageID uint, previews LinkPreviews) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case messageID == 0:
		return ErrInvalidMessageID
	case len(previews) == 0:
		return ErrEmptyPreviews
	}

	payload := &struct {
		Previews LinkPreviews `json:"link_previews"`
	}{
		Previews: previews,
	}

	err := c.sendRequest(
		req.POST,
		getURL("/messages/%d/link_previews", messageID),
		nil, payload, nil,
	)

	if err != nil {
		return fmt.Errorf("can't add previews to message %d: %w", messageID, err)
	}

	return nil
}

// SendMessageToUser helper to send message to user with given ID
func (c *Client) SendMessageToUser(userID uint, text string) (*Message, error) {
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
func (c *Client) SendMessageToChat(chatID uint, text string) (*Message, error) {
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
func (c *Client) SendMessageToThread(threadID uint, text string) (*Message, error) {
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

// UpdateMessage helper to change message text
func (c *Client) UpdateMessage(messageID uint, text string) (*Message, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case messageID == 0:
		return nil, ErrInvalidMessageID
	case text == "":
		return nil, ErrEmptyMessage
	}

	return c.EditMessage(messageID, &MessageRequest{Content: text})
}

// DeleteMessageButtons is a helper for deleting buttons from a message
func (c *Client) DeleteMessageButtons(messageID uint) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case messageID == 0:
		return ErrInvalidMessageID
	}

	msg, err := c.GetMessage(messageID)

	if err != nil {
		return err
	}

	if len(msg.Buttons) == 0 {
		return nil
	}

	_, err = c.EditMessage(messageID, &MessageRequest{
		Content: msg.Content,
		Buttons: Buttons{},
	})

	return err
}

// THREADS ////////////////////////////////////////////////////////////////////////// //

// GetThread returns info about thread
//
// https://dev.pachca.com/thread/get
func (c *Client) GetThread(threadID uint) (*Thread, error) {
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
		return nil, fmt.Errorf("can't fetch thread info: %w", err)
	}

	return resp.Data, nil
}

// NewThread creates a new tread
//
// https://dev.pachca.com/thread/add-thread
func (c *Client) NewThread(messageID uint) (*Thread, error) {
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
		return nil, fmt.Errorf("can't create thread for message %d: %w", messageID, err)
	}

	return resp.Data, nil
}

// AddThreadMessage helper to create thread and add new message to it
func (c *Client) AddThreadMessage(messageID uint, message *MessageRequest) (*Thread, *Message, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, nil, ErrNilClient
	case messageID == 0:
		return nil, nil, ErrInvalidMessageID
	case message == nil:
		return nil, nil, ErrNilMessageRequest
	}

	thread, err := c.NewThread(messageID)

	if err != nil {
		return nil, nil, err
	}

	message.EntityID = thread.ID
	message.EntityType = ENTITY_TYPE_THREAD

	msg, err := c.AddMessage(message)

	if err != nil {
		return nil, nil, err
	}

	return thread, msg, err
}

// AddThreadMessageText helper to create thread and add new message with given text to it
func (c *Client) AddThreadMessageText(messageID uint, text string) (*Thread, *Message, error) {
	return c.AddThreadMessage(messageID, &MessageRequest{Content: text})
}

// FILES //////////////////////////////////////////////////////////////////////////// //

// UploadFile uploads new file and returns key of it
//
// https://dev.pachca.com/common/direct-url
func (c *Client) UploadFile(file string) (*File, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case file == "":
		return nil, ErrEmptyFilePath
	}

	fd, err := os.Open(file)

	if err != nil {
		return nil, fmt.Errorf("can't open file %q to upload: %w", file, err)
	}

	defer fd.Close()

	stat, err := fd.Stat()

	if err != nil {
		return nil, fmt.Errorf("can't get file %q info: %w", file, err)
	}

	if stat.Size() >= c.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds the limit (%d ≥ %d)", stat.Size(), c.MaxFileSize)
	}

	upload := &Upload{}
	err = c.sendRequest(req.POST, getURL("/uploads"), nil, nil, upload)

	if err != nil {
		return nil, fmt.Errorf("can't create upload for file %q: %w", file, err)
	}

	fileName := path.Base(fd.Name())
	key := strings.ReplaceAll(upload.Key, "${filename}", fileName)
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	contentType := mw.FormDataContentType()

	go func() {
		defer pw.Close()
		defer mw.Close()

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
			pw.CloseWithError(fmt.Errorf("can't create multipart upload: %w", errs.First()))
			return
		}

		fw, err := mw.CreateFormFile("file", fileName)

		if err != nil {
			pw.CloseWithError(err)
			return
		}

		_, err = io.Copy(fw, fd)

		if err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	resp, err := c.engine.Post(
		req.Request{
			URL:         upload.DirectURL,
			ContentType: contentType,
			Body:        pr,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("can't send request to API: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf(
			"can't upload file %q data (key: %s | status: %d): %w",
			file, upload.Key, resp.StatusCode, extractS3Error(resp.String()),
		)
	}

	return &File{
		Key:  key,
		Name: fileName,
		Size: stat.Size(),
		Type: guessFileType(fileName),
	}, nil
}

// BOTS ///////////////////////////////////////////////////////////////////////////// //

// UpdateBot updates bot webhook URL
//
// https://dev.pachca.com/bots/update
func (c *Client) UpdateBot(botID uint, webhookURL string) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case botID == 0:
		return ErrInvalidBotID
	case webhookURL == "":
		return ErrEmptyFilePath
	}

	payload := struct {
		Bot struct {
			Webhook struct {
				URL string `json:"outgoing_url"`
			} `json:"webhook"`
		} `json:"bot"`
	}{}

	payload.Bot.Webhook.URL = webhookURL

	err := c.sendRequest(req.PUT, getURL("/bots/%d", botID), nil, &payload, nil)

	if err != nil {
		return fmt.Errorf("can't update bot settings: %w", err)
	}

	return nil
}

// GetWebhookEvents returns webhook events. Each event's Payload field contains
// raw JSON that must be decoded using webhook.DecodeJSON to extract the specific
// webhook type (Message, Reaction, etc.).
//
// Example:
//
//	events, _ := client.GetWebhookEvents(10)
//	for _, ev := range events {
//	    wh, _ := webhook.DecodeJSON(ev.Payload)
//	    // Handle webhook based on type
//	}
//
// https://dev.pachca.com/bots/list-events
func (c *Client) GetWebhookEvents(maxPages int) ([]*WebhookEvent, error) {
	switch {
	case c == nil || c.engine == nil:
		return nil, ErrNilClient
	case maxPages < 1:
		return nil, ErrInvalidMaxPages
	}

	var result []*WebhookEvent
	query := req.Query{}

	for i := 0; i < min(maxPages, MAX_PAGES); i++ {
		resp := &struct {
			Data []*WebhookEvent `json:"data"`
			Meta *Metadata       `json:"meta"`
		}{}

		err := c.sendRequest(req.GET, getURL("/webhooks/events"), query, nil, resp)

		if err != nil {
			return nil, fmt.Errorf("can't fetch webhook events: %w", err)
		}

		result = append(result, resp.Data...)

		if len(resp.Data) != c.getBatchSize() {
			break
		}

		query.SetIf(
			resp.Meta != nil && resp.Meta.Paginate != nil,
			"cursor", resp.Meta.Paginate.NextPage,
		)
	}

	return result, nil
}

// DeleteWebhookEvent deletes webhook event with given ID
//
// https://dev.pachca.com/bots/remove-event
func (c *Client) DeleteWebhookEvent(eventID string) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case eventID == "" || len(eventID) != 26:
		return ErrInvalidEventID
	}

	err := c.sendRequest(
		req.DELETE, getURL("/webhooks/events/%s", eventID),
		nil, nil, nil,
	)

	if err != nil {
		return fmt.Errorf("can't delete webhook event with ID %s: %w", eventID, err)
	}

	return nil
}

// FORMS //////////////////////////////////////////////////////////////////////////// //

// OpenView opens view with form
//
// https://dev.pachca.com/views/create-open
func (c *Client) OpenView(view *ViewRequest) error {
	switch {
	case c == nil || c.engine == nil:
		return ErrNilClient
	case view == nil:
		return ErrNilViewRequest
	case view.View == nil:
		return ErrNilView
	case view.TriggerID == "":
		return ErrEmptyTriggerID
	case view.Type != VIEW_MODAL:
		return fmt.Errorf("unknown form type %q", view.Type)
	case len(view.View.Blocks) == 0:
		return ErrViewHasNoBlocks
	}

	for _, b := range view.View.Blocks {
		b.Init()
	}

	err := c.sendRequest(req.POST, getURL("/views/open"), nil, view, nil)

	if err != nil {
		return fmt.Errorf("can't open view: %w", err)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Get returns custom property with given ID
func (p Properties) Get(id uint) *Property {
	for _, pp := range p {
		if pp.ID == id {
			return pp
		}
	}

	return nil
}

// Has returns true if properties contains property with given name
func (p Properties) Has(name string) bool {
	return len(p) > 0 && p.Find(name) != nil
}

// HasAny returns true if properties contains property with one of given names
func (p Properties) HasAny(name ...string) bool {
	return len(p) > 0 && p.FindAny(name...) != nil
}

// Find returns custom property with given name
func (p Properties) Find(name string) *Property {
	name = strings.ToLower(name)

	for _, pp := range p {
		if strings.ToLower(pp.Name) == name {
			return pp
		}
	}

	return nil
}

// FindAny returns first found property with one of given names
func (p Properties) FindAny(name ...string) *Property {
	for _, n := range name {
		p := p.Find(n)

		if p != nil {
			return p
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

// IsSet returns true if property has a value
func (p *Property) IsSet() bool {
	return p != nil && p.Value != ""
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
		return time.Time{}, ErrNilPropertyRequest
	case p.Value == "":
		return time.Time{}, nil
	case p.Type != PROP_TYPE_DATE:
		return time.Time{}, fmt.Errorf("invalid property type for date (%s)", p.Type)
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
		return 0, ErrNilPropertyRequest
	case p.Value == "":
		return 0, nil
	case p.Type != PROP_TYPE_NUMBER:
		return 0, fmt.Errorf("invalid property type for date (%s)", p.Type)
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
	switch {
	case u == nil, u.FirstName == "" && u.LastName == "":
		return ""
	case u.FirstName == "" && u.LastName != "":
		return u.LastName
	case u.FirstName != "" && u.LastName == "":
		return u.FirstName
	}

	return u.FirstName + " " + u.LastName
}

// HasAvatar returns true if user has custom avatar
func (u *User) HasAvatar() bool {
	return u != nil && u.ImageURL != ""
}

// HasTag returns true if user has given tag
func (u *User) HasTag(tag string) bool {
	return u != nil && slices.Contains(u.Tags, tag)
}

// IsActive returns true if user is active
func (u *User) IsActive() bool {
	return u != nil && !u.IsSuspended && u.InviteStatus == INVITE_CONFIRMED
}

// IsInvited returns true if user is invited
func (u *User) IsInvited() bool {
	return u != nil && !u.IsSuspended && u.InviteStatus == INVITE_SENT
}

// IsGuest returns true if user is guest or multi-guest
func (u *User) IsGuest() bool {
	return u != nil && (u.Role == ROLE_MULTI_GUEST || u.Role == ROLE_GUEST)
}

// IsMultiGuest returns true if user is multi-guest
func (u *User) IsMultiGuest() bool {
	return u != nil && u.Role == ROLE_MULTI_GUEST
}

// IsPaid returns true if user is paid (not bot or guest)
func (u *User) IsPaid() bool {
	return u != nil && (u.Role == ROLE_MULTI_GUEST || u.Role == ROLE_REGULAR || u.Role == ROLE_ADMIN)
}

// IsAdmin returns true if user is admin
func (u *User) IsAdmin() bool {
	return u != nil && u.Role == ROLE_ADMIN
}

// IsRegular returns true if user just regular user
func (u *User) IsRegular() bool {
	return u != nil && u.Role == ROLE_REGULAR
}

// Get returns user with given ID or nil
func (u Users) Get(id uint) *User {
	for _, uu := range u {
		if uu.ID == id {
			return uu
		}
	}

	return nil
}

// InChat only returns users that are present in the given chat
func (u Users) InChat(chat *Chat) Users {
	if chat == nil {
		return nil
	}

	var result Users

	for _, id := range chat.Members {
		user := u.Get(id)

		if user != nil {
			result = append(result, user)
		}
	}

	return result
}

// Find returns user with given nickname or email
func (u Users) Find(nicknameOrEmail string) *User {
	nicknameOrEmail = strings.ToLower(nicknameOrEmail)

	for _, uu := range u {
		if strings.ToLower(uu.Nickname) == nicknameOrEmail ||
			strings.ToLower(uu.Email) == nicknameOrEmail {
			return uu
		}
	}

	return nil
}

// Active returns slice with active users
func (u Users) Active() Users {
	return sliceutil.Filter(u, func(uu *User, _ int) bool {
		return uu.IsActive()
	})
}

// Suspended returns slice with inactive users
func (u Users) Suspended() Users {
	return sliceutil.Filter(u, func(uu *User, _ int) bool {
		return uu.IsSuspended
	})
}

// Invited returns all invited users
func (u Users) Invited() Users {
	return sliceutil.Filter(u, func(uu *User, _ int) bool {
		return uu.IsInvited()
	})
}

// Bots returns slice with bots
func (u Users) Bots() Users {
	return sliceutil.Filter(u, func(uu *User, _ int) bool {
		return uu.IsBot
	})
}

// People returns slice with non-bot users
func (u Users) People() Users {
	return sliceutil.Filter(u, func(uu *User, _ int) bool {
		return !uu.IsBot
	})
}

// Admins returns slice with admins
func (u Users) Admins() Users {
	return sliceutil.Filter(u, func(uu *User, _ int) bool {
		return uu.IsAdmin()
	})
}

// Regular returns slice with regular users
func (u Users) Regular() Users {
	return sliceutil.Filter(u, func(uu *User, _ int) bool {
		return uu.IsRegular()
	})
}

// Guests returns slice with guests or multi-guests
func (u Users) Guests() Users {
	return sliceutil.Filter(u, func(uu *User, _ int) bool {
		return uu.IsGuest()
	})
}

// MultiGuests returns slice with multi-guests
func (u Users) MultiGuests() Users {
	return sliceutil.Filter(u, func(uu *User, _ int) bool {
		return uu.IsMultiGuest()
	})
}

// Paid returns slice with paid users
func (u Users) Paid() Users {
	return sliceutil.Filter(u, func(uu *User, _ int) bool {
		return uu.IsPaid()
	})
}

// WithoutGuests returns slice with users without guests
func (u Users) WithoutGuests() Users {
	return sliceutil.Filter(u, func(uu *User, _ int) bool {
		return !uu.IsGuest()
	})
}

// WithTag returns users with given tag
func (u Users) WithTag(tag string) Users {
	return sliceutil.Filter(u, func(uu *User, _ int) bool {
		return uu.HasTag(tag)
	})
}

// Get returns chat with given ID
func (c Chats) Get(id uint) *Chat {
	for _, cc := range c {
		if cc.ID == id {
			return cc
		}
	}

	return nil
}

// Find returns chat with given name
func (c Chats) Find(name string) *Chat {
	name = strings.ToLower(name)

	for _, cc := range c {
		if strings.ToLower(cc.Name) == name {
			return cc
		}
	}

	return nil
}

// Public returns slice with public chats
func (c Chats) Public() Chats {
	return sliceutil.Filter(c, func(cc *Chat, _ int) bool {
		return cc.IsPublic
	})
}

// Channels returns slice with channels
func (c Chats) Channels() Chats {
	return sliceutil.Filter(c, func(cc *Chat, _ int) bool {
		return cc.IsChannel
	})
}

// Personal returns p2p chats
func (c Chats) Personal() Chats {
	return sliceutil.Filter(c, func(cc *Chat, _ int) bool {
		return cc.Name == ""
	})
}

// Communal returns communal chats (non-p2p)
func (c Chats) Communal() Chats {
	return sliceutil.Filter(c, func(cc *Chat, _ int) bool {
		return cc.Name != ""
	})
}

// Get returns tag with given ID
func (t Tags) Get(id uint) *Tag {
	for _, tt := range t {
		if tt.ID == id {
			return tt
		}
	}

	return nil
}

// Find returns tag with given name
func (t Tags) Find(name string) *Tag {
	name = strings.ToLower(name)

	for _, tt := range t {
		if strings.ToLower(tt.Name) == name {
			return tt
		}
	}

	return nil
}

// Names returns names of all tags
func (t Tags) Names() []string {
	var result []string

	for _, tt := range t {
		result = append(result, tt.Name)
	}

	return result
}

// InChat only returns tags that are present in the given chat
func (t Tags) InChat(chat *Chat) Tags {
	if chat == nil {
		return nil
	}

	var result Tags

	for _, id := range chat.GroupTags {
		tag := t.Get(id)

		if tag != nil {
			result = append(result, tag)
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

// AddBlock adds new blocks to the view
func (v *View) AddBlocks(blocks ...block.Block) *View {
	if v == nil || len(blocks) == 0 {
		return v
	}

	for _, b := range blocks {
		if b != nil {
			b.Init()
			v.Blocks = append(v.Blocks, b)
		}
	}

	return v
}

// AddBlocksIf conditionally adds new blocks to the view
func (v *View) AddBlocksIf(cond bool, blocks ...block.Block) *View {
	if !cond {
		return v
	}

	return v.AddBlocks(blocks...)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// isZero is special method for omitzero
func (f Files) isZero() bool {
	return f == nil
}

// isZero is special method for omitzero
func (b Buttons) isZero() bool {
	return b == nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Error returns error message
func (e *S3Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Message
}

// ////////////////////////////////////////////////////////////////////////////////// //

// ToQuery converts filter struct to request query
func (f ChatFilter) ToQuery() req.Query {
	query := req.Query{}

	query.SetIf(f.Public, "availability", "public")
	query.SetIf(!f.LastMessageBefore.IsZero(),
		"last_message_at_before", formatDate(f.LastMessageBefore),
	)
	query.SetIf(!f.LastMessageAfter.IsZero(),
		"last_message_at_after", formatDate(f.LastMessageAfter),
	)

	if len(f.Sort) != 0 {
		for k, v := range f.Sort {
			query["sort["+k+"]"] = v
		}
	}

	return query
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getBatchSize returns batch size for paginated responses
func (c *Client) getBatchSize() int {
	return mathutil.Between(c.BatchSize, 5, MAX_PER_PAGE)
}

// sendRequest sends request to Pachca API
func (c *Client) sendRequest(method, url string, query req.Query, payload any, response any) error {
	r := req.Request{
		Method: method,
		URL:    url,
		Query:  query,
		Accept: req.CONTENT_TYPE_JSON,
		Auth:   req.AuthBearer{c.token},
	}

	if payload != nil {
		r.ContentType = req.CONTENT_TYPE_JSON
		r.Body = payload
	}

	resp, err := c.engine.Do(r)

	if err != nil {
		return fmt.Errorf("can't send request to API: %w", err)
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
			return fmt.Errorf("can't decode API response: %w", err)
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

// extractS3Error extracts error text from S3 error message
func extractS3Error(errorMessage string) error {
	found := s3ErrorExtractRegex.FindStringSubmatch(errorMessage)

	if len(found) == 2 {
		return &S3Error{
			Message: found[1],
			Full:    errorMessage,
		}
	}

	return &S3Error{
		Message: errorMessage,
		Full:    errorMessage,
	}
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
