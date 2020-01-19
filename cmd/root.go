package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bcs",
	Short: "Obtain student details from the bootcampspot API",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitError(err)
	}
}

func exitError(err error) {
	fmt.Printf("error: %v\n", err)
	os.Exit(1)
}
