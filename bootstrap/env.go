package bootstrap

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Env struct {
	AppEnv string `mapstructure:"APP_ENV"`

	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`

	MigrationURL string `mapstructure:"MIGRATION_URL"`

	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	//ContextTimeout        int    `mapstructure:"CONTEXT_TIMEOUT"`
	TokenExpiryHour        int `mapstructure:"TOKEN_EXPIRY_HOUR"`
	MaxFailedLoginAttempts int `mapstructure:"MAX_FAILED_LOGIN_ATTEMPTS"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("can't find the file .env: %v", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		logrus.Fatal("environment can't be loaded: %v", err)
	}

	if env.AppEnv == "development" {
		logrus.Println("the App is running in development env")
	}

	return &env
}
