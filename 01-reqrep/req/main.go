//
//  Hello World client.
//  Connects REQ socket to tcp://localhost:5555
//  Sends "Hello" to server, expects "World" back
//

package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
)

func main() {
	//  Socket to talk to server
	fmt.Println("Connecting to hello world server...")
	requester, _ := zmq.NewSocket(zmq.REQ)
	defer func(requester *zmq.Socket) {
		err := requester.Close()
		if err != nil {
			fmt.Println("Error closing requester")
		}
	}(requester)
	err := requester.Connect("tcp://localhost:5555")
	if err != nil {
		fmt.Println("Error connecting to hello world")
	}

	for requestNbr := 0; requestNbr != 10; requestNbr++ {
		// send hello
		msg := fmt.Sprintf("Hello %d", requestNbr)
		fmt.Println("Sending ", msg)
		requester.Send(msg, 0)

		// Wait for reply:
		reply, _ := requester.Recv(0)
		fmt.Println("Received ", reply)
	}
}
