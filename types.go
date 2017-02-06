package tgbotapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// APIResponse is a response from the Telegram API with the result
// stored raw.
type APIResponse struct {
	Ok          bool                `json:"ok"`
	Result      json.RawMessage     `json:"result"`
	ErrorCode   int                 `json:"error_code"`
	Description string              `json:"description"`
	Parameters  *ResponseParameters `json:"parameters"`
}

// ResponseParameters are various errors that can be returned in APIResponse.
type ResponseParameters struct {
	MigrateToChatID int64 `json:"migrate_to_chat_id"` // optional
	RetryAfter      int   `json:"retry_after"`        // optional
}

// This object represents an incoming update.
// Only one of the optional parameters can be present in any given update.
type Update struct {
	UpdateID int `json:"update_id"` // The update‘s unique identifier.
	// 	Update identifiers start from a certain positive number
	// 	and increase sequentially.
	// 	This ID becomes especially handy if you’re using Webhooks,
	// 	since it allows you to ignore repeated updates or to restore
	// 	the correct update sequence, should they get out of order.
	Message            *Message            `json:"message"` // Optional. New incoming message of any kind — text, photo, sticker, etc.
	EditedMessage      *Message            `json:"edited_message"`
	ChannelPost        *Message            `json:"channel_post"`
	EditedChannelPost  *Message            `json:"edited_channel_post"`
	InlineQuery        *InlineQuery        `json:"inline_query"`         // Optional. New incoming inline query
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"` // Optional. The result of an inline query that was chosen by a user and sent to their chat partner.
	CallbackQuery      *CallbackQuery      `json:"callback_query"`       // Optional. New incoming callback query
}

// UpdatesChannel is the channel for getting updates.
type UpdatesChannel <-chan Update

// Clear discards all unprocessed incoming updates.
func (ch UpdatesChannel) Clear() {
	for len(ch) != 0 {
		<-ch
	}
}

// User is a user on Telegram.
type User struct {
	ID        int    `json:"id"`         // Unique identifier for this user or bot
	FirstName string `json:"first_name"` // User‘s or bot’s first name
	LastName  string `json:"last_name"`  // Optional. User‘s or bot’s last name
	UserName  string `json:"username"`   // Optional. User‘s or bot’s username
}

// String displays a simple text version of a user.
//
// It is normally a user's username, but falls back to a first/last
// name as available.
func (u *User) String() string {
	if u.UserName != "" {
		return u.UserName
	}

	name := u.FirstName
	if u.LastName != "" {
		name += " " + u.LastName
	}

	return name
}

// GroupChat is a group chat.
type GroupChat struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// This object represents a chat.
type Chat struct {
	ID                  int64  `json:"id"`                             // Unique identifier for this chat, not exceeding 1e13 by absolute value
	Type                string `json:"type"`                           // Type of chat, can be either “private”, “group”, “supergroup” or “channel”
	Title               string `json:"title"`                          // Optional. Title, for channels and group chats
	UserName            string `json:"username"`                       // Optional. Username, for private chats and channels if available
	FirstName           string `json:"first_name"`                     // Optional. First name of the other party in a private chat
	LastName            string `json:"last_name"`                      // Optional. Last name of the other party in a private chat
	AllMembersAreAdmins bool   `json:"all_members_are_administrators"` // optional
}

// IsPrivate returns if the Chat is a private conversation.
func (c Chat) IsPrivate() bool {
	return c.Type == "private"
}

// IsGroup returns if the Chat is a group.
func (c Chat) IsGroup() bool {
	return c.Type == "group"
}

// IsSuperGroup returns if the Chat is a supergroup.
func (c Chat) IsSuperGroup() bool {
	return c.Type == "supergroup"
}

// IsChannel returns if the Chat is a channel.
func (c Chat) IsChannel() bool {
	return c.Type == "channel"
}

// ChatConfig returns a ChatConfig struct for chat related methods.
func (c Chat) ChatConfig() ChatConfig {
	return ChatConfig{ChatID: c.ID}
}

// Message is returned by almost every request, and contains data about
// almost anything.
type Message struct {
	MessageID            int      `json:"message_id"`              // Unique message identifier
	From                 *User    `json:"from"`                    // Optional. Sender, can be empty for messages sent to channels
	Date                 int      `json:"date"`                    // Date the message was sent in Unix time
	Chat                 *Chat    `json:"chat"`                    // Conversation the message belongs to
	ForwardFrom          *User    `json:"forward_from"`            // Optional. For forwarded messages, sender of the original message
	ForwardFromChat      *Chat    `json:"forward_from_chat"`       // optional
	ForwardFromMessageID int      `json:"forward_from_message_id"` // optional
	ForwardDate          int      `json:"forward_date"`            // Optional. For forwarded messages, date the original message was sent in Unix time
	ReplyToMessage       *Message `json:"reply_to_message"`        // Optional. For replies, the original message.
	// 	Note that the Message object in this field
	// 	will not contain further reply_to_message fields
	// 	even if it itself is a reply.
	EditDate              int              `json:"edit_date"`               // optional
	Text                  string           `json:"text"`                    // Optional. For text messages, the actual UTF-8 text of the message, 0-4096 characters.
	Entities              *[]MessageEntity `json:"entities"`                // Optional. For text messages, special entities like usernames, URLs, bot commands, etc. that appear in the text
	Audio                 *Audio           `json:"audio"`                   // Optional. Message is an audio file, information about the file
	Document              *Document        `json:"document"`                // Optional. Message is a general file, information about the file
	Game                  *Game            `json:"game"`                    // optional
	Photo                 *[]PhotoSize     `json:"photo"`                   // Optional. Message is a photo, available sizes of the photo
	Sticker               *Sticker         `json:"sticker"`                 // Optional. Message is a sticker, information about the sticker
	Video                 *Video           `json:"video"`                   // Optional. Message is a video, information about the video
	Voice                 *Voice           `json:"voice"`                   // Optional. Message is a voice message, information about the file
	Caption               string           `json:"caption"`                 // Optional. Caption for the document, photo or video, 0-200 characters
	Contact               *Contact         `json:"contact"`                 // Optional. Message is a shared contact, information about the contact
	Location              *Location        `json:"location"`                // Optional. Message is a shared location, information about the location
	Venue                 *Venue           `json:"venue"`                   // Optional. Message is a venue, information about the venue
	NewChatMember         *User            `json:"new_chat_member"`         // Optional. A new member was added to the group, information about them (this member may be the bot itself)
	LeftChatMember        *User            `json:"left_chat_member"`        // Optional. A member was removed from the group, information about them (this member may be the bot itself)
	NewChatTitle          string           `json:"new_chat_title"`          // Optional. A chat title was changed to this value
	NewChatPhoto          *[]PhotoSize     `json:"new_chat_photo"`          // Optional. A chat photo was change to this value
	DeleteChatPhoto       bool             `json:"delete_chat_photo"`       // Optional. Service message: the chat photo was deleted
	GroupChatCreated      bool             `json:"group_chat_created"`      // Optional. Service message: the group has been created
	SuperGroupChatCreated bool             `json:"supergroup_chat_created"` // Optional. Service message: the supergroup has been created
	ChannelChatCreated    bool             `json:"channel_chat_created"`    // Optional. Service message: the channel has been created
	MigrateToChatID       int64            `json:"migrate_to_chat_id"`      // Optional. The group has been migrated to a supergroup with the specified
	// 	identifier, not exceeding 1e13 by absolute value
	MigrateFromChatID int64 `json:"migrate_from_chat_id"` // Optional. The supergroup has been migrated from a group with the specified
	// 	identifier, not exceeding 1e13 by absolute value
	PinnedMessage *Message `json:"pinned_message"` // Optional. Specified message was pinned. Note that the Message object in this
	// 	field will not contain further reply_to_message fields even if it is itself a reply.
}

// Time converts the message timestamp into a Time.
func (m *Message) Time() time.Time {
	return time.Unix(int64(m.Date), 0)
}

// IsCommand returns true if message starts with '/'.
func (m *Message) IsCommand() bool {
	return m.Text != "" && m.Text[0] == '/'
}

// Command checks if the message was a command and if it was, returns the
// command. If the Message was not a command, it returns an empty string.
//
// If the command contains the at bot syntax, it removes the bot name.
func (m *Message) Command() string {
	if !m.IsCommand() {
		return ""
	}

	command := strings.SplitN(m.Text, " ", 2)[0][1:]

	if i := strings.Index(command, "@"); i != -1 {
		command = command[:i]
	}

	return command
}

// CommandArguments checks if the message was a command and if it was,
// returns all text after the command name. If the Message was not a
// command, it returns an empty string.
func (m *Message) CommandArguments() string {
	if !m.IsCommand() {
		return ""
	}

	split := strings.SplitN(m.Text, " ", 2)
	if len(split) != 2 {
		return ""
	}

	return strings.SplitN(m.Text, " ", 2)[1]
}

// This object represents one special entity in a text message. For example, hashtags, usernames, URLs, etc.
type MessageEntity struct {
	Type string `json:"type"` //Type of the entity. One of mention (@username), hashtag, bot_command, url, email, bold (bold text),
	//	italic (italic text), code (monowidth string), pre (monowidth block), text_link (for clickable text URLs)
	Offset int    `json:"offset"` // Offset in UTF-16 code units to the start of the entity
	Length int    `json:"length"` // Length of the entity in UTF-16 code units
	URL    string `json:"url"`    // Optional. For “text_link” only, url that will be opened after user taps on the text
	User   *User  `json:"user"`   // optional
}

// ParseURL attempts to parse a URL contained within a MessageEntity.
func (entity MessageEntity) ParseURL() (*url.URL, error) {
	if entity.URL == "" {
		return nil, errors.New(ErrBadURL)
	}

	return url.Parse(entity.URL)
}

// This object represents one size of a photo or a file / sticker thumbnail.
type PhotoSize struct {
	FileID   string `json:"file_id"`   // Unique identifier for this file
	Width    int    `json:"width"`     // Photo width
	Height   int    `json:"height"`    // Photo height
	FileSize int    `json:"file_size"` // Optional. File size
}

// This object represents an audio file to be treated as music by the Telegram clients.
type Audio struct {
	FileID    string `json:"file_id"`   // Unique identifier for this file
	Duration  int    `json:"duration"`  // Duration of the audio in seconds as defined by sender
	Performer string `json:"performer"` // Optional. Performer of the audio as defined by sender or by audio tags
	Title     string `json:"title"`     // Optional. Title of the audio as defined by sender or by audio tags
	MimeType  string `json:"mime_type"` // Optional. MIME type of the file as defined by sender
	FileSize  int    `json:"file_size"` // Optional. File size
}

// This object represents a general file (as opposed to photos, voice messages and audio files).
type Document struct {
	FileID    string     `json:"file_id"`   // Unique file identifier
	Thumbnail *PhotoSize `json:"thumb"`     // Optional. Document thumbnail as defined by sender
	FileName  string     `json:"file_name"` // Optional. Original filename as defined by sender
	MimeType  string     `json:"mime_type"` // Optional. MIME type of the file as defined by sender
	FileSize  int        `json:"file_size"` // Optional. File size
}

// This object represents a sticker.
type Sticker struct {
	FileID    string     `json:"file_id"`   // Unique identifier for this file
	Width     int        `json:"width"`     // Sticker width
	Height    int        `json:"height"`    // Sticker height
	Thumbnail *PhotoSize `json:"thumb"`     // Optional. Sticker thumbnail in .webp or .jpg format
	Emoji     string     `json:"emoji"`     // optional
	FileSize  int        `json:"file_size"` // Optional. File size
}

// This object represents a video file.
type Video struct {
	FileID    string     `json:"file_id"`   // Unique identifier for this file
	Width     int        `json:"width"`     // Video width as defined by sender
	Height    int        `json:"height"`    // Video height as defined by sender
	Duration  int        `json:"duration"`  // Duration of the video in seconds as defined by sender
	Thumbnail *PhotoSize `json:"thumb"`     // Optional. Video thumbnail
	MimeType  string     `json:"mime_type"` // Optional. Mime type of a file as defined by sender
	FileSize  int        `json:"file_size"` // Optional. File size
}

// This object represents a voice note.
type Voice struct {
	FileID   string `json:"file_id"`   // Unique identifier for this file
	Duration int    `json:"duration"`  // Duration of the audio in seconds as defined by sender
	MimeType string `json:"mime_type"` // Optional. MIME type of the file as defined by sender
	FileSize int    `json:"file_size"` // Optional. File size
}

// This object represents a phone contact.
type Contact struct {
	PhoneNumber string `json:"phone_number"` // Contact's phone number
	FirstName   string `json:"first_name"`   // Contact's first name
	LastName    string `json:"last_name"`    // Optional. Contact's last name
	UserID      int    `json:"user_id"`      // Optional. Contact's user identifier in Telegram
}

// This object represents a point on the map.
type Location struct {
	Longitude float64 `json:"longitude"` // Longitude as defined by sender
	Latitude  float64 `json:"latitude"`  // Latitude as defined by sender
}

// This object represents a venue.
type Venue struct {
	Location     Location `json:"location"`      // Venue location
	Title        string   `json:"title"`         // Name of the venue
	Address      string   `json:"address"`       // Address of the venue
	FoursquareID string   `json:"foursquare_id"` // Optional. Foursquare identifier of the venue
}

// This object represent a user's profile pictures.
type UserProfilePhotos struct {
	TotalCount int           `json:"total_count"` // Total number of profile pictures the target user has
	Photos     [][]PhotoSize `json:"photos"`      // Requested profile pictures (in up to 4 sizes each)
}

// This object represents a file ready to be downloaded.
// The file can be downloaded via the link https://api.telegram.org/file/bot<token>/<file_path>.
// It is guaranteed that the link will be valid for at least 1 hour.
// When the link expires, a new one can be requested by calling getFile.
// !!!Maximum file size to download is 20 MB
type File struct {
	FileID   string `json:"file_id"`   // Unique identifier for this file
	FileSize int    `json:"file_size"` // Optional. File size, if known
	FilePath string `json:"file_path"` // Optional. File path. Use https://api.telegram.org/file/bot<token>/<file_path> to get the file.
}

// Link returns a full path to the download URL for a File.
//
// It requires the Bot Token to create the link.
func (f *File) Link(token string) string {
	return fmt.Sprintf(FileEndpoint, token, f.FilePath)
}

// This object represents a custom keyboard with reply options.
type ReplyKeyboardMarkup struct {
	Keyboard       [][]KeyboardButton `json:"keyboard"`        // Array of button rows, each represented by an Array of KeyboardButton objects
	ResizeKeyboard bool               `json:"resize_keyboard"` // Optional. Requests clients to resize the keyboard vertically for optimal fit
	// 	(e.g., make the keyboard smaller if there are just two rows of buttons).
	// 	Defaults to false, in which case the custom keyboard is always of the same height as the app's
	// 	standard keyboard.
	OneTimeKeyboard bool `json:"one_time_keyboard"` // Optional. Requests clients to hide the keyboard as soon as it's been used.
	// 	The keyboard will still be available, but clients will automatically
	// 	display the usual letter-keyboard in the chat – the user can press a special
	// 	button in the input field to see the custom keyboard again. Defaults to false.
	Selective bool `json:"selective"` // Optional. Use this parameter if you want to show the keyboard to specific users only.
	// 	Targets: 1) users that are @mentioned in the text of the Message object;
	// 	2) if the bot's message is a reply (has reply_to_message_id), sender of the original message.
}

// This object represents one button of the reply keyboard.
// For simple text buttons String can be used instead of this object to specify text of the button.
// Optional fields are mutually exclusive.
type KeyboardButton struct {
	Text            string `json:"text"`             // Text of the button. If none of the optional fields are used, it will be sent to the bot as a message when the button is pressed
	RequestContact  bool   `json:"request_contact"`  // Optional. If True, the user's phone number will be sent as a contact when the button is pressed. Available in private chats only
	RequestLocation bool   `json:"request_location"` // Optional. If True, the user's current location will be sent when the button is pressed. Available in private chats only
}

// Upon receiving a message with this object,
// Telegram clients will hide the current custom keyboard and display the default letter-keyboard.
// By default, custom keyboards are displayed until a new keyboard is sent by a bot.
// An exception is made for one-time keyboards that are hidden immediately after the user presses a button
type ReplyKeyboardHide struct {
	HideKeyboard bool `json:"hide_keyboard"` // Requests clients to hide the custom keyboard
	Selective    bool `json:"selective"`     // Optional. Use this parameter if you want to hide keyboard for specific users only.
	// 	Targets: 1) users that are @mentioned in the text of the Message object;
	// 	2) if the bot's message is a reply (has reply_to_message_id), sender of the original message.
}

// ReplyKeyboardRemove allows the Bot to hide a custom keyboard.
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}

// InlineKeyboardMarkup is a custom keyboard presented for an inline bot.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"` // Array of button rows, each represented by an Array of InlineKeyboardButton objects
}

// InlineKeyboardButton is a button within a custom keyboard for
// inline query responses.
//
// Note that some values are references as even an empty string
// will change behavior.
//
// CallbackGame, if set, MUST be first button in first row.
type InlineKeyboardButton struct {
	Text                         string        `json:"text"`
	URL                          *string       `json:"url,omitempty"`                              // optional
	CallbackData                 *string       `json:"callback_data,omitempty"`                    // optional
	SwitchInlineQuery            *string       `json:"switch_inline_query,omitempty"`              // optional
	SwitchInlineQueryCurrentChat *string       `json:"switch_inline_query_current_chat,omitempty"` // optional
	CallbackGame                 *CallbackGame `json:"callback_game,omitempty"`                    // optional
}

// CallbackQuery is data sent when a keyboard button with callback data
// is clicked.
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message"`           // optional
	InlineMessageID string   `json:"inline_message_id"` // optional
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`            // optional
	GameShortName   string   `json:"game_short_name"` // optional
}

// ForceReply allows the Bot to have users directly reply to it without
// additional interaction.
type ForceReply struct {
	ForceReply bool `json:"force_reply"` // Shows reply interface to the user, as if they manually selected the bot‘s message and tapped ’Reply'
	Selective  bool `json:"selective"`   // Optional. Use this parameter if you want to force reply from specific users only.
	// 	Targets: 1) users that are @mentioned in the text of the Message object;
	// 	2) if the bot's message is a reply (has reply_to_message_id), sender of the original message.
}

// ChatMember is information about a member in a chat.
type ChatMember struct {
	User   *User  `json:"user"`
	Status string `json:"status"`
}

// IsCreator returns if the ChatMember was the creator of the chat.
func (chat ChatMember) IsCreator() bool { return chat.Status == "creator" }

// IsAdministrator returns if the ChatMember is a chat administrator.
func (chat ChatMember) IsAdministrator() bool { return chat.Status == "administrator" }

// IsMember returns if the ChatMember is a current member of the chat.
func (chat ChatMember) IsMember() bool { return chat.Status == "member" }

// HasLeft returns if the ChatMember left the chat.
func (chat ChatMember) HasLeft() bool { return chat.Status == "left" }

// WasKicked returns if the ChatMember was kicked from the chat.
func (chat ChatMember) WasKicked() bool { return chat.Status == "kicked" }

// Game is a game within Telegram.
type Game struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Photo        []PhotoSize     `json:"photo"`
	Text         string          `json:"text"`
	TextEntities []MessageEntity `json:"text_entities"`
	Animation    Animation       `json:"animation"`
}

// Animation is a GIF animation demonstrating the game.
type Animation struct {
	FileID   string    `json:"file_id"`
	Thumb    PhotoSize `json:"thumb"`
	FileName string    `json:"file_name"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
}

// GameHighScore is a user's score and position on the leaderboard.
type GameHighScore struct {
	Position int  `json:"position"`
	User     User `json:"user"`
	Score    int  `json:"score"`
}

// CallbackGame is for starting a game in an inline keyboard button.
type CallbackGame struct{}

// WebhookInfo is information about a currently set webhook.
type WebhookInfo struct {
	URL                  string `json:"url"`
	HasCustomCertificate bool   `json:"has_custom_certificate"`
	PendingUpdateCount   int    `json:"pending_update_count"`
	LastErrorDate        int    `json:"last_error_date"`    // optional
	LastErrorMessage     string `json:"last_error_message"` // optional
}

// IsSet returns true if a webhook is currently set.
func (info WebhookInfo) IsSet() bool {
	return info.URL != ""
}

// InlineQuery is a Query from Telegram for an inline request.
type InlineQuery struct {
	ID       string    `json:"id"`       // Unique identifier for this query
	From     *User     `json:"from"`     // Sender
	Location *Location `json:"location"` // Optional. Sender location, only for bots that request user location
	Query    string    `json:"query"`    // Text of the query
	Offset   string    `json:"offset"`   // Offset of the results to be returned, can be controlled by the bot
}

// Represents a link to an article or web page.
type InlineQueryResultArticle struct {
	Type                string                `json:"type"`                            // Type of the result, must be article
	ID                  string                `json:"id"`                              // Unique identifier for this result, 1-64 Bytes
	Title               string                `json:"title"`                           // Title of the result
	InputMessageContent interface{}           `json:"input_message_content,omitempty"` // Content of the message to be sent
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	URL                 string                `json:"url"`                             // Optional. URL of the result
	HideURL             bool                  `json:"hide_url"`                        // Optional. Pass True, if you don't want the URL to be shown in the message
	Description         string                `json:"description"`                     // Optional. Short description of the result
	ThumbURL            string                `json:"thumb_url"`                       // Optional. Url of the thumbnail for the result
	ThumbWidth          int                   `json:"thumb_width"`                     // Optional. Thumbnail width
	ThumbHeight         int                   `json:"thumb_height"`                    // Optional. Thumbnail height
}

// Represents a link to a photo. By default, this photo will be sent by the user with optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the photo.
type InlineQueryResultPhoto struct {
	Type                string                `json:"type"`                            // Type of the result, must be photo
	ID                  string                `json:"id"`                              // Unique identifier for this result, 1-64 bytes
	URL                 string                `json:"photo_url"`                       // A valid URL of the photo. Photo must be in jpeg format. Photo size must not exceed 5MB
	ThumbURL            string                `json:"thumb_url"`                       // Optional. Title for the result
	Width               int                   `json:"photo_width"`                     // Optional. Width of the photo
	Height              int                   `json:"photo_height"`                    // Optional. Height of the photo
	Title               string                `json:"title"`                           // Optional. Title for the result
	Description         string                `json:"description"`                     // Optional. Short description of the result
	Caption             string                `json:"caption"`                         // Optional. Caption of the photo to be sent, 0-200 characters
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the photo
}

// Represents a link to an animated GIF file.
// By default, this animated GIF file will be sent by the user with optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.
type InlineQueryResultGIF struct {
	Type                string                `json:"type"`                            // Type of the result, must be gif
	ID                  string                `json:"id"`                              // Unique identifier for this result, 1-64 bytes
	URL                 string                `json:"gif_url"`                         // A valid URL for the GIF file. File size must not exceed 1MB
	Width               int                   `json:"gif_width"`                       // Optional. Width of the GIF
	Height              int                   `json:"gif_height"`                      // Optional. Height of the GIF
	ThumbURL            string                `json:"thumb_url"`                       // URL of the static thumbnail for the result (jpeg or gif)
	Title               string                `json:"title"`                           // Optional. Title for the result
	Caption             string                `json:"caption"`                         //  	Optional. Caption of the GIF file to be sent, 0-200 characters
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the GIF animation
}

// Represents a link to a video animation (H.264/MPEG-4 AVC video without sound).
// By default, this animated MPEG-4 file will be sent by the user with optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the animation.
type InlineQueryResultMPEG4GIF struct {
	Type                string                `json:"type"`                            // Type of the result, must be mpeg4_gif
	ID                  string                `json:"id"`                              // Unique identifier for this result, 1-64 bytes
	URL                 string                `json:"mpeg4_url"`                       // A valid URL for the MP4 file. File size must not exceed 1MB
	Width               int                   `json:"mpeg4_width"`                     // Optional. Video width
	Height              int                   `json:"mpeg4_height"`                    // Optional. Video height
	ThumbURL            string                `json:"thumb_url"`                       // URL of the static thumbnail (jpeg or gif) for the result
	Title               string                `json:"title"`                           // Optional. Title for the result
	Caption             string                `json:"caption"`                         //  	Optional. Caption of the MPEG-4 file to be sent, 0-200 characters
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the video animation
}

// Represents a link to a page containing an embedded video player or a video file.
// By default, this video file will be sent by the user with an optional caption.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the video.
type InlineQueryResultVideo struct {
	Type                string                `json:"type"`                            // Type of the result, must be video
	ID                  string                `json:"id"`                              // Unique identifier for this result, 1-64 bytes
	URL                 string                `json:"video_url"`                       // A valid URL for the embedded video player or video file
	MimeType            string                `json:"mime_type"`                       // Mime type of the content of video url, “text/html” or “video/mp4”
	ThumbURL            string                `json:"thumb_url"`                       // URL of the thumbnail (jpeg only) for the video
	Title               string                `json:"title"`                           // Title for the result
	Caption             string                `json:"caption"`                         // Optional. Caption of the video to be sent, 0-200 characters
	Width               int                   `json:"video_width"`                     // Optional. Video width
	Height              int                   `json:"video_height"`                    // Optional. Video height
	Duration            int                   `json:"video_duration"`                  // Optional. Video duration in seconds
	Description         string                `json:"description"`                     // Optional. Short description of the result
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the video
}

// Represents a link to an mp3 audio file. By default, this audio file will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the audio.
type InlineQueryResultAudio struct {
	Type                string                `json:"type"`      // required
	ID                  string                `json:"id"`        // required
	URL                 string                `json:"audio_url"` // required
	Title               string                `json:"title"`     // required
	Caption             string                `json:"caption"`
	Performer           string                `json:"performer"`
	Duration            int                   `json:"audio_duration"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultVoice is an inline query response voice.
type InlineQueryResultVoice struct {
	Type                string                `json:"type"`      // required
	ID                  string                `json:"id"`        // required
	URL                 string                `json:"voice_url"` // required
	Title               string                `json:"title"`     // required
	Caption             string                `json:"caption"`
	Duration            int                   `json:"voice_duration"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
}

// InlineQueryResultDocument is an inline query response document.
type InlineQueryResultDocument struct {
	Type                string                `json:"type"`                            // Type of the result, must be document
	ID                  string                `json:"id"`                              // Unique identifier for this result, 1-64 bytes
	Title               string                `json:"title"`                           // Title for the result
	Caption             string                `json:"caption"`                         // Optional. Caption of the document to be sent, 0-200 characters
	URL                 string                `json:"document_url"`                    // A valid URL for the file
	MimeType            string                `json:"mime_type"`                       // Mime type of the content of the file, either “application/pdf” or “application/zip”
	Description         string                `json:"description"`                     // Optional. Short description of the result
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`          // Optional. Inline keyboard attached to the message
	InputMessageContent interface{}           `json:"input_message_content,omitempty"` // Optional. Content of the message to be sent instead of the file
	ThumbURL            string                `json:"thumb_url"`                       // Optional. URL of the thumbnail (jpeg only) for the file
	ThumbWidth          int                   `json:"thumb_width"`                     // Optional. Thumbnail width
	ThumbHeight         int                   `json:"thumb_height"`                    //  	Optional. Thumbnail height
}

// Represents a location on a map. By default, the location will be sent by the user.
// Alternatively, you can use input_message_content to send a message with the specified content instead of the location.
type InlineQueryResultLocation struct {
	Type                string                `json:"type"`      // required
	ID                  string                `json:"id"`        // required
	Latitude            float64               `json:"latitude"`  // required
	Longitude           float64               `json:"longitude"` // required
	Title               string                `json:"title"`     // required
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent interface{}           `json:"input_message_content,omitempty"`
	ThumbURL            string                `json:"thumb_url"`
	ThumbWidth          int                   `json:"thumb_width"`
	ThumbHeight         int                   `json:"thumb_height"`
}

// InlineQueryResultGame is an inline query response game.
type InlineQueryResultGame struct {
	Type          string                `json:"type"`
	ID            string                `json:"id"`
	GameShortName string                `json:"game_short_name"`
	ReplyMarkup   *InlineKeyboardMarkup `json:"reply_markup"`
}

// Represents a result of an inline query that was chosen by the user and sent to their chat partner.
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`         // The unique identifier for the result that was chosen
	From            *User     `json:"from"`              // The user that chose the result
	Location        *Location `json:"location"`          // Optional. Sender location, only for bots that require user location
	InlineMessageID string    `json:"inline_message_id"` // Optional. Identifier of the sent inline message. Available only if there is an inline keyboard attached to the message. Will be also received in callback queries and can be used to edit the message.
	Query           string    `json:"query"`             // The query that was used to obtain the result
}

// InputTextMessageContent contains text for displaying
// as an inline query result.
type InputTextMessageContent struct {
	Text                  string `json:"message_text"`             // Text of the message to be sent, 1-4096 characters
	ParseMode             string `json:"parse_mode"`               // Optional. Send Markdown or HTML, if you want Telegram apps to show bold, italic, fixed-width text or inline URLs in your bot's message.
	DisableWebPagePreview bool   `json:"disable_web_page_preview"` // Optional. Disables link previews for links in the sent message
}

// Represents the content of a location message to be sent as the result of an inline query.
type InputLocationMessageContent struct {
	Latitude  float64 `json:"latitude"`  // Latitude of the location in degrees
	Longitude float64 `json:"longitude"` // Longitude of the location in degrees
}

// Represents the content of a venue message to be sent as the result of an inline query.
type InputVenueMessageContent struct {
	Latitude     float64 `json:"latitude"`      // Latitude of the venue in degrees
	Longitude    float64 `json:"longitude"`     // Longitude of the venue in degrees
	Title        string  `json:"title"`         // Name of the venue
	Address      string  `json:"address"`       // Address of the venue
	FoursquareID string  `json:"foursquare_id"` // Optional. Foursquare identifier of the venue, if known
}

// Represents the content of a contact message to be sent as the result of an inline query.
type InputContactMessageContent struct {
	PhoneNumber string `json:"phone_number"` // Contact's phone number
	FirstName   string `json:"first_name"`   //  	Contact's first name
	LastName    string `json:"last_name"`    // Optional. Contact's last name
}
