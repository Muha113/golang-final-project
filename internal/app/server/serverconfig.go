package server

import (
	"encoding/json"
	"io/ioutil"
)

//Config : represent configuration json file
type Config struct {
	Host               string `json:"host"`
	Port               string `json:"port"`
	LogLevel           string `json:"logLevel"`
	DbConnectionString string `json:"dbConnectionString"`
	DbName             string `json:"dbName"`
	Jwt                string `json:"jwt"`
}

//NewConfig : read json file into Config struct
func NewConfig(cfgPath string) (*Config, error) {
	cfgJSON, err := ioutil.ReadFile(cfgPath)
	cfg := Config{}
	err = json.Unmarshal(cfgJSON, &cfg)
	return &cfg, err
}
