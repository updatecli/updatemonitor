package cmd

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/updatecli/updateserver/pkg/engine"
)

var (
	// Server configuration file
	cfgFile string

	// Verbose allows to enable/disable debug logging
	verbose bool
	rootCmd = &cobra.Command{
		Use:   "updatefactory",
		Short: "updatefactory is the server alternative to Updatecli",
		Long:  `A long running Updatecli pipeline`,
		Run: func(cmd *cobra.Command, args []string) {

			var o engine.Options

			if err := viper.Unmarshal(&o); err != nil {
				logrus.Errorln(err)
				os.Exit(1)
			}

			e := engine.Engine{
				Options: o,
			}

			e.Start()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			logrus.Infoln("See you next time")
		},
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("%s", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "set config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "debug", "", false, "set log level")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}
	}

	viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug"))

	rootCmd.AddCommand(
		versionCmd,
	)

}

func initConfig() {

	viper.SetConfigName("config") // name of config file (without extension)
	if cfgFile != "" {
		viper.SetConfigName(cfgFile)
	}

	viper.SetConfigType("yaml")                 // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/updatefactory/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.updatefactory") // call multiple times to add many search paths
	viper.AddConfigPath(".")                    // optionally look for config in the working directory
	viper.ReadInConfig()                        // Find and read the config file

	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Infof("Config file changed:", e.Name)
	})
	viper.WatchConfig()

}
