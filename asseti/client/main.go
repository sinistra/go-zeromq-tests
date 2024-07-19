//
//  ROUTER-to-DEALER example.
//

package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
	"math/rand"
	"time"
)

const (
	NbrWorkers = 10
)

func workerTask() {
	worker, _ := zmq.NewSocket(zmq.DEALER)
	defer func(worker *zmq.Socket) {
		err := worker.Close()
		if err != nil {
			fmt.Println("Error closing worker:", err)
		}
	}(worker)
	setId(worker) //  Set a printable identity
	err := worker.Connect("tcp://localhost:5555")
	if err != nil {
		fmt.Println("Error connecting to worker:", err)
	}

	total := 0
	for {
		//  Tell the broker we're ready for work
		_, err := worker.Send("", zmq.SNDMORE)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
		_, err = worker.Send("Hi Boss", 0)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}

		//  Get workload from broker, until finished
		_, err = worker.Recv(0)
		if err != nil {
			fmt.Println("Error receiving message:", err)
		} //  Envelope delimiter
		workload, _ := worker.Recv(0)
		if workload == "Fired!" {
			fmt.Printf("Completed: %d tasks\n", total)
			break
		}
		total++

		//  Do some random work
		time.Sleep(time.Duration(rand.Intn(500)+1) * time.Millisecond)
	}
}

func main() {
	broker, _ := zmq.NewSocket(zmq.ROUTER)
	defer func(broker *zmq.Socket) {
		err := broker.Close()
		if err != nil {
			fmt.Println("Error closing broker:", err)
		}
	}(broker)

	err := broker.Bind("tcp://*:5555")
	if err != nil {
		fmt.Println("Error connecting to broker:", err)
	}
	rand.Seed(time.Now().UnixNano())

	for workerNbr := 0; workerNbr < NbrWorkers; workerNbr++ {
		go workerTask()
	}
	//  Run for five seconds and then tell workers to end
	startTime := time.Now()
	workersFired := 0
	for {
		//  Next message gives us least recently used worker
		identity, _ := broker.Recv(0)
		_, err := broker.Send(identity, zmq.SNDMORE)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
		_, err = broker.Recv(0)
		if err != nil {
			fmt.Println("Error receiving message:", err)
		} //  Envelope delimiter
		_, err = broker.Recv(0)
		if err != nil {
			fmt.Println("Error receiving message:", err)
		} //  Response from worker
		_, err = broker.Send("", zmq.SNDMORE)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}

		//  Encourage workers until it's time to fire them
		if time.Since(startTime) < 5*time.Second {
			_, err := broker.Send("Work harder", 0)
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
		} else {
			_, err := broker.Send("Fired!", 0)
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
			workersFired++
			if workersFired == NbrWorkers {
				break
			}
		}
	}

	time.Sleep(time.Second)
}

func setId(soc *zmq.Socket) {
	identity := fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
	err := soc.SetIdentity(identity)
	if err != nil {
		fmt.Println("Error setting identity:", err)
	}
}
