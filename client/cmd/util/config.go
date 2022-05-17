package util

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Description    string `mapstructure:"DESCRIPTION"`
	Account        string `mapstructure:"ACCOUNT"`
	Environment    string `mapstructure:"ENVIRONMENT"`
	LogRed         string `mapstructure:"LOGRED"`
	LogGreen       string `mapstructure:"LOGGREEN"`
	LogYellow      string `mapstructure:"LOGYELLOW"`
	MinSpecialChar int    `mapstructure:"MINSPECIALCHAR"`
	MinNum         int    `mapstructure:"MINNUM"`
	MinUppercase   int    `mapstructure:"MINUPPERCASE"`
	PasswordLength int    `mapstructure:"PASSWORDLENGTH"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
