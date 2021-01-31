package configx

import (
	"fmt"
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

type AuthenticationConfig struct {
	Keystore struct {
		URL string `json:"url"`
	} `json:"keystore"`
}

type AppConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Port        int    `json:"port"`
	Environment string `json:"environment"`
	Domain      string `json:"domain"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type UrlProp struct {
	Url string `json:"url"`
}

type Services struct {
	MDB UrlProp `json:"mdb"`
}

func (c DatabaseConfig) Validate() error {
	if c.Host == "" {
		return errors.New("database host is required")
	}

	if c.Port == 0 {
		return errors.New("database port is required")
	}

	if c.User == "" {
		return errors.New("database user is required")
	}

	if c.Database == "" {
		return errors.New("database name is required")
	}

	return nil
}

func (c DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=require statement_cache_mode=describe",
		c.Host, c.Port, c.User, c.Password, c.Database)
}

type Configuration struct {
	App   AppConfig            `json:"app"`
	Auth  AuthenticationConfig `json:"auth"`
	Log   LogConfig            `json:"log"`
	DB    DatabaseConfig       `json:"db"`
	Redis RedisConfig          `json:"redis"`
	Nats  struct {
		Url string `json:"url"`
	} `json:"nats"`
	Disks struct {
		S3 struct {
			Key      string
			Secret   string
			Endpoint string
			Bucket   string
		}
	}
	Services Services `json:"services"`
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
