package config

import (
	_ "crypto/sha256"
	"encoding/json"
	_ "github.com/tent/hawk-go"
	"github.com/tent/tent-client-go"
	"io/ioutil"
	"os"
	// "fmt"
)

type Config struct {
	ID  string
	Key string
	App string

	Servers []tent.MetaPostServer

	Player string
}

func Write(config Config, configFilePath string) (err error) {
	enc, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(configFilePath, enc, 0644)
	if err != nil {
		return
	}

	return
}

func Read(configFilePath string) (config Config, err error) {
	configFile, err := ioutil.ReadFile(configFilePath)

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return
	}

	return
}

func Exists(path string) (exists bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
