package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type configurations struct {
	AppPort  int64    `yaml:"appPort"`
	DBConfig DBConfig `yaml:"dbConfig"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	Port     int64  `yaml:"port"`
	DBName   string `yaml:"dbName"`
}

var config configurations

func ReadAndInitConfig(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	c := &configurations{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return err
	}

	config = *c
	return nil
}

func GetAppPort() string {
	return fmt.Sprint(config.AppPort)
}

func GetDBConfig() DBConfig {
	return config.DBConfig
}
