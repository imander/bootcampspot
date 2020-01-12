package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/imander/bootcampspot/util"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// Configuration stores all configuration information for the bcs CLI
type Configuration struct {
	URL      string
	User     string
	Password string
}

const defaultURL = "https://bootcampspot.com/api/instructor/v1"

// BCS config is made available to other packages
var BCS = &Configuration{}

// Load loads the bcs configurtion into the runtime
func Load(prompt bool) {
	viper.SetConfigFile(configFile())

	if err := viper.ReadInConfig(); err != nil {
		if strings.Contains(err.Error(), "no such file") {
			if prompt {
				promptConfig()
				os.Exit(0)
			}
			return
		}
		log.Fatal(err)
	}

	setConfig()
}

// Set prompts for user input to build the configuration file for the bcs CLI
func Set() {
	setValue("url", "Enter API endpoint", false)
	setValue("user", "Enter user name", false)
	setValue("password", "Enter password", true)

	if err := viper.WriteConfigAs(configFile()); err != nil {
		log.Fatal(err)
	}
}

func getInput(hide bool) (str string) {
	if hide {
		val, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Printf("error reading std-in: %s", err.Error())
			os.Exit(1)
		}
		fmt.Println()
		return string(val)
	}

	util.ReadInput("", &str)
	return
}

func setValue(name, promp string, hide bool) {
	curr := viper.Get(name)
	if curr == nil {
		if name == "url" {
			viper.Set(name, string(defaultURL))
			curr = defaultURL
		} else {
			curr = ""
		}
	}
	if hide {
		curr = "*****"
	}

	fmt.Printf("%s [%s]: ", promp, curr)
	val := getInput(hide)
	if len(val) != 0 {
		viper.Set(name, string(val))
	}
}

func setConfig() {
	BCS = &Configuration{
		URL:      viper.GetString("url"),
		User:     viper.GetString("user"),
		Password: viper.GetString("password"),
	}
}

func promptConfig() {
	fmt.Printf("Config not found. Generate now? [Y/n]: ")
	var ans string
	fmt.Scanf("%s", &ans)
	if strings.ToLower(ans) != "n" {
		Set()
	}
}

func configFile() string {
	return os.Getenv("HOME") + "/.bcs.yaml"
}
