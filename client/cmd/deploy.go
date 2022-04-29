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

	color "github.com/TwiN/go-color"
	wapi "github.com/iamacarpet/go-win64api"
	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"
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

func logger(message string) {
	f, err := os.OpenFile("crypter-laps.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "INFO ", log.LstdFlags)
	logger.Println(message)
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
		// Main Laps Methods
		fmt.Println("Executing 'crypter deploy api' command")
		logger("Starting Crypter")

		// Set config file
		viper.SetConfigName("config")
		viper.AddConfigPath("./config")
		logger("Config File Set to ./config")

		// Find and read the config file
		logger("Loading Variables from config file")
		username := getConfigStr("username")
		description := getConfigStr("description")
		environment := getConfigStr("environment")
		minSpecialChar := getConfigInt("minSpecialChar")
		minNum := getConfigInt("minNum")
		minUpperCase := getConfigInt("minUpperCase")
		passwordLength := getConfigInt("passwordLength")

		// Generate Keys
		logger("Generating Keys")
		rand.Seed(time.Now().Unix())
		password := generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase)

		// Identify Hostname
		hostname, err := os.Hostname()
		if err == nil {
			println(color.Colorize(color.Yellow, "HostName is "+hostname))
			logger("HostName is " + hostname)
		}

		// Identify Serial
		command := exec.Command("powershell.exe", "(Get-WmiObject -class win32_bios).SerialNumber")
		logger("Running Powershell and Collecting Serial Number")
		println(color.Colorize(color.Yellow, "Running Powershell and Collecting Serial Number"))
		serialnumber, err := command.CombinedOutput()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
			logger("cmd.Run() failed, DEBUG ME")
			println(color.Colorize(color.Red, "cmd.Run() failed, DEBUG ME"))
		}
		logger("Serial Number is " + string(serialnumber))
		println(color.Colorize(color.Yellow, "Serial Number is "+string(serialnumber)))

		// Create User
		wapi.UserAdd(username, description, password)
		logger("Creating User " + username)
		println(color.Colorize(color.Yellow, "Creating User "+username))

		// Modify User Permissions
		wapi.SetAdmin(username)
		logger("Elevating  " + username + " to Administrator")
		println(color.Colorize(color.Yellow, "Elevating  "+username+" to Administrator"))

		// Send to crypt
		logger("Sending " + hostname + " " + username + " " + string(serialnumber) + " to Crypt")
		sendtocrypt(password, username, hostname, string(serialnumber), environment)
		println(color.Colorize(color.Green, "Crypter Run Completed"))

		logger("Crypter Run Completed")
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
