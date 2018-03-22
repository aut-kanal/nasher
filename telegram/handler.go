package telegram

import (
	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/kanalbot/nasher/models"
	telegramAPI "gopkg.in/telegram-bot-api.v4"
)

func sessionStartHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.WithField("user", *userSession).Debugf("new session started")
}

func unknownMessageHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.WithField("user", *userSession).Debugf("unknown message received")
}

func newAcceptedMessageHandler(msg amqp.Delivery) {
	logrus.Debug("new message arrived")

	// Decode message
	decodedMsg := &models.Message{}
	decodeBinary(string(msg.Body), decodedMsg)

	// Send to nasher group
	bot.Send(generateNasherMessage(decodedMsg))
}

func likeHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.Debug("like handler")

	messageID := update.(*telegramAPI.Update).CallbackQuery.Message.MessageID
	callbackID := update.(*telegramAPI.Update).CallbackQuery.ID

	likeReaction := &models.LikeReaction{
		MessageID: messageID,
		UserID:    userSession.UserID,
	}

	addReactionToMsg(callbackID, userSession.UserID, userSession.ChatID, messageID,
		likeReaction, &models.LolReaction{}, &models.FacepalmReaction{})
}

func lolHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.Debug("lol handler")

	messageID := update.(*telegramAPI.Update).CallbackQuery.Message.MessageID
	callbackID := update.(*telegramAPI.Update).CallbackQuery.ID

	lolReaction := &models.LolReaction{
		MessageID: messageID,
		UserID:    userSession.UserID,
	}

	addReactionToMsg(callbackID, userSession.UserID, userSession.ChatID, messageID,
		lolReaction, &models.LikeReaction{}, &models.FacepalmReaction{})
}

func facepalmHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.Debug("facepalm handler")

	messageID := update.(*telegramAPI.Update).CallbackQuery.Message.MessageID
	callbackID := update.(*telegramAPI.Update).CallbackQuery.ID

	facePalmReaction := &models.FacepalmReaction{
		MessageID: messageID,
		UserID:    userSession.UserID,
	}

	addReactionToMsg(callbackID, userSession.UserID, userSession.ChatID, messageID,
		facePalmReaction, &models.LolReaction{}, &models.LikeReaction{})
}
