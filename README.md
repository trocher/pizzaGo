# PizzaGo

This is an implementation of the Pizza Maker, it runs a pizzeria that must serve a given amount of pizzas.

## Requirement

To be able to use a yaml configuration file, the package [yaml](https://github.com/go-yaml/yaml) was used, to install it, run:

    go get gopkg.in/yaml.v2


## Configuration

To run the program, two possibilities are available :

### Using a yaml config file

The folder config provide a [sample file](config/config.yml) that can be used and modified. An example of such file is shown below.
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

### By specifying all arguments on the commandline

For the sake of being able to run the program many times with different parameters, it is also possible to input all the parameters values in the ```go run``` command in the order described below.
## How to run it

The project can be run as follow :

    go run [Path to main.go] [path to config.yml]

or as follow :

    go run [Path to main.go] [processTime] [preparationTime] [bakingTime] [qualityCheckTime] [nbOfWorkers] [nbOfOvens] [nbOfOrders]

For example, when trying to run the project from its root using the sample config file:

    go run cmd/pizzeria/main.go configs/config.yml

Or when trying to run the project from its root using the values as above :

    go run cmd/pizzeria/main.go 1 2 5 1 2 2 500


## Descriptions of the project

### [Config](internal/config/config.go)
The config package is simply here to read the yaml file, verify some property on the contained values and populate a ```Config``` struct that will be used by the rest of the packages.

### [Timing](internal/components/timing.go)
The timing file is used to provide an abstraction to the process of waiting for the differents times given in the config file.

It also provide the calculation of the ```expected time```. This is the time it would take to treat all the order from a theoritical point of view. It does not consider the overhead produced by go and hence should always be smaller than the real elapsed time. This will be discussed later when analysing the performance of the implementation.

### [Pizza Bakers](internal/components/pizzaBaker.go)

Pizza bakers are working together concurently to serve all pizzas as fast as possible. The number of pizza bakers can be configured using the ```NumberOfBakers``` field of the config file.

A baker as an uid that goes from 1 to ```NumberOfBakers```.

A pizza baker job is simple, if all the clients's orders haven't been taken yet, he takes an order and prepare the pizza. 

To bake the pizza, two cases can happen :

- If there are there are at least as many oven than pizza bakers, then an oven as been attributed to the baker so he can just use it to bake the pizza.
- Otherwise, he try to find an oven he can use, if no oven are available, he wait until one is.

Once the pizza is baked, the pizza baker can finally check its quality and repead the process for the next client.

If all clients have been served, he can go to sleep.

### [Clients and Orders](internal/components/order.go)

We assume that clients are already waiting in front of the store at the begining of the "day", (i.e. bakers will never wait for clients)

The number of clients is given in the ```NumberOfOrders``` field of the config file.

To monitor which client's order has already been taken, the counter ```OrderTaken``` is implemented in the file [order.go](internal/components/order.go). It is incremented by one each time a pizza baker takes a new order. Such order will have an uid corresponding to the value of the counter after this incrementation.

An order will hence have an uid that goes from 1 to ```OrderTaken```

To ensure that two bakers does not take the same order twice in the confusion, access to this counter are done using atomic operations. Another option would have been to attribute a specific slice of the client list to each baker. (for example baker1 has orders 1 to 100, baker2 has orders 101 to 200, etc...) however this would mean that clients order's are not taking in order, which is not consistent with the reality.

### [Ovens](internal/components/oven.go)

The number of ovens can be configured using the ```NumberOfOvens``` field of the config file.

An oven can hold one pizza at a time and takes the time specified in the config file's ```bake``` field to cook a pizza.

The field ```isUsed``` of the ```PizzaOven``` struct is used by the pizza bakers looking for an empty oven to know if it is available. If it is equal to 0, this means that the oven is available, otherwise, it is set to the uid of the pizza worker currently using it.
This field is always read or written using atomic operations to ensure that two bakers doesn't try to use it a the same times. It act as a lock.

### [Pizzeria](internal/components/pizzeria.go)

This is the entry point of the pizzeria. Its function ```StartPizzeria``` is in charge of reading the config file, creating the ovens and bakers and starting each baker's Goroutine. The function will return once each goroutine has finished.

## Concurency in the project

At serveral point of the project, concurency has been used, indeed, each worker is runing in its own Goroutine, all elements that are accessed by more than one of them must been accessed carefully.

### The order counter
To increment the client queue ```OrderTaken```, bakers are using ```atomic.AddUint64``` to ensure their incrementation is indeed done properly.

Before taking the order of a new client, bakers are checking the following condition to know when there are no more clients:

````go
atomic.LoadUint64(&OrderTaken) < Config.Parameters.NumberOfOrders)
````

One issue is that between this check and the incrementation of ```OrderTaken```, it is possible that more bakers than needed manage to pass the condition and hence we would end up with more order taken than the actual number of order.

To avoid this issues we could have been using a lock, however we check the following condition after incrementing the counter instead :

````go
orderId := atomic.AddUint64(&OrderTaken, 1)
	if orderId > Config.Parameters.NumberOfOrders {
		atomic.SwapUint64(&OrderTaken, Config.Parameters.NumberOfOrders)
		return Order{orderId}, errors.New("Oups, I took an order that was out of bound")
	}
````
It consist in reverting the counter if one worker went "too far".

### The ovens

Ovens are the other objects that are accessed by all the workers. To improve a bit the efficiency of the pizzeria, as explained earlier, two cases are possible :

- If there are there are at least as many oven than pizza bakers, then an oven as been attributed to the baker based on its uid. Hence each worker will only access his own oven and there are no concurrency involved

- The second case is more interesting. If there are more pizza bakers than ovens, then at some points it might be the case that 2 pizza bakers try to use the same oven concurrently. 

We will talk about the second case here. To ensure that only one baker use a given oven at a time, we used atomic operations again. The code that is in charge of locking an oven and releasing it can be found in the two methods ```findOven``` and ```releaseOven``` of [pizzaBaker](internal/components/pizzaBaker.go).

The field ```isUsed``` of an oven is its lock, indeed, when a baker is looking for an oven, he will call ```findOven``` that will look throught all the oven trying to perform the following :
````go
atomic.CompareAndSwapUint64(&(ovenList[i].isUsed), 0, w.Name)
````
When the ```CompareAndSwapUint64``` returns true, he know for sure that his uid has been placed in ```isUsed``` and that the oven is his as anyone trying to claim using the ```CompareAndSwapUint64``` will fail as ```isUsed``` is no longer 0.

To release the oven, the worker atomically swap the value back to 0 in a similar fashion :
```go
atomic.CompareAndSwapUint64(&(o.isUsed), w.Name, 0)
```

### The delivered Orders counter

The counter ```OrderDelivered``` is the last object accessed concurrently, the same logic as the order counter is used here as it is only incremented using the following function:
```go
atomic.AddUint64(&OrderDelivered, 1)
```
This value is mostly used to debug the program and ensure its correctness so we wont discuss it further here.

## Analysis of the throughput and the latency

When analysing both the throughput and the latency of the implementation, we are going to use two notions : the theoritical values and the real ones, we describe now the difference.

### Theoritical versus Real values

Compared to the theoritical values, the real one can varie a lot depending on multiple factors. For example, the OS or the load on the computer at the moment the code is ran can have influence on its running time.

One major variable found when implementing the project was the time taken by the following call used to simulate the time taken by each step of the pizza making:
```go
time.Sleep(time.Duration(t) * time.Millisecond)
```
indeed, according to the ```time``` package [documentation](https://pkg.go.dev/time#Sleep),
> Sleep pauses the current goroutine for at least the duration d

We hence have no guarranty that the sleep function will sleep exactly for the time specified, depending on the environment, it might in fact be much more. For example using two differents environment, when trying to sleep for 5ms multiple times, one was in fact sleeping for 6.25 ms in average when the other one was sleeping for 5.20 ms in average. When doing some tests in differents machine, it appear that linux configurations tends to be the more precise, we will hence benchmark the pizzeria on a computer running ubuntu 18.04.

### Theoritical values

The throughput of the implementation is quite easy to compute. We first compute the bottleneck of the pizzeria :

````
bottleneck = min(NumberOfWorkers,NumberOfOvens)
````
we then compute the time taken to produce a single pizza:
````go
pizzaCookingTime = bakingTime+preparationTime+processTime+qualityCheckTime (in ms)
````
We hence get :
```
throughput = bottleneck/pizzaCookingTime (pizza/ms)
```

For example, using 2 workers 