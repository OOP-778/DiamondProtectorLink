package config

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
)

const DefaultLocation = "config.yml"

type Configuration struct {
	RedisHostname string
	RedisPort     int
	RedisPassword string
}

func Get() Configuration {

	// Check for file existence
	_, err := os.Stat(DefaultLocation)
	var exists = !os.IsNotExist(err)

	var currentConfig = GetDefaultConfig()

	if !exists {
		saveConfig(currentConfig)
		return currentConfig
	}

	currentConfig = readConfig()
	saveConfig(currentConfig)
	return currentConfig
}

func readConfig() Configuration {
	readFile, _ := ioutil.ReadFile(DefaultLocation)
	var config = GetDefaultConfig()

	err2 := yaml.Unmarshal(readFile, &config)
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println(config)
	return config
}

func saveConfig(config Configuration) {
	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile(DefaultLocation, data, 0)
	if err2 != nil {
		log.Fatal(err2)
	}
}

func GetDefaultConfig() Configuration {
	return Configuration{"localhost", 6379, ""}
}
