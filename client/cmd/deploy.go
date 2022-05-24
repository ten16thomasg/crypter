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
	"syscall"
	"time"
	"unsafe"

	wapi "github.com/iamacarpet/go-win64api"
	"github.com/spf13/cobra"
	"github.com/ten16thomasg/crypter/client/cmd/util"
	"golang.org/x/sys/windows/registry"
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

func winevent(evtxType, evtxLocation, evtxMessage, evtxId string) {
	command := "EventCreate"
	args := []string{"/T", "" + evtxType, "/ID", "" + evtxId, "/L", "" + evtxLocation,
		"/SO", "Go-Crypter", "/D", "" + evtxMessage}
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func creatRegKey(key string) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, "Software", registry.QUERY_VALUE)
	if err != nil {
		fmt.Println(err)
	}
	defer k.Close()

	keyName := key

	tK, exist, err := registry.CreateKey(k, keyName, registry.CREATE_SUB_KEY)
	if err != nil {
		fmt.Println(err)
	}
	defer tK.Close()

	if exist {
		fmt.Println("key %q already exists", keyName)
	}
}

func createRegKeyValue(path string) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	if err := k.SetStringValue("RotateTime", time.Now().Format("01-02-2006 15:04:05")); err != nil {
		log.Fatal(err)
	}
	if err := k.Close(); err != nil {
		log.Fatal(err)
	}
}

func isDomainJoined() bool {
	var domain *uint16
	var status uint32
	err := syscall.NetGetJoinInformation(nil, &domain, &status)
	if err != nil {
		return false
	}
	syscall.NetApiBufferFree((*byte)(unsafe.Pointer(domain)))
	return status == syscall.NetSetupDomainName
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

		fmt.Println("Executing 'crypter deploy laps' command")
		winevent("INFORMATION", "APPLICATION", "Executing 'crypter deploy laps' command", "359")

		// Set config file
		config, err := util.LoadConfig(".")
		if err != nil {
			fmt.Println("cannot load config")
			log.Fatal("cannot load config:", err)
			winevent("ERROR", "APPLICATION", "cannot load config:", "359")
		}
		winevent("INFORMATION", "APPLICATION", "Loading Configs", "359")

		// Assigning Config Values
		username := config.Account
		description := config.Description
		environment := config.Environment
		minSpecialChar := config.MinSpecialChar
		minNum := config.MinNum
		minUpperCase := config.MinUppercase
		passwordLength := config.PasswordLength
		RegKeyPath := config.RegKeyPath
		red := "\033[31m"
		green := "\033[32m"
		yellow := "\033[33m"

		logger("Loading Variables from config file", yellow)

		// Generate Keys
		winevent("INFORMATION", "APPLICATION", "Generating Keys", "359")
		logger("Generating Key", yellow)
		rand.Seed(time.Now().Unix())
		password := generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase)

		// Identify Hostname
		hostname, err := os.Hostname()
		if err == nil {
			logger("HostName is "+hostname, yellow)
			winevent("INFORMATION", "APPLICATION", "Hostname Identified as "+hostname, "359")
		}

		// Identify Serial
		command := exec.Command("powershell.exe", "(Get-WmiObject -class win32_bios).SerialNumber") //Get From the Registry!!!
		logger("Running Powershell and Collecting Serial Number", yellow)
		serialnumber, err := command.CombinedOutput()
		if err != nil {
			logger("cmd.Run() failed, DEBUG ME", red)
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		winevent("INFORMATION", "APPLICATION", "Serial Number is "+string(serialnumber), "359")
		logger("Serial Number is "+string(serialnumber), yellow)

		// Remove User
		wapi.UserDelete(username)
		logger("Removing User if Exists "+username, yellow)
		winevent("INFORMATION", "APPLICATION", "Removing User if Exists "+username, "359")

		// Create User
		wapi.UserAdd(username, description, password)
		logger("Creating User "+username, yellow)
		winevent("INFORMATION", "APPLICATION", "Creating User "+username, "359")

		// Modify User Permissions
		wapi.SetAdmin(username)
		logger("Elevating  "+username+" to Administrator", yellow)
		winevent("INFORMATION", "APPLICATION", "Elevating  "+username+" to Administrator", "359")

		// Send to crypt
		logger("Sending "+hostname+" "+username+" "+string(serialnumber)+" to Crypt", yellow)
		sendtocrypt(password, username, hostname, string(serialnumber), environment)
		winevent("INFORMATION", "APPLICATION", "Sending "+hostname+" "+username+" "+string(serialnumber)+" to Crypt", "359")

		// Add Rotate time
		creatRegKey("Crypter")
		createRegKeyValue(RegKeyPath)
		winevent("INFORMATION", "APPLICATION", "Adding Crypter RotateTime to Registry", "359")

		// Finish Crypt Run
		winevent("INFORMATION", "APPLICATION", "Crypter Run Completed", "359")
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

		if isDomainJoined() {
			fmt.Print("Is Domain Joined")
		} else {
			fmt.Print("Is NOT Domain Joined")
		}
	},
}
