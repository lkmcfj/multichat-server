package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type wssConfig struct {
	Key  string `json:"keyfile"`
	Cert string `json:"certfile"`
}

type config struct {
	WSpath    string     `json:"ws-path"`
	Host      string     `json:"host"`
	SecretKey string     `json:"secret-key"`
	Wss       *wssConfig `json:"wss"`
}

var globalConfig config

func loadConfig() error {
	f, err := os.Open("config.json")
	if err != nil {
		return err
	}
	configText, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(configText, &globalConfig)
	if err != nil {
		return err
	}
	if (globalConfig.WSpath == "") || (globalConfig.Host == "") || (globalConfig.SecretKey == "") {
		return errors.New("ws-path, host and secret-key are required (and should not be empty string)")
	}
	if globalConfig.Wss != nil {
		if (globalConfig.Wss.Cert == "") || (globalConfig.Wss.Key == "") {
			return errors.New("wss.keyfile and wss.certfile are required when wss is enabled")
		}
	}
	return nil
}
