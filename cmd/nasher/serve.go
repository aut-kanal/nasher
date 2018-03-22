package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/kanalbot/nasher/mq"
	"gitlab.com/kanalbot/nasher/telegram"
)

var (
	serveCmd = &cobra.Command{
		Use:   "start",
		Short: "Start bot",
		Run:   start,
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

func start(cmd *cobra.Command, args []string) {
	logVersion()

	mq.InitMessageQueue()
	defer mq.Close()

	telegram.StartBot()
}
