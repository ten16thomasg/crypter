package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	wapi "github.com/iamacarpet/go-win64api"
	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"
	"github.com/ten16thomasg/crypter/client/cmd/util"
)

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.AddCommand(deployLockerCmd)
	deployCmd.AddCommand(deployLAPSCmd)
	deployCmd.AddCommand(deployEncryptCmd)

}

func generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
	var (
		lowerCharSet   = "abcdedfghijklmnopqrst"
		upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		specialCharSet = "!@#$%&"
		numberSet      = "0123456789"
		allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
	)
	var password strings.Builder

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func getConfigStr(key string) string {

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(key).(string)

	// If the type is a string then ok will be true
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

func getConfigInt(key string) int {

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(key).(int)

	// If the type is a string then ok will be true
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

func sendtocrypt(pass, username, hostname, serialnumber, environment string) {

	data := url.Values{
		"recovery_password": {pass},
		"serial":            {serialnumber},
		"username":          {username},
		"macname":           {hostname},
	}
	var (
		env = environment
	)
	resp, err := http.PostForm(env, data)

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res["form"])
}

func logger(message, clr string) {
	f, err := os.OpenFile("crypter-laps.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "INFO ", log.LstdFlags)
	logger.Println(message)
	fmt.Println(string(clr), message)
}

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"dep", "depl"},
	Short:   "Deploy artifacts (laps, locker or encrypt)",
	Long:    `This command can be used together with locker, laps or encrypt sub-commands to interact with crypt`,
}

var deployLockerCmd = &cobra.Command{
	Use:   "locker",
	Short: "Deploy locker artifacts",
	Long:  `This command can be used to deploy Locker`,
	Run: func(cmd *cobra.Command, args []string) {
		// *** add code to invoke locker ***
		fmt.Println("Executing 'crypter deploy Locker' command")
	},
}

var deployLAPSCmd = &cobra.Command{
	Use:   "laps",
	Short: "Deploy LAPS artifacts",
	Long:  `This command can be used to deploy LAPS artifacts`,
	Run: func(cmd *cobra.Command, args []string) {
		// Format Logs
		red := "\033[31m"
		green := "\033[32m"
		yellow := "\033[33m"
		logger("Executing 'crypter deploy laps' command", yellow)

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
		description := config.Description
		environment := config.Environment
		minSpecialChar := config.MinSpecialChar
		minNum := config.MinNum
		minUpperCase := config.MinUppercase
		passwordLength := config.PasswordLength

		// Generate Keys
		logger("Generating Keys", yellow)
		rand.Seed(time.Now().Unix())
		password := generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase)

		// Identify Hostname
		hostname, err := os.Hostname()
		if err == nil {
			logger("HostName is "+hostname, yellow)
		}

		// Identify Serial
		command := exec.Command("powershell.exe", "(Get-WmiObject -class win32_bios).SerialNumber")
		logger("Running Powershell and Collecting Serial Number", yellow)
		serialnumber, err := command.CombinedOutput()
		if err != nil {
			logger("cmd.Run() failed, DEBUG ME", "Red")
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		logger("Serial Number is "+string(serialnumber), yellow)

		// Create User
		wapi.UserAdd(username, description, password)
		logger("Creating User "+username, yellow)

		// Modify User Permissions
		wapi.SetAdmin(username)
		logger("Elevating  "+username+" to Administrator", yellow)

		// Send to crypt
		logger("Sending "+hostname+" "+username+" "+string(serialnumber)+" to Crypt", yellow)
		sendtocrypt(password, username, hostname, string(serialnumber), environment)

		logger("Crypter Run Completed", green)
	},
}

var deployEncryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Deploy encrypt artifacts",
	Long:  `This command can be used to deploy encrypt artifacts`,
	Run: func(cmd *cobra.Command, args []string) {
		// *** add code to encrypt ***
		fmt.Println("Executing 'crypter deploy encrypt'  command")
	},
}
