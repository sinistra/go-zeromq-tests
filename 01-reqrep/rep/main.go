//
//  Hello World server.
//  Binds REP socket to tcp://*:5555
//  Expects "Hello" from client, replies with "World"
//

package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
	"time"
)

func main() {
	//  Socket to talk to clients
	responder, _ := zmq.NewSocket(zmq.REP)
	defer func(responder *zmq.Socket) {
		err := responder.Close()
		if err != nil {
			fmt.Println("Error closing zmq socket:", err)
		}
	}(responder)
	err := responder.Bind("tcp://*:5555")
	if err != nil {
		fmt.Println("Error binding to tcp://*:5555:", err)
	}

	for {
		//  Wait for next request from client
		msg, _ := responder.Recv(0)
		fmt.Println("Received ", msg)

		//  Do some 'work'
		time.Sleep(time.Second)

		//  Send reply back to client
		reply := "World"
		responder.Send(reply, 0)
	}
}
