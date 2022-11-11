# STATSD MUX

Statsd is yet another statsd server implementation using Golang.
This project is as POC that we can create our own statsd server and modify the logic based on our needs.

## How to run


### Run statsd server

```shell
$ go mod download
$ go run main.go
```

It will run the UDP statsd server on port 8125 by default.

Try to push the metrics using `nc` command:

```shell
echo 'gorets:1|c\nglork:320|ms\ngaugor:333|g\nuniques:765|s' | nc -w 1 -u localhost 8125
```

[Pure statsd](https://github.com/statsd/statsd) implementation supports multi-line.


### Run testing

To test with more data use `testing` directory by run:

```shell
$ cd testing
$ go mod download
$ go run main.go
```

