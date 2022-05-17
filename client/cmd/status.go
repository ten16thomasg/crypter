package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.AddCommand(statuslockerCmd)
	statusCmd.AddCommand(statusLAPSCmd)
	statusCmd.AddCommand(statusEncryptCmd)
}

var registryKey string = "Software\\Crypter"

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check status of deployed binaries (locker, laps or encrypt)",
	Long:  `This command can be used together with locker, laps or encrypt sub-commands to check status of respective deployments`,
}

var statuslockerCmd = &cobra.Command{
	Use:   "locker",
	Short: "Check status of locker deployment",
	Long:  `This command can be used to status of locker deployment`,
	Run: func(cmd *cobra.Command, args []string) {
		// *** add code to invoke automation end points below ***
		fmt.Println("Executing 'crypter status locker' placeholder command")
	},
}

var statusLAPSCmd = &cobra.Command{
	Use:   "laps",
	Short: "Check status of LAPS deployment",
	Long:  `This command can be used to status of LAPS deployment`,
	Run: func(cmd *cobra.Command, args []string) {
		// *** add code to invoke automation end points below ***
		// fmt.Println("Executing 'crypter status laps' placeholder command")
		fmt.Print(GetRegKeyStringValue(registryKey))
	},
}

var statusEncryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Check status of encrypt deployment",
	Long:  `This command can be used to status of encrypt deployment`,
	Run: func(cmd *cobra.Command, args []string) {
		// *** add code to invoke automation end points below ***
		fmt.Println("Executing 'crypter status encrypt' placeholder command")
	},
}
