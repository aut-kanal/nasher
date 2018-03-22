package models

import telegramAPI "gopkg.in/telegram-bot-api.v4"

type Message struct {
	*telegramAPI.Message
	FileURL string
}
