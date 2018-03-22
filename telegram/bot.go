package telegram

import (
	"fmt"

	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"gitlab.com/kanalbot/nasher/configuration"
	"gitlab.com/kanalbot/nasher/mq"
	"gitlab.com/kanalbot/nasher/ui/keyboard"
)

var (
	bot *miyanbor.Bot
)

func StartBot() {
	botDebug := configuration.GetInstance().GetBool("bot.telegram.debug")
	botToken := configuration.GetInstance().GetString("bot.telegram.token")
	botSessionTimeout := configuration.GetInstance().GetInt("bot.telegram.session-timeout")
	botUpdaterTimeout := configuration.GetInstance().GetInt("bot.telegram.updater-timeout")

	var err error
	bot, err = miyanbor.NewBot(botToken, botDebug, botSessionTimeout)
	if err != nil {
		logrus.WithError(err).Fatalf("can't init bot")
	}
	logrus.Infof("telegram bot initialized completely")

	mq.SubscribeAcceptedMsgs(newAcceptedMessageHandler)
	logrus.Info("subscribed on msgs queue")
	logrus.Infof("===================================")

	setCallbacks(bot)
	bot.StartUpdater(0, botUpdaterTimeout)
}

func setCallbacks(bot *miyanbor.Bot) {
	bot.SetSessionStartCallbackHandler(sessionStartHandler)
	bot.SetFallbackCallbackHandler(unknownMessageHandler)

	bot.AddCallbackHandler(fmt.Sprintf("^%s$", keyboard.KeyboardLikeButtonData), likeHandler)
	bot.AddCallbackHandler(fmt.Sprintf("^%s$", keyboard.KeyboardLolButtonData), lolHandler)
	bot.AddCallbackHandler(fmt.Sprintf("^%s$", keyboard.KeyboardFacepalmButtonData), facepalmHandler)
}
