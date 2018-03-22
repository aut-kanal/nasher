package telegram

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.com/kanalbot/nasher/configuration"
	"gitlab.com/kanalbot/nasher/db"
	"gitlab.com/kanalbot/nasher/models"
	"gitlab.com/kanalbot/nasher/ui/keyboard"
	"gitlab.com/kanalbot/nasher/ui/text"

	telegramAPI "gopkg.in/telegram-bot-api.v4"
)

func decodeBinary(enc string, out interface{}) {
	b64, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		logrus.WithError(err).Error("base64 decode failed")
	}
	buf := bytes.Buffer{}
	buf.Write(b64)
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(out)
	if err != nil {
		logrus.WithError(err).Error("can't decode message")
	}
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func generateNasherMessage(msg *models.Message) telegramAPI.Chattable {
	kanalUsername := configuration.GetInstance().GetString("kanal-username")

	// Keyboard
	inlineKeyboard := keyboard.NewReactionInlineKeyboard(0, 0, 0)

	// Separate text and msgId
	var text string
	var id int
	if len(msg.Text) > 0 {
		text, id = separateTextAndMsgID(msg.Text)
	} else {
		text, id = separateTextAndMsgID(msg.Caption)
	}

	// Send messages with media attached
	if msg.FileURL != "" {
		// Download attached media
		fileBytes, err := downloadFile(msg.FileURL)
		if err != nil {
			logrus.WithError(err).Error("can't download media file")
		}
		mediaFile := telegramAPI.FileBytes{
			Bytes: fileBytes,
		}

		// Kanal's BaseChat
		baseChat := telegramAPI.BaseChat{
			ChannelUsername: kanalUsername,
		}

		// Create message
		if msg.Audio != nil {
			mediaFile.Name = msg.Audio.Title
			audio := telegramAPI.AudioConfig{
				BaseFile: telegramAPI.BaseFile{
					BaseChat:    baseChat,
					File:        mediaFile,
					UseExisting: false,
				},
			}
			audio.Caption = text
			audio.ReplyMarkup = inlineKeyboard
			if id != 0 {
				audio.ReplyToMessageID = id
			}
			return audio
		}
		if msg.Voice != nil {
			mediaFile.Name = "voice"
			voice := telegramAPI.VoiceConfig{
				BaseFile: telegramAPI.BaseFile{
					BaseChat:    baseChat,
					File:        mediaFile,
					UseExisting: false,
				},
			}
			voice.Caption = text
			voice.ReplyMarkup = inlineKeyboard
			if id != 0 {
				voice.ReplyToMessageID = id
			}
			return voice
		}
		if msg.Video != nil {
			mediaFile.Name = "video"
			video := telegramAPI.VideoConfig{
				BaseFile: telegramAPI.BaseFile{
					BaseChat:    baseChat,
					File:        mediaFile,
					UseExisting: false,
				},
			}
			video.Caption = text
			video.ReplyMarkup = inlineKeyboard
			if id != 0 {
				video.ReplyToMessageID = id
			}
			return video
		}
		if msg.Document != nil {
			mediaFile.Name = msg.Document.FileName
			document := telegramAPI.DocumentConfig{
				BaseFile: telegramAPI.BaseFile{
					BaseChat:    baseChat,
					File:        mediaFile,
					UseExisting: false,
				},
			}
			document.Caption = text
			document.ReplyMarkup = inlineKeyboard
			if id != 0 {
				document.ReplyToMessageID = id
			}
			return document
		}
		if msg.Photo != nil {
			mediaFile.Name = "photo"
			photo := telegramAPI.PhotoConfig{
				BaseFile: telegramAPI.BaseFile{
					BaseChat:    baseChat,
					File:        mediaFile,
					UseExisting: false,
				},
			}
			photo.Caption = text
			photo.ReplyMarkup = inlineKeyboard
			if id != 0 {
				photo.ReplyToMessageID = id
			}
			return photo
		}
	}

	// Message without media
	replyMsg := telegramAPI.NewMessageToChannel(kanalUsername, text)
	replyMsg.ReplyMarkup = inlineKeyboard
	if id != 0 {
		replyMsg.ReplyToMessageID = id
	}
	return replyMsg
}

func reactionRemoved(callbackQueryID string) {
	callback := telegramAPI.CallbackConfig{
		CallbackQueryID: callbackQueryID,
		Text:            text.MsgReactionRemoved,
	}
	bot.AnswerCallbackQuery(callback)
}

func reactionSet(callbackQueryID string) {
	callback := telegramAPI.CallbackConfig{
		CallbackQueryID: callbackQueryID,
		Text:            text.MsgReactionSet,
	}
	bot.AnswerCallbackQuery(callback)
}

func deleteReactionsFromDB(reactions ...interface{}) {
	for _, reaction := range reactions {
		db.GetInstance().Delete(reaction)
	}
}

func addReactionToMsg(callbackID string, userID int, chatID int64, msgID int,
	chosenReaction models.Reaction, otherReactions ...models.Reaction) {
	oldChosenReaction := chosenReaction
	var chosenCount int
	db.GetInstance().Where("user_id = ? AND message_id = ?", userID, msgID).First(oldChosenReaction).Count(&chosenCount)

	if chosenCount > 0 {
		// Send chat action
		go reactionRemoved(callbackID)

		// Update DB
		db.GetInstance().Delete(oldChosenReaction)
	} else {
		// Send chat action
		go reactionSet(callbackID)

		// Add like to DB
		db.GetInstance().Create(chosenReaction)

		// Remove others from DB
		for _, reaction := range otherReactions {
			db.GetInstance().Where("user_id = ? AND message_id = ?", userID, msgID).Delete(reaction)
		}
	}

	updateMessageReactionKeys(chatID, msgID)
}

func updateMessageReactionKeys(chatID int64, msgID int) {
	var likeCount, lolCount, facepalmCount int
	db.GetInstance().Model(&models.LikeReaction{}).Where("message_id = ?", msgID).Count(&likeCount)
	db.GetInstance().Model(&models.LolReaction{}).Where("message_id = ?", msgID).Count(&lolCount)
	db.GetInstance().Model(&models.FacepalmReaction{}).Where("message_id = ?", msgID).Count(&facepalmCount)

	keyboard := keyboard.NewReactionInlineKeyboard(likeCount, lolCount, facepalmCount)
	bot.Send(telegramAPI.NewEditMessageReplyMarkup(chatID, msgID, keyboard))
}

var (
	messageLinkPattern *regexp.Regexp
)

func separateTextAndMsgID(message string) (string, int) {
	if messageLinkPattern == nil {
		urlRegex := fmt.Sprintf(`https:\/\/t\.me\/%s\/(?P<ID>\d+)`,
			configuration.GetInstance().GetString("kanal-username")[1:])
		messageLinkPattern = regexp.MustCompile(urlRegex)
	}

	id := -1
	if submatches := messageLinkPattern.FindStringSubmatch(message); len(submatches) > 0 {
		id, _ = strconv.Atoi(submatches[1])
		message = messageLinkPattern.ReplaceAllString(message, "")
	}
	return message, id
}
