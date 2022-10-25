package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	agentCmd = &cobra.Command{
		Use:   "agent",
		Long:  "The Updatemonitor agent is responsible to retrieve Updatecli manifest from a mongo database and then run Updatecli pipeline for information update",
		Short: "The Updatemonitor agent",
	}

	agentStartCmd = &cobra.Command{
		Use:   "start",
		Short: "starts an Updatemonitor agent",
		Run: func(cmd *cobra.Command, args []string) {
			err := run("agentStart")
			if err != nil {
				logrus.Errorf("command failed")
				os.Exit(1)
			}
		},
	}
)

func init() {
	agentCmd.AddCommand(
		agentStartCmd,
	)
}
