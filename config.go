package x2o

import (
	"encoding/json"
	"fmt"
	"os"
)

var Info Config

type Config struct {
	Host string
	Port int
}

type ConfigFile struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func Load(configFilePath string) {
	cfgFile := ConfigFile{}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read config failed: %v\n", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(data, &cfgFile); err != nil {
		fmt.Fprintf(os.Stderr, "parse config failed: %v\n", err)
		os.Exit(1)
	}

	if cfgFile.Host == "" {
		// 0.0.0.0 and ::
		Info.Host = "wildcard_address"
	} else {
		Info.Host = cfgFile.Host
	}

	Info.Port = cfgFile.Port
}
