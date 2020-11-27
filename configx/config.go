package configx

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
	"time"
)

var once sync.Once

type LogConfig struct {
	Level string `json:"level"`
}

type AppConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Port        int    `json:"port"`
	Environment string `json:"environment"`
}

type Configuration struct {
	App AppConfig `json:"app"`
	Log LogConfig `json:"log"`
}

var (
	instance *Configuration
)

func GetConfig() *Configuration {

	once.Do(func() {
		v := viper.New()

		// Viper settings
		v.SetConfigName("config")
		v.AddConfigPath(".")
		v.AddConfigPath("$CONFIG_DIR/")

		// Environment variable settings
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
		v.AllowEmptyEnv(true)
		v.AutomaticEnv()

		// Global configuration
		v.SetDefault("shutdownTimeout", 15*time.Second)
		if _, ok := os.LookupEnv("NO_COLOR"); ok {
			v.SetDefault("no_color", true)
		}

		// Database configuration
		_ = v.BindEnv("db.host")
		v.SetDefault("db.port", 5432)
		_ = v.BindEnv("db.user")
		_ = v.BindEnv("db.password")
		_ = v.BindEnv("db.database")

		err := v.ReadInConfig()
		if err != nil {
			panic(errors.Wrap(err, "Cant read configuration file"))
		}

		instance = &Configuration{}

		err = v.Unmarshal(&instance)
		if err != nil {
			panic(errors.Wrap(err, "Cant unmarshall configuration"))
		}

	})

	return instance
}
