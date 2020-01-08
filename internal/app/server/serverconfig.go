package server

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Host               string `json:"host"`
	Port               string `json:"port"`
	LogLevel           string `json:"logLevel"`
	DbConnectionString string `json:"dbConnectionString"`
	DbName             string `json:"dbName"`
	Jwt                string `json:"jwt"`
}

func NewConfig(cfgPath string) (*Config, error) {
	cfgJSON, err := ioutil.ReadFile(cfgPath)
	// fmt.Println(cfgJSON)
	// fmt.Println(cfgPath)
	cfg := Config{}
	err = json.Unmarshal(cfgJSON, &cfg)
	return &cfg, err
}
