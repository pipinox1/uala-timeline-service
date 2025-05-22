package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	ServiceName string      `mapstructure:"service_name"`
	Env         string      `mapstructure:"env"`
	Port        string      `mapstructure:"port"`
	Postgres    Postgres    `mapstructure:"postgres"`
	AWS         AWS         `mapstructure:"aws"`
	RestConfigs RestConfigs `mapstructure:"rest_configs"`
	Nats        Nats        `mapstructure:"nats"`
}

type Nats struct {
	Host string `mapstructure:"host"`
}

type RestConfigs struct {
	PostService      RestConfig `mapstructure:"post_service"`
	FollowersService RestConfig `mapstructure:"followers_service"`
}

type RestConfig struct {
	BasePath string `mapstructure:"base_path"`
	Timeout  int    `mapstructure:"timeout"`
}

type AWS struct {
	Region  string `mapstructure:"region"`
	Table   string `mapstructure:"table"`
	Secret  string `mapstructure:"secret"`
	Account string `mapstructure:"account"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	UseSSL   bool   `mapstructure:"use_ssl"`
}

func ReadConfig() (*Config, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("unable to get current file")
	}

	configDir := filepath.Join(filepath.Dir(filename))
	viper.SetConfigName(getConfigName())
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func getConfigName() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		return "local"
	}
	return env
}
