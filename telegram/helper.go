package telegram

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"io/ioutil"
	"net/http"

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

	inlineKeyboard := keyboard.NewReactionInlineKeyboard(0, 0, 0)

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
			audio.Caption = msg.Caption
			audio.ReplyMarkup = inlineKeyboard
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
			voice.Caption = msg.Caption
			voice.ReplyMarkup = inlineKeyboard
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
			video.Caption = msg.Caption
			video.ReplyMarkup = inlineKeyboard
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
			document.Caption = msg.Caption
			document.ReplyMarkup = inlineKeyboard
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
			photo.Caption = msg.Caption
			photo.ReplyMarkup = inlineKeyboard
			return photo
		}
	}

	// Message without media
	replyMsg := telegramAPI.NewMessageToChannel(kanalUsername, msg.Text)
	replyMsg.ReplyMarkup = inlineKeyboard
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
		db.GetInstance().Delete(chosenReaction)
	} else {
		// Send chat action
		go reactionSet(callbackID)

		// Add like to DB
		db.GetInstance().Create(chosenReaction)

		// Remove others from DB
		for _, reaction := range otherReactions {
			db.GetInstance().Delete(reaction)
		}
	}

	updateMessageReactionKeys(chatID, msgID)
}

func updateMessageReactionKeys(chatID int64, msgID int) {
	var likeCount, lolCount, facepalmCount int
	db.GetInstance().Model(&models.LikeReaction{}).Where("message_id = ?", msgID).Count(&likeCount)
	db.GetInstance().Model(&models.LolReaction{}).Where("message_id = ?", msgID).Count(&lolCount)
	db.GetInstance().Model(&models.FacepalmReaction{}).Where("message_id = ?", msgID).Count(&facepalmCount)

	logrus.Errorln(msgID, ":", likeCount, lolCount, facepalmCount)
	keyboard := keyboard.NewReactionInlineKeyboard(likeCount, lolCount, facepalmCount)
	bot.Send(telegramAPI.NewEditMessageReplyMarkup(chatID, msgID, keyboard))
}
