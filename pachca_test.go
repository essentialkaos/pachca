package pachca

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"testing"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

type PachcaSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&PachcaSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *PachcaSuite) TestTokenValidator(c *C) {
	c.Assert(ValidateToken(""), NotNil)
	c.Assert(ValidateToken("ABCD"), NotNil)
	c.Assert(ValidateToken("YQlf-6Vce7jM1RMZZUs_iWKYPt24PeR4c7k_RwzqjI5"), IsNil)
}

func (s *PachcaSuite) TestNewClientErrors(c *C) {
	cc, err := NewClient("test")
	c.Assert(cc, IsNil)
	c.Assert(err, NotNil)

	cc, err = NewClient("YQlf-6Vce7jM1RMZZUs_iWKYPt24PeR4c7k_RwzqjI5")
	c.Assert(cc, NotNil)
	c.Assert(err, IsNil)
	cc.SetUserAgent("Test", "1.0.0")
	c.Assert(cc.Engine(), NotNil)
}

func (s *PachcaSuite) TestNilClient(c *C) {
	var cc *Client

	cc.SetUserAgent("Test", "1.0.0")

	c.Assert(cc.Engine(), IsNil)

	// CUSTOM PROPS

	_, err := cc.GetProperties()
	c.Assert(err, Equals, ErrNilClient)

	// REACTIONS

	_, err = cc.GetReactions(1)
	c.Assert(err, Equals, ErrNilClient)

	c.Assert(cc.AddReaction(1, "😊"), Equals, ErrNilClient)

	c.Assert(cc.DeleteReaction(1, "😊"), Equals, ErrNilClient)

	// USERS

	_, err = cc.GetUser(1)
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.GetUsers()
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.AddUser(&UserRequest{Email: "test@example.com"})
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.EditUser(1, &UserRequest{Email: "test@example.com"})
	c.Assert(err, Equals, ErrNilClient)

	c.Assert(cc.DeleteUser(1), Equals, ErrNilClient)

	// TAGS

	_, err = cc.GetTags()
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.GetTag(1)
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.GetTagUsers(1)
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.AddTag("test")
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.EditTag(1, "test")
	c.Assert(err, Equals, ErrNilClient)

	c.Assert(cc.DeleteTag(1), Equals, ErrNilClient)

	// CHATS

	_, err = cc.GetChats()
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.GetChat(1)
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.AddChat(&ChatRequest{Name: "Test"})
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.EditChat(1, &ChatRequest{Name: "Test"})
	c.Assert(err, Equals, ErrNilClient)

	c.Assert(cc.AddChatUsers(1, []ID{1, 2, 3}, true), Equals, ErrNilClient)
	c.Assert(cc.AddChatTags(1, []ID{1, 2, 3}), Equals, ErrNilClient)
	c.Assert(cc.ExcludeChatUser(1, 1), Equals, ErrNilClient)
	c.Assert(cc.ExcludeChatTag(1, 1), Equals, ErrNilClient)

	// MESSAGES

	_, err = cc.GetMessage(1)
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.AddMessage(&MessageRequest{EntityID: 1})
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.EditMessage(1, &MessageRequest{EntityID: 1})
	c.Assert(err, Equals, ErrNilClient)

	c.Assert(cc.DeleteMessage(1), Equals, ErrNilClient)
	c.Assert(cc.PinMessage(1), Equals, ErrNilClient)
	c.Assert(cc.UnpinMessage(1), Equals, ErrNilClient)

	_, err = cc.SendMessageToUser(1, "Test")
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.SendMessageToChat(1, "Test")
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.SendMessageToThread(1, "Test")
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.ChangeMessageText(1, "Test")
	c.Assert(err, Equals, ErrNilClient)

	// THREADS

	_, err = cc.GetThread(1)
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.NewThread(1)
	c.Assert(err, Equals, ErrNilClient)

	_, err = cc.AddThreadMessage(1, &MessageRequest{EntityID: 1})
	c.Assert(err, Equals, ErrNilClient)

	// FILES

	_, err = cc.UploadFile("test.txt")
	c.Assert(err, Equals, ErrNilClient)
}

func (s *PachcaSuite) TestErrors(c *C) {
	cc, err := NewClient("YQlf-6Vce7jM1RMZZUs_iWKYPt24PeR4c7k_RwzqjI5")
	c.Assert(cc, NotNil)
	c.Assert(err, IsNil)

	// REACTIONS

	_, err = cc.GetReactions(0)
	c.Assert(err, Equals, ErrInvalidMessageID)

	c.Assert(cc.AddReaction(0, "😊"), Equals, ErrInvalidMessageID)
	c.Assert(cc.AddReaction(1, ""), Equals, ErrBlankEmoji)

	c.Assert(cc.DeleteReaction(0, "😊"), Equals, ErrInvalidMessageID)
	c.Assert(cc.DeleteReaction(1, ""), Equals, ErrBlankEmoji)

	// USERS

	_, err = cc.GetUser(0)
	c.Assert(err, Equals, ErrInvalidUserID)

	_, err = cc.AddUser(nil)
	c.Assert(err, Equals, ErrNilUserRequest)
	_, err = cc.AddUser(&UserRequest{})
	c.Assert(err, Equals, ErrEmptyUserEmail)

	_, err = cc.EditUser(0, nil)
	c.Assert(err, Equals, ErrInvalidUserID)
	_, err = cc.EditUser(1, nil)
	c.Assert(err, Equals, ErrNilUserRequest)

	c.Assert(cc.DeleteUser(0), Equals, ErrInvalidUserID)

	// TAGS

	_, err = cc.GetTag(0)
	c.Assert(err, Equals, ErrInvalidTagID)

	_, err = cc.GetTagUsers(0)
	c.Assert(err, Equals, ErrInvalidTagID)

	_, err = cc.AddTag("")
	c.Assert(err, Equals, ErrEmptyTag)

	_, err = cc.EditTag(0, "test")
	c.Assert(err, Equals, ErrInvalidTagID)
	_, err = cc.EditTag(1, "")
	c.Assert(err, Equals, ErrEmptyTag)

	c.Assert(cc.DeleteTag(0), Equals, ErrInvalidTagID)

	// CHATS

	_, err = cc.GetChat(0)
	c.Assert(err, Equals, ErrInvalidChatID)

	_, err = cc.AddChat(nil)
	c.Assert(err, Equals, ErrNilChatRequest)
	_, err = cc.AddChat(&ChatRequest{})
	c.Assert(err, Equals, ErrEmptyChatName)

	_, err = cc.EditChat(0, nil)
	c.Assert(err, Equals, ErrInvalidChatID)
	_, err = cc.EditChat(1, nil)
	c.Assert(err, Equals, ErrNilChatRequest)

	c.Assert(cc.AddChatUsers(0, []ID{1, 2, 3}, true), Equals, ErrInvalidChatID)
	c.Assert(cc.AddChatUsers(1, nil, true), Equals, ErrEmptyUsersIDS)

	c.Assert(cc.AddChatTags(0, []ID{1, 2, 3}), Equals, ErrInvalidChatID)
	c.Assert(cc.AddChatTags(1, nil), Equals, ErrEmptyTagsIDS)

	c.Assert(cc.ExcludeChatUser(0, 1), Equals, ErrInvalidChatID)
	c.Assert(cc.ExcludeChatUser(1, 0), Equals, ErrInvalidUserID)

	c.Assert(cc.ExcludeChatTag(0, 1), Equals, ErrInvalidChatID)
	c.Assert(cc.ExcludeChatTag(1, 0), Equals, ErrInvalidTagID)

	// MESSAGES

	_, err = cc.GetMessage(0)
	c.Assert(err, Equals, ErrInvalidMessageID)

	_, err = cc.AddMessage(nil)
	c.Assert(err, Equals, ErrNilMessageRequest)
	_, err = cc.AddMessage(&MessageRequest{})
	c.Assert(err, Equals, ErrInvalidEntityID)

	_, err = cc.EditMessage(0, nil)
	c.Assert(err, Equals, ErrInvalidMessageID)
	_, err = cc.EditMessage(1, nil)
	c.Assert(err, Equals, ErrNilMessageRequest)

	c.Assert(cc.DeleteMessage(0), Equals, ErrInvalidMessageID)
	c.Assert(cc.PinMessage(0), Equals, ErrInvalidMessageID)
	c.Assert(cc.UnpinMessage(0), Equals, ErrInvalidMessageID)

	_, err = cc.SendMessageToUser(0, "test")
	c.Assert(err, Equals, ErrInvalidUserID)
	_, err = cc.SendMessageToUser(1, "")
	c.Assert(err, Equals, ErrEmptyMessage)

	_, err = cc.SendMessageToChat(0, "test")
	c.Assert(err, Equals, ErrInvalidChatID)
	_, err = cc.SendMessageToChat(1, "")
	c.Assert(err, Equals, ErrEmptyMessage)

	_, err = cc.SendMessageToThread(0, "test")
	c.Assert(err, Equals, ErrInvalidThreadID)
	_, err = cc.SendMessageToThread(1, "")
	c.Assert(err, Equals, ErrEmptyMessage)

	_, err = cc.ChangeMessageText(0, "test")
	c.Assert(err, Equals, ErrInvalidMessageID)
	_, err = cc.ChangeMessageText(1, "")
	c.Assert(err, Equals, ErrEmptyMessage)

	// THREADS

	_, err = cc.GetThread(0)
	c.Assert(err, Equals, ErrInvalidThreadID)

	_, err = cc.NewThread(0)
	c.Assert(err, Equals, ErrInvalidMessageID)

	_, err = cc.AddThreadMessage(0, &MessageRequest{})
	c.Assert(err, Equals, ErrInvalidMessageID)
	_, err = cc.AddThreadMessage(1, nil)
	c.Assert(err, Equals, ErrNilMessageRequest)

	// FILES

	_, err = cc.UploadFile("")
	c.Assert(err, Equals, ErrEmptyFilePath)
}

func (s *PachcaSuite) TestPropertiesHelpers(c *C) {
	p := Properties{
		{ID: 1, Type: PROP_TYPE_DATE, Name: "test1", Value: "2024-08-08T09:11:50.368Z"},
		{ID: 2, Type: PROP_TYPE_LINK, Name: "test2", Value: "https://domain.com"},
		{ID: 3, Type: PROP_TYPE_NUMBER, Name: "test3", Value: "314"},
		{ID: 4, Type: PROP_TYPE_TEXT, Name: "test4", Value: "Test1"},
		{ID: 5, Type: PROP_TYPE_NUMBER, Name: "test5", Value: ""},
		{ID: 6, Type: PROP_TYPE_DATE, Name: "test6", Value: ""},
	}

	c.Assert(p.Get("test"), IsNil)
	c.Assert(p.Get("test1"), NotNil)

	c.Assert(p.GetS("test2"), Equals, "https://domain.com")
	c.Assert(p.GetS("test4"), Equals, "Test1")
	c.Assert(p.GetS("test100"), Equals, "")

	i, err := p.GetI("test3")
	c.Assert(err, IsNil)
	c.Assert(i, Equals, 314)
	i, err = p.GetI("test5")
	c.Assert(err, IsNil)
	c.Assert(i, Equals, 0)
	_, err = p.GetI("test4")
	c.Assert(err, NotNil)
	_, err = p.GetI("test100")
	c.Assert(err, NotNil)

	d, err := p.GetD("test1")
	c.Assert(err, IsNil)
	c.Assert(d.IsZero(), Equals, false)
	d, err = p.GetD("test6")
	c.Assert(err, IsNil)
	_, err = p.GetD("test4")
	c.Assert(err, NotNil)
	_, err = p.GetD("test100")
	c.Assert(err, NotNil)
}

func (s *PachcaSuite) TestUsersHelpers(c *C) {
	var u *User
	c.Assert(u.FullName(), Equals, "")
	u = &User{ID: 1234, FirstName: "John", LastName: "Doe", Nickname: "j.doe"}
	c.Assert(u.FullName(), Equals, "John Doe")

	uu := Users{
		{ID: 1, IsSuspended: true, InviteStatus: INVITE_SENT, IsBot: false, Role: ROLE_USER},
		{ID: 2, IsSuspended: false, InviteStatus: INVITE_CONFIRMED, IsBot: false, Role: ROLE_USER},
		{ID: 3, IsSuspended: false, InviteStatus: INVITE_CONFIRMED, IsBot: false, Role: ROLE_ADMIN},
		{ID: 4, IsSuspended: false, InviteStatus: INVITE_CONFIRMED, IsBot: false, Role: ROLE_MULTI_GUEST},
		{ID: 5, IsSuspended: false, InviteStatus: INVITE_CONFIRMED, IsBot: true, Role: ROLE_USER},
	}

	c.Assert(uu.Active(), HasLen, 4)

	c.Assert(uu.Invited(), HasLen, 1)
	c.Assert(uu.Invited()[0].ID, Equals, ID(1))
	c.Assert(uu.Bots()[0].ID, Equals, ID(5))
	c.Assert(uu.Admins()[0].ID, Equals, ID(3))
	c.Assert(uu.Regular(), HasLen, 3)
	c.Assert(uu.Regular()[0].ID, Equals, ID(1))
	c.Assert(uu.Guests()[0].ID, Equals, ID(4))
}

func (s *PachcaSuite) TestChatsHelpers(c *C) {
	cc := Chats{
		{ID: 1, Name: "test1", IsPublic: false, IsChannel: false},
		{ID: 2, Name: "test2", IsPublic: false, IsChannel: false},
		{ID: 3, Name: "test3", IsPublic: true, IsChannel: false},
		{ID: 4, Name: "test4", IsPublic: false, IsChannel: true},
	}

	c.Assert(cc.Get("test"), IsNil)
	c.Assert(cc.Get("test1"), NotNil)

	c.Assert(cc.Public()[0].ID, Equals, ID(3))
	c.Assert(cc.Channels()[0].ID, Equals, ID(4))
}

func (s *PachcaSuite) TestURLHelpers(c *C) {
	var user *User
	var chat *Chat
	var message *Message
	var thread *Thread

	c.Assert(user.URL(), Equals, "")
	c.Assert(chat.URL(), Equals, "")
	c.Assert(message.URL(), Equals, "")
	c.Assert(thread.URL(), Equals, "")

	user = &User{ID: 89, FirstName: "John", LastName: "Doe", Nickname: "j.doe"}
	chat = &Chat{ID: 15, Name: "test1", IsPublic: false, IsChannel: false}
	message = &Message{ID: 145, ChatID: 15, Content: "Test"}
	thread = &Thread{ID: 238, ChatID: 15, MessageID: 145}

	c.Assert(user.URL(), Equals, "https://app.pachca.com/chats?user_id=89")
	c.Assert(chat.URL(), Equals, "https://app.pachca.com/chats/15")
	c.Assert(message.URL(), Equals, "https://app.pachca.com/chats/15?message=145")
	c.Assert(thread.URL(), Equals, "https://app.pachca.com/chats?thread_message_id=145&sidebar_message=238")
}

func (s *PachcaSuite) TestWebhookHelpers(c *C) {
	message := &Webhook{Type: WEBHOOK_TYPE_MESSAGE, Content: "/find-user j.doe"}
	reaction := &Webhook{Type: WEBHOOK_TYPE_REACTION}
	button := &Webhook{Type: WEBHOOK_TYPE_BUTTON}
	new := &Webhook{Type: WEBHOOK_TYPE_MESSAGE, Event: WEBHOOK_EVENT_NEW}
	update := &Webhook{Type: WEBHOOK_TYPE_MESSAGE, Event: WEBHOOK_EVENT_UPDATE}
	delete := &Webhook{Type: WEBHOOK_TYPE_MESSAGE, Event: WEBHOOK_EVENT_DELETE}

	var nilWebhook *Webhook

	c.Assert(nilWebhook.IsMessage(), Equals, false)
	c.Assert(nilWebhook.IsReaction(), Equals, false)
	c.Assert(nilWebhook.IsButton(), Equals, false)
	c.Assert(nilWebhook.IsNew(), Equals, false)
	c.Assert(nilWebhook.IsUpdate(), Equals, false)
	c.Assert(nilWebhook.IsDelete(), Equals, false)
	c.Assert(nilWebhook.Command(), Equals, "")

	c.Assert(message.IsMessage(), Equals, true)
	c.Assert(reaction.IsReaction(), Equals, true)
	c.Assert(button.IsButton(), Equals, true)
	c.Assert(new.IsNew(), Equals, true)
	c.Assert(update.IsUpdate(), Equals, true)
	c.Assert(delete.IsDelete(), Equals, true)
	c.Assert(message.Command(), Equals, "find-user")
}
