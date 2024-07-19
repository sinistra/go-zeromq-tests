//
//  Reading from multiple sockets.
//  This version uses zmq.Poll()
//

package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
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
	err := receiver.Connect("tcp://localhost:5555")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
	}

	//  Connect to weather server
	//subscriber, _ := zmq.NewSocket(zmq.SUB)
	//defer func(subscriber *zmq.Socket) {
	//	err := subscriber.Close()
	//	if err != nil {
	//		fmt.Println("Error closing subscriber:", err)
	//	}
	//}(subscriber)
	//err = subscriber.Connect("tcp://localhost:5556")
	//if err != nil {
	//	fmt.Println("Error connecting to server:", err)
	//}
	//err = subscriber.SetSubscribe("10001 ")
	//if err != nil {
	//	fmt.Println("Error subscribing:", err)
	//}

	//  Initialize poll set
	poller := zmq.NewPoller()
	poller.Add(receiver, zmq.POLLIN)
	//poller.Add(subscriber, zmq.POLLIN)
	//  Process messages from both sockets
	for {
		sockets, _ := poller.Poll(-1)
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case receiver:
				task, _ := s.Recv(0)
				//  Process task
				fmt.Println("Got task:", task)
			//case subscriber:
			//	update, _ := s.Recv(0)
			//	 Process weather update
			//fmt.Println("Got weather update:", update)
			default:
				fmt.Println("Unrecognized socket:", s)
			}
		}
	}
}
