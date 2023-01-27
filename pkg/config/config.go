package config

import (
	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Server   ServerConf      `mapstructure:"server"`
	Postgres PostgresSection `mapstructure:"postgres"`
}

type ServerConf struct {
	Port      int    `mapstructure:"port"`
	Secret    string `mapstructure:"secret"`
	Interval  int    `mapstructure:"interval"`
	WorkerNum int    `mapstructure:"worker_num"`
}

type PostgresSection struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Db       string `mapstructure:"database"`
}

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&Conf)
	return
}
