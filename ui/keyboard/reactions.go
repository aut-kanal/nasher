package keyboard

import (
	"fmt"

	"gitlab.com/kanalbot/nasher/ui/text"
	telegramAPI "gopkg.in/telegram-bot-api.v4"
)

const (
	KeyboardLikeButtonData     = "li"
	KeyboardLolButtonData      = "lo"
	KeyboardFacepalmButtonData = "fp"
)

func NewReactionInlineKeyboard(likeCount, lolCount, facepalmCount int) telegramAPI.InlineKeyboardMarkup {
	var row []telegramAPI.InlineKeyboardButton

	like := fmt.Sprintf(text.KeyboardEmojiLike, likeCount)
	likeButton := telegramAPI.NewInlineKeyboardButtonData(like, KeyboardLikeButtonData)
	lol := fmt.Sprintf(text.KeyboardEmojiLol, facepalmCount)
	lolButton := telegramAPI.NewInlineKeyboardButtonData(lol, KeyboardLolButtonData)
	facepalm := fmt.Sprintf(text.KeyboardEmojiFacepalm, facepalmCount)
	facepalmButton := telegramAPI.NewInlineKeyboardButtonData(facepalm, KeyboardFacepalmButtonData)

	row = append(row, likeButton, lolButton, facepalmButton)
	return telegramAPI.NewInlineKeyboardMarkup(row)
}
