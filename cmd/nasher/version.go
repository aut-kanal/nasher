package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/kanalbot/nasher"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "nasher's version",
	Run: func(cmd *cobra.Command, args []string) {
		logVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func logVersion() {
	logrus.Info("version   > ", nasher.Version)
	logrus.Info("buildtime > ", nasher.BuildTime)
	logrus.Info("commit    > ", nasher.Commit)
}
