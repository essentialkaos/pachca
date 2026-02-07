## Changelog

### [0.21.1](https://kaos.sh/pachca/0.21.1)

- Code refactoring

### [0.21.0](https://kaos.sh/pachca/0.21.0)

- Improved method `UploadFile` to pipe data instead of buffering
- Added helper `GetLatestMessages`

### [0.20.2](https://kaos.sh/pachca/0.20.2)

- **`[block]`** Added a quick hack for parsing date block values, to support different date layouts

### [0.20.1](https://kaos.sh/pachca/0.20.1)

- Updated method `GetTags` to support filtering tags

### [0.20.0](https://kaos.sh/pachca/0.20.0)

- Added method `GetWebhookEvents`
- Added method `DeleteWebhookEvent`
- **`[webhook]`** Added method `DecodeJSON` to decode webhook data without age validation
- **`[webhook]`** Method `Decode` renamed to `DecodeBytes`
- Code refactoring

### [0.19.0](https://kaos.sh/pachca/0.19.0)

- Added method `View.AddBlocksIf`
- `View.AddBlocks` now returns pointer to view
- Ignore `nil` blocks in `View.AddBlocks`
- **`[block]`** Added method `Select.AddOptionIf`
- **`[block]`** Added method `Radio.AddOptionIf`
- **`[block]`** Added method `Checkbox.AddOptionIf`
- **`[block]`** Added method `Date.SetIf`
- **`[block]`** Added method `Date.SetWithDateIf`
- **`[block]`** Added method `Time.SetIf`
- **`[block]`** Added method `Time.SetWithDateIf`

### [0.18.0](https://kaos.sh/pachca/0.18.0)

- **`[webhook]`** Added helper `Message.Command`
- Added helper `User.HasTag`
- Added helper `Users.WithTag`

### [0.17.0](https://kaos.sh/pachca/0.17.0)

- Added sub-package for blocks data

### [0.16.0](https://kaos.sh/pachca/0.16.0)

- Added helper `DeleteMessageButtons`
- Dependencies update

### [0.15.1](https://kaos.sh/pachca/0.15.1)

- Improved working with paginated methods
- Dependencies update

### [0.15.0](https://kaos.sh/pachca/0.15.0)

- Added `UpdateBot` method
- `GetChats` migrated to cursor pagination
- Improved cursor pagination for `GetChatUsers`
- Dependencies update

### [0.14.0](https://kaos.sh/pachca/0.14.0)

- Refactored working with webhooks

### [0.13.0](https://kaos.sh/pachca/0.13.0)

- Added [forms](https://crm.pachca.com/dev/forms/getting-started/) support
- Webhook support has been rewritten from scratch
- Dependencies update

### [0.12.0](https://kaos.sh/pachca/0.12.0)

- Added custom reactions support to `AddReaction` and `DeleteReaction`
- Added name of reaction to `Reaction`

### [0.11.0](https://kaos.sh/pachca/0.11.0)

- Improved message update helper (`UpdateMessage`)
- Fixed issue with addition buttons to message
- Fixed type of parent message ID in `MessageRequest`

### [0.10.2](https://kaos.sh/pachca/0.10.2)

- Added helper `User.IsMultiGuest`
- Added helper `User.IsPaid`
- Added helper `Users.MultiGuests`
- Added helper `Users.Paid`
- Code refactoring

### [0.10.1](https://kaos.sh/pachca/0.10.1)

- Code refactoring
- Dependencies update

### [0.10.0](https://kaos.sh/pachca/0.10.0)

- Added new API method `GetChatUsers`
- Added `Name` field to `Webhook` struct
- Improved method `GetMessageReads`
- Added chat roles `CHAT_ROLE_OWNER` and `CHAT_ROLE_ANY`

### [0.9.0](https://kaos.sh/pachca/0.9.0)

- Added method `ReadWebhook` for reading webhook data from `http.Request`
- Added method `ReadWebhookSigned` for reading signed webhook data from `http.Request`
- Dependencies update

### [0.8.2](https://kaos.sh/pachca/0.8.2)

- Dependencies update

### [0.8.1](https://kaos.sh/pachca/0.8.1)

- Added new field `IsPublic` to `Chat`
- Dependencies update

### [0.8.0](https://kaos.sh/pachca/0.8.0)

- Added new API method `GetMessages`

### [0.7.0](https://kaos.sh/pachca/0.7.0)

- Added method `CurrentUser` for fetching info about current user

### [0.6.3](https://kaos.sh/pachca/0.6.3)

- Added support of sorting chats using `ChatFilter`

### [0.6.2](https://kaos.sh/pachca/0.6.2)

- Type of `File.Size` changed from `uint` to `int64`
- Type of `File.Width` changed from `uint` to `int`
- Type of `File.Height` changed from `uint` to `int`

### [0.6.1](https://kaos.sh/pachca/0.6.1)

- Updated compatibility with the latest version of API
- Dependencies update

### [0.6.0](https://kaos.sh/pachca/0.6.0)

- `AddThreadMessage` and `AddThreadMessageText` now also returns `Thread`
- `ChangeMessageText` renamed to `UpdateMessage`

### [0.5.5](https://kaos.sh/pachca/0.5.5)

- Dependencies update
- Set custom user-agent by default

### [0.5.4](https://kaos.sh/pachca/0.5.4)

- Minor fixes in JSON encoding of structs
- Better errors

### [0.5.3](https://kaos.sh/pachca/0.5.3)

- More info in error messages from S3
- Dependencies update

### [0.5.2](https://kaos.sh/pachca/0.5.2)

- Removed debug output

### [0.5.1](https://kaos.sh/pachca/0.5.1)

- Added helper `Users.WithoutGuests`

### [0.5.0](https://kaos.sh/pachca/0.5.0)

- Added method `AddLinkPreview`
- Dependencies update

### [0.4.0](https://kaos.sh/pachca/0.4.0)

- Added method [`GetMessageReads`](https://crm.pachca.com/dev/read_members/list/)
- Added method [`SetChatUserRole`](https://crm.pachca.com/dev/members/users/update/)

### [0.3.2](https://kaos.sh/pachca/0.3.2)

- Dependencies update
- Code refactoring

### [0.3.1](https://kaos.sh/pachca/0.3.1)

- Minor fix with deferring response body close
- Dependencies update

### [0.3.0](https://kaos.sh/pachca/0.3.0)

- Added new API method `ArchiveChat`
- Added new API method `UnarchiveChat`
- Dependencies update

### [0.2.3](https://kaos.sh/pachca/0.2.3)

- Make encoding of all variables of all request optional

### [0.2.2](https://kaos.sh/pachca/0.2.2)

- Fixed bug with using wrong HTTP method for `AddChatUsers`
- Minor typos fixes

### [0.2.1](https://kaos.sh/pachca/0.2.1)

- Added `Users.People` helper

### [0.2.0](https://kaos.sh/pachca/0.2.0)

- Added helper `Property.IsSet`
- Added helper `AddThreadMessageText`
- Dependencies update

### [0.1.1](https://kaos.sh/pachca/0.1.1)

- Added helper `Tags.Names`
- `Properties` replaced by `PropertiesRequest` for `UserRequest`
- All `Find` helpers now case-insensitive
- Improved `APIError` message format

### [0.1.0](https://kaos.sh/pachca/0.1.0)

- Added helper `User.IsActive`
- Added helper `User.IsInvited`
- Added helper `User.IsGuest`
- Added helper `User.IsAdmin`
- Added helper `User.IsRegular`
- Added helper `User.HasAvatar`
- Added helper `Users.InChat`
- Added helper `Chats.Personal`
- Added helper `Chats.Communal`
- Added helper `Tags.Get`
- Added helper `Tags.InChat`
- Added helper `Properties.Get`
- Added helper `Properties.Has`
- Added helper `Properties.HasAny`
- `Chat.GroupTagIDs` renamed to `Chat.GroupTags`
- `Properties.Get` renamed to `Properties.Find`
- `Properties.GetAny` renamed to `Properties.FindAny`
- `ROLE_USER` renamed to `ROLE_REGULAR`
- Removed bots from `Chats.Invited` output
- Code refactoring

### [0.0.4](https://kaos.sh/pachca/0.0.4)

- Improved `User.FullName` helper
- Improved `Chats.Suspended` helper

### [0.0.3](https://kaos.sh/pachca/0.0.3)

- Added helper `Chats.Suspended`

### [0.0.2](https://kaos.sh/pachca/0.0.2)

- Better helpers for properties

### [0.0.1](https://kaos.sh/pachca/0.0.1)

_The very first version_
