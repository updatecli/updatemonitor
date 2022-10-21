package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "The updatefactory server",
	}

	serverStartCmd = &cobra.Command{
		Use:   "start",
		Short: "starts an Updatefactory server",
		Run: func(cmd *cobra.Command, args []string) {
			err := run("serverStart")
			if err != nil {
				logrus.Errorf("command failed")
				os.Exit(1)
			}
		},
	}
)

func init() {
	serverCmd.AddCommand(
		serverStartCmd,
	)
}
