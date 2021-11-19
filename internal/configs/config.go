// The package configs is used to configure the pizzeria parameters
package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config type used to store all the parameters that can be modified into the yaml
type Config struct {
	Times struct {
		Process      uint64 `yaml:"process"`
		Prepare      uint64 `yaml:"prepare"`
		Bake         uint64 `yaml:"bake"`
		QualityCheck uint64 `yaml:"qualityCheck"`
	} `yaml:"times"`
	Parameters struct {
		NumberOfWorkers uint64 `yaml:"NumberOfWorkers"`
		NumberOfOvens   uint64 `yaml:"NumberOfOvens"`
		NumberOfOrders  uint64 `yaml:"NumberOfOrders"`
	} `yaml:"parameters"`
}

func ReadConfig(config *Config) {
	if len(os.Args) == 1 {
		log.Fatal("Please specify a config file")
	}
	f, error := os.Open(os.Args[1])
	if error != nil {
		log.Fatal(error)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	error = decoder.Decode(&config)
	if error != nil {
		log.Fatal(error)
	}
	if !verifyConfig(config) {
		log.Fatal("The config file is not valid")
	}

}
func verifyConfig(config *Config) bool {
	if config.Times.Process == 0 ||
		config.Times.Prepare == 0 ||
		config.Times.Bake == 0 ||
		config.Times.QualityCheck == 0 ||
		config.Parameters.NumberOfWorkers == 0 ||
		config.Parameters.NumberOfOvens == 0 ||
		config.Parameters.NumberOfOrders == 0 {
		return false
	}
	return true
}
