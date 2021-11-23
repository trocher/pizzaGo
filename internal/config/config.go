// The package config is used to configure the pizzeria parameters
package config

import (
	"log"
	"os"
	"strconv"

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

// Read the yaml specified as an argument to initialize the values of the given config
func ReadConfig(config *Config) {
	// If no path to yaml is specified, exit
	switch len(os.Args) {
	case 2:
		yamlConfig(config)
	case 8:
		manualConfig(config)
	default:
		log.Fatal("Incorrect arguments")

	}
	// Verify that the values of the configurations are correct
	if !verifyConfig(config) {
		log.Fatal("The config file is not valid")
	}

}

func readArg(i uint64) uint64 {
	x, err := strconv.ParseUint(os.Args[i], 10, 64)
	if err != nil {
		log.Fatal("Incorrect arguments")
	}
	return x

}
func manualConfig(config *Config) {
	config.Times.Process = readArg(1)
	config.Times.Prepare = readArg(2)
	config.Times.Bake = readArg(3)
	config.Times.QualityCheck = readArg(4)
	config.Parameters.NumberOfWorkers = readArg(5)
	config.Parameters.NumberOfOvens = readArg(6)
	config.Parameters.NumberOfOrders = readArg(7)
}

func yamlConfig(config *Config) {
	f, error := os.Open(os.Args[1])
	// If there was an error while opening the file, exit
	if error != nil {
		log.Fatal(error)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	error = decoder.Decode(&config)
	if error != nil {
		log.Fatal(error)
	}
}

// Verify that the yaml does not contain null values
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
