package telegram

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/kanalbot/nasher/configuration"
	"gitlab.com/kanalbot/nasher/models"
	"gitlab.com/kanalbot/nasher/ui/keyboard"

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
