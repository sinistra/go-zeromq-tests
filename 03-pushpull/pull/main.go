//
//  Task sink.
//  Binds PULL socket to tcp://localhost:5558
//  Collects results from workers via that socket
//

package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
	"time"
)

func main() {
	//  Prepare our socket
	receiver, _ := zmq.NewSocket(zmq.PULL)
	defer func(receiver *zmq.Socket) {
		err := receiver.Close()
		if err != nil {
			fmt.Println("Error closing receiver:", err)
		}
	}(receiver)
	err := receiver.Bind("tcp://*:5558")
	if err != nil {
		fmt.Println("Error binding to tcp://*:5558:", err)
	}

	//  Wait for start of batch
	_, err = receiver.Recv(0)
	if err != nil {
		fmt.Println("Error receiving message:", err)
	}

	//  Start our clock now
	start_time := time.Now()

	//  Process 100 confirmations
	for task_nbr := 0; task_nbr < 100; task_nbr++ {
		_, err := receiver.Recv(0)
		if err != nil {
			fmt.Println("Error receiving message:", err)
		}
		if task_nbr%10 == 0 {
			fmt.Print(":")
		} else {
			fmt.Print(".")
		}
	}

	//  Calculate and report duration of batch
	fmt.Println("\nTotal elapsed time:", time.Since(start_time))
}
