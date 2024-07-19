//
//  Reading from multiple sockets.
//  This version uses a simple recv loop
//

package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
	"time"
)

func main() {

	//  Connect to task ventilator
	receiver, _ := zmq.NewSocket(zmq.PULL)
	defer func(receiver *zmq.Socket) {
		err := receiver.Close()
		if err != nil {
			fmt.Println("Error closing receiver:", err)
		}
	}(receiver)
	err := receiver.Connect("tcp://localhost:5557")
	if err != nil {
		fmt.Println("Error connecting to receiver:", err)
	}

	//  Connect to weather server
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer func(subscriber *zmq.Socket) {
		err := subscriber.Close()
		if err != nil {
			fmt.Println("Error closing subscriber:", err)
		}
	}(subscriber)
	err = subscriber.Connect("tcp://localhost:5556")
	if err != nil {
		fmt.Println("Error connecting to subscriber:", err)
	}
	err = subscriber.SetSubscribe("10001 ")
	if err != nil {
		fmt.Println("Error subscribing:", err)
	}

	//  Process messages from both sockets
	//  We prioritize traffic from the task ventilator
	for {

		//  Process any waiting tasks
		for {
			task, err := receiver.Recv(zmq.DONTWAIT)
			if err != nil {
				break
			}
			//  process task
			fmt.Println("Got task:", task)
		}

		//  Process any waiting weather updates
		for {
			udate, err := subscriber.Recv(zmq.DONTWAIT)
			if err != nil {
				break
			}
			//  process weather update
			fmt.Println("Got weather update:", udate)
		}

		//  No activity, so sleep for 1 msec
		time.Sleep(time.Millisecond)
	}
}
