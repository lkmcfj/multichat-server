package main

type config struct {
	WSpath    string `json:"ws-path"`
	Host      string `json:"host"`
	SecretKey string `json:"secret-key"`
	//TODO: wss
}

var globalConfig config

func loadConfig() error {
	globalConfig.WSpath = "/"
	globalConfig.Host = "0.0.0.0:6780"
	globalConfig.SecretKey = "password"
	return nil
}
