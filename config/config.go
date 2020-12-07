package config

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"strings"
)

type config struct {
	Database Database `yaml:"Database"`
	App      App      `yaml:"App"`
	Redis    Redis    `yaml:"Redis"`
}

type Database struct {
	Username           string `yaml:"USERNAME"`
	Password           string `yaml:"PASSWORD"`
	Port               string `yaml:"PORT"`
	Host               string `yaml:"HOST"`
	DBName             string `yaml:"DB_NAME"`
	MaxOpenConnections int    `yaml:"MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections int    `yaml:"MAX_IDLE_CONNECTIONS"`
}

type Redis struct {
	Address string `yaml:"ADDRESS"`
}

type App struct {
	LogLevel     bool   `yaml:"LOG_LEVEL"`
	ComQueueName string `yaml:"COM_QUEUE_NAME"`
}

var cfg config

type ConfiguredApp struct {
	DB     *sql.DB
	RDB    *redis.Client
	Config config
}

func InitializeConfig() *ConfiguredApp {
	viper.SetEnvPrefix("R1WALLET")
	viper.AddConfigPath(".")
	viper.SetConfigName("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.MergeInConfig()
	if err != nil {
		fmt.Println("Error in reading config")
		panic(err)
	}

	err = viper.Unmarshal(&cfg, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
	})
	if err != nil {
		fmt.Println("Error in unmarshaling config")
		panic(err)
	}

	fmt.Printf("\n loaded config: %#v \n", cfg)

	db, err := initializeMysql(cfg.Database)
	if err != nil {
		panic(err)
	}

	re, err := initializeRedis(cfg.Redis)
	if err != nil {
		panic(err)
	}

	return &ConfiguredApp{
		DB:     db,
		RDB:    re,
		Config: cfg,
	}
}
