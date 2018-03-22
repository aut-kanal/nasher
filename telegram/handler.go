package telegram

import (
	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/kanalbot/nasher/models"
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
}

func lolHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
}

func facepalmHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
}
