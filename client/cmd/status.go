package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/ten16thomasg/crypter/client/cmd/util"
	"golang.org/x/sys/windows/registry"
)

type State struct {
	RotateTime     string
	CrypterEnabled bool
}

func GetRegKeyStringValue(registryKey string) string {
	var access uint32 = registry.QUERY_VALUE
	regKey, err := registry.OpenKey(registry.LOCAL_MACHINE, registryKey, access)
	if err != nil {
		if err != registry.ErrNotExist {
			panic(err)
		}
		fmt.Println(false)
	}

	id, _, err := regKey.GetStringValue("RotateTime")
	if err != nil {
		panic(err)
		fmt.Println(false)
	}
	data := State{
		RotateTime:     id,
		CrypterEnabled: true,
	}
	file, _ := json.MarshalIndent(data, "", " ")
	return string(file)

}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.AddCommand(statuslockerCmd)
	statusCmd.AddCommand(statusLAPSCmd)
	statusCmd.AddCommand(statusEncryptCmd)
}

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
		// Set config file
		config, err := util.LoadConfig(".")
		if err != nil {
			fmt.Println("cannot load config")
			log.Fatal("cannot load config:", err)
			winevent("ERROR", "APPLICATION", "cannot load config:", "359")
		}
		winevent("INFORMATION", "APPLICATION", "Loading Configs", "359")

		// Assigning Config Values
		RegKeyPath := config.RegKeyPath

		fmt.Print(GetRegKeyStringValue(RegKeyPath))
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
