package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	ContestID       string
	Handles         []string
	RefreshInterval int
}

func parseConfigFile() *tomlConfig {
	var conf tomlConfig
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		if err != nil {
			log.Fatalln(err)
		}
	}

	return &conf
}
