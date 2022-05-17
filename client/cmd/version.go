package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print version number of crypter tool",
	Aliases: []string{"ver", "v"},
	Long:    `This command can be used get the version number of crypter tool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("crypter-v3.0.0-alpha")
	},
}
