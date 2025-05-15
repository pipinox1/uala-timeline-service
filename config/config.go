package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
)

type Config struct {
	ServiceName string   `json:" "`
	Env         string   `json:"env"`
	Port        string   `json:"port"`
	Postgres    Postgres `json:"postgres"`
	Redis       Redis    `json:"redis"`
}

type Redis struct {
	Host string `json:"redis_host"`
}

type Postgres struct {
	Host     string `json:"postgres_host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Database string `json:"database"`
	User     string `json:"user"`
	UseSSL   bool   `json:"use_ssl"`
}

func ReadConfig() (*Config, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("unable to get current file")
	}

	configDir := filepath.Join(filepath.Dir(filename))
	viper.SetConfigName("local")
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
