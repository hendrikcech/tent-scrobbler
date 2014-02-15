package config

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/tent/hawk-go"
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

func Write(client *tent.Client, configFilePath string) (err error) {
	config := Config{
		ID:      client.Credentials.ID,
		Key:     client.Credentials.Key,
		App:     client.Credentials.App,
		Servers: client.Servers,
	}

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

func Read(configFilePath string) (client *tent.Client, err error) {
	configFile, err := ioutil.ReadFile(configFilePath)

	var config Config

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return
	}

	return &tent.Client{
		Credentials: &hawk.Credentials{
			ID:   config.ID,
			Key:  config.Key,
			App:  config.App,
			Hash: sha256.New,
		},
		Servers: config.Servers,
	}, nil
}

func Exists(path string) (exists bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
