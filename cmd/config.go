package cmd

import (
	"github.com/imander/bootcampspot/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "configure",
	Aliases: []string{"config", "conf"},
	Short:   "set bootcamp url and login details",
	Run: func(cmd *cobra.Command, args []string) {
		config.Set()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
