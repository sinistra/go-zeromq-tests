//
//  Task ventilator.
//  Binds PUSH socket to tcp://localhost:5557
//  Sends batch of tasks to workers via that socket
//

package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
	"math/rand"
	"time"
)

func main() {
	//  Socket to send messages on
	sender, _ := zmq.NewSocket(zmq.PUSH)
	defer func(sender *zmq.Socket) {
		err := sender.Close()
		if err != nil {
			fmt.Println("Error while closing socket:", err)
		}
	}(sender)
	err := sender.Bind("tcp://*:5557")
	if err != nil {
		fmt.Println("Error while binding to tcp://*:5557:", err)
	}

	//  Socket to send start of batch message on
	sink, _ := zmq.NewSocket(zmq.PUSH)
	defer func(sink *zmq.Socket) {
		err := sink.Close()
		if err != nil {
			fmt.Println("Error while closing socket:", err)
		}
	}(sink)
	err = sink.Connect("tcp://localhost:5558")
	if err != nil {
		fmt.Println("Error while connecting to tcp://localhost:5558:", err)
	}

	fmt.Print("Press Enter when the workers are ready: ")
	var line string
	_, err = fmt.Scanln(&line)
	if err != nil {
		fmt.Println("Error while reading input:", err)
	}
	fmt.Println("Sending tasks to workers...")

	//  The first message is "0" and signals start of batch
	_, err = sink.Send("0", 0)
	if err != nil {
		fmt.Println("Error while sending tasks to workers:", err)
	}

	//  Initialize random number generator
	rand.Seed(time.Now().UnixNano())

	//  Send 100 tasks
	total_msec := 0
	for task_nbr := 0; task_nbr < 100; task_nbr++ {
		//  Random workload from 1 to 100msecs
		workload := rand.Intn(100) + 1
		total_msec += workload
		s := fmt.Sprintf("%d", workload)
		_, err := sender.Send(s, 0)
		if err != nil {
			fmt.Println("Error while sending tasks to workers:", err)
		}
	}
	fmt.Println("Total expected cost:", time.Duration(total_msec)*time.Millisecond)
	time.Sleep(time.Second) //  Give 0MQ time to deliver

}
