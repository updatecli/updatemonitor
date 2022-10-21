package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	agentCmd = &cobra.Command{
		Use:   "agent",
		Long:  "The updatefactory agent is responsible to retrieve Updatecli manifest from a mongo database and then run Updatecli pipeline for information update",
		Short: "The Updatefactory agent",
	}

	agentStartCmd = &cobra.Command{
		Use:   "start",
		Short: "starts an Updatefactory agent",
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
