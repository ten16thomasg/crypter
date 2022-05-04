package cmd

import (
	"fmt"
	"log"

	wapi "github.com/iamacarpet/go-win64api"
	"github.com/spf13/cobra"
	"github.com/ten16thomasg/crypter/client/cmd/util"
)

func init() {
	rootCmd.AddCommand(undeployCmd)

	undeployCmd.AddCommand(undeployLockerCmd)
	undeployCmd.AddCommand(undeployLAPSCmd)
	undeployCmd.AddCommand(undeployEncryptCmd)

}

var undeployCmd = &cobra.Command{
	Use:     "undeploy",
	Aliases: []string{"undep", "undepl"},
	Short:   "Undeploy artifacts (locker, laps or encrypt)",
	Long:    `This command can be used together with locker, laps or encrypt sub-commands to undeploy respective artifacts`,
}

var undeployLockerCmd = &cobra.Command{
	Use:   "locker",
	Short: "Undeploy locker artifacts",
	Long:  `This command can be used to undeploy locker artifacts`,
	Run: func(cmd *cobra.Command, args []string) {
		// *** add code to invoke automation end points below ***
		fmt.Println("Executing 'crypter undeploy locker' ")
	},
}

var undeployLAPSCmd = &cobra.Command{
	Use:   "laps",
	Short: "Uneploy LAPS artifacts",
	Long:  `This command can be used to undeploy LAPS artifacts`,
	Run: func(cmd *cobra.Command, args []string) {
		// Format Logs
		red := "\033[31m"
		green := "\033[32m"
		yellow := "\033[33m"

		// Start crypter Deploy
		logger("Executing 'crypter deploy laps' command", yellow)
		winevent("INFORMATION", "APPLICATION", "Executing 'crypter deploy laps' command", "359")

		// Set config file
		logger("Loading Configuration File", yellow)
		config, err := util.LoadConfig(".")
		if err != nil {
			logger("cannot load config", red)
			log.Fatal("cannot load config:", err)
		}

		// Assigning Config Values
		logger("Loading Variables from config file", yellow)
		username := config.Account

		// Revoke Admin
		wapi.RevokeAdmin(username)
		logger("Revoking User "+username+" Permissions ", yellow)

		// UserDisabled adds or removes the flag that disables a user's account
		wapi.UserDisabled(username, true)
		logger("Disabling  "+username, yellow)

		// Return if Undeploy Failed
		logger("Checking if "+username+" Is Admin", yellow)
		t, err := wapi.IsLocalUserAdmin(username)
		if err != nil {
			fmt.Println(err)
			logger("Undeploy Failed", red)
		}
		fmt.Println(t)
		logger("Crypter Undeploy Completed ", green)
	},
}

var undeployEncryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Undeploy encrypt artifacts",
	Long:  `This command can be used to undeploy encrypt artifacts`,
	Run: func(cmd *cobra.Command, args []string) {
		// *** add code to invoke automation end points below ***
		fmt.Println("Executing 'crypter undeploy encrypt' ")
	},
}
