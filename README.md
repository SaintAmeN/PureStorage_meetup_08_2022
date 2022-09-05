## PureStorage - GoLang Meetup (31 Aug 2022 )
#### <i>"share by communicating"</i>
This code is an introduction and example of threading and testing. Repository consists of following components:
- example usage of channels
- example usage of channels with additional periodic thread 
- worker pool with collector
- sequential task executor

#### Example usage of channels (1_channels)
This code is a simple example showing how channels work ona Pub-Sub example.
```
go run 1_channels/*
```

#### Example usage of channels (2_channels_with_periodic_thread)
This code is a simple example showing how channels work ona Pub-Sub example. The key difference between this and the `1_channels` example is additional thread that runs asynchronously and periodically.
```
go run 2_channels_with_periodic_thread/*
```

#### Using channels, wait groups and locks to create worker pool (3_worker_pool)
In some cases we want to parallelize some work and then collect results of that work. For that, we've created executor with collector.
This code is an example of how to create a worker pool that will execute some work asynchronously. The main thread will be blocked and wait until the workers finish. 
```
go run 3_worker_pool/*
```

#### Using channels to communicate instead of locks (4_sequential_task_executor)
In some cases we want to make our work sequentiual, and for that, we've implemented sequential task executor that executes tasks one by one.
This is an example on how to replace traditional locked up code with channels. Change the way of thinking, 'share by communicating'. 
```
go run 4_sequential_task_executor/*
```


#### Final notes:
This repository contains utility that allows parallel testing with a race flag. In file: `4_sequential_task_executor/executor/queue_thread_safety_test.go` we use that utility to test our structure against any possible races.
To execute tests with data race detection use:
```
cd 4_sequential_task_executor/executor
go clean -testcache
go test . -race -v
```