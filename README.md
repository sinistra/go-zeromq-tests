go-zeromq-tests
===============

Testing 0mq (using Go)

Basically, copying [pebbe/zmq4/examples](https://github.com/pebbe/zmq4/tree/master/examples).

## Examples

### 01-reqrep

`clients (REQ) * --- 1 server (REP)`

Synchronous requests?

* Client sends a request
* Server does something with the request
* Server returns something to client

### 02-pubsub

`clients (SUB) * --- 1 server (PUB)`

Event based?

* Publisher... published a number of messages with a prefix
* Subscriber... subscribes to messages starting with a given prefix
* Publisher only sends messages to a subscriber if he subscribed to it
* If publisher dies, subscribers don't break
* When publisher ressurects, subscribers carry on as if nothing happened

### 03-pushpull

`pusher nodes (PUSH) * --- * puller nodes (PULL)`

Round-robin / load-balancer / fan-out / workers?

* Pusher sends messages to registered pullers following a rota
* Pullers receive a message and do something with it
* Pullers can be added / removed dynamically and pusher always knows how many "workers" it can push to