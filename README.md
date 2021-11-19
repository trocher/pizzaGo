# PizzaGo

In this implementation of the Pizza Maker, multiple parameters are configurable, you can look at the 

## Requirement
To have a clean yaml configuration file, the package [yaml](https://github.com/go-yaml/yaml) was used, to install it, run:

    go get gopkg.in/yaml.v2


## Configuration
To run the program, a yaml config file must be provided. The package internal/config provide a [sample file](config/config.yml) that can be used and modified.
## How to run it
The project can be run as follow :

    go run [Path to main.go] [path to config.yml]

For example, when trying to run the project from its root using the sample config file:

    go run cmd/pizzeria/main.go config/config.yml


## 

