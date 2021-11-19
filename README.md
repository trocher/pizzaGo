# PizzaGo

This is an implementation of the Pizza Maker, it runs a pizzeria that must serve a given amount of pizzas.

## Requirement

----------

To have a clean yaml configuration file, the package [yaml](https://github.com/go-yaml/yaml) was used, to install it, run:

    go get gopkg.in/yaml.v2


## Configuration

----------

To run the program, a yaml config file must be provided. The package internal/config provide a [sample file](config/config.yml) that can be used and modified. An example of such file is shown below.
```yml
# Config file for PizzaGo
# Each value must be an integer strictly positive

times: # Times to perform each step in miliseconds
  process: 1
  prepare: 2
  bake: 5
  qualityCheck: 1
parameters: 
  NumberOfWorkers: 2   # Number of workers in the pizzeria
  NumberOfOvens: 2     # Number of Ovens in the pizzeria
  NumberOfOrders: 500  # Number of Orders to take

```
## How to run it

----------

The project can be run as follow :

    go run [Path to main.go] [path to config.yml]

For example, when trying to run the project from its root using the sample config file:

    go run cmd/pizzeria/main.go config/config.yml


## Descriptions of the project

----------

### [Config](internal/configs/config.go)
The config package is simply here to read the yaml file, verify some property on the contained values and populate a ```Config``` struct that will be used by the rest of the packages

### [Timing](internal/components/timing.go)
The timing file is used to provide an abstraction to the process of waiting for the differents times given in the config file.

It also provide the calculation of the ```expected time```. This is the time it would take to treat all the order from a theoritical point of view. It does not consider the overhead produced by go and hence will always be smaller than the real elapsed time. This will be discussed later when analysing the performance of the implementation.

----------

### [Pizza Bakers](internal/components/pizzaBaker.go)

Pizza bakers are working together concurently to serve all pizzas as fast as possible. The number of pizza bakers can be configured using the file ```NumberOfBakers``` of the config file.

A pizza baker job is simple, if all the clients's orders haven't been taken yet, he takes an order and prepare the pizza. 

To bake the pizza, two cases can happen :

- If there are there are at least as many oven than pizza bakers, then an oven as been attributed to the baker so he can just use it to bake the pizza.
- Otherwise, he try to find an oven he can use, if no oven are available, he wait until one is.

Once the pizza is baked, the pizza baker can finally check its quality and repead the process for the next client.

If all clients have been served, he can go to sleep.

----------

### [Clients and Orders](internal/components/order.go)

We assume that clients are already waiting in front of the store at the begining of the "day", (i.e. bakers will never wait for client)

The number of clients is given in the ```NumberOfOrders``` field of the config file.

To monitor which client's order has already been taken, the counter ```OrderTaken``` is implemented in the file [order.go](internal/components/order.go). It is incremented by one each time a pizza baker takes a new order. Such order will have an uid corresponding to the value of the counter after this incrementation.

To ensure that two bakers does not take the same order twice in the confusion, access to this counter are done using atomic operations. Another option would have been to attribute a specific slice of the clients to each baker. (for example baker1 has orders 1 to 100, baker2 has orders 101 to 200, etc...) however this would mean that clients order's are not taking in order, which is not consistent with the reality.

----------

### [Ovens](internal/components/oven.go)

The number of ovens can be configured using the ```NumberOfOvens``` field of the config file.

An oven can hold one pizza at a time and takes the time specified in the ```bake``` field of the config file to cook a pizza.

The field ```isUsed``` of the ```PizzaOven``` struct is used by the pizza bakers looking for an empty oven to know if it is available. If it is equal to 0, this means that the oven is available, otherwise, it is set to the uid of the pizza worker currently using it.
This field is always accessed using atomic operations to ensure that two bakers doesn't try to use it a the same times. It is hences a lock.

----------

### [Pizzeria](internal/components/pizzeria.go)

This is the entry point of the pizzeria. Its function ```StartPizzeria``` is in charge of reading the config file, creating the ovens and bakers and starting each baker's Goroutine.