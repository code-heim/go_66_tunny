package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Jeffail/tunny"
)

type myWorker struct {
	jobID int
	state string
}

func (w myWorker) Process(payload interface{}) interface{} {
	w.jobID, _ = payload.(int)
	w.state = "processing"
	fmt.Printf("Processing job %v, state: %s\n", payload, w.state)
	time.Sleep(2 * time.Second)
	return nil
}

func (w myWorker) BlockUntilReady() {
	w.state = "starting"
	fmt.Printf("State: %s\n", w.state)
	time.Sleep(10 * time.Millisecond)
}

func (w myWorker) Interrupt() {
	w.state = "interrupted"
	fmt.Printf("State: %s\n", w.state)
	time.Sleep(10 * time.Millisecond)
}

func (w myWorker) Terminate() {
	w.state = "terminated"
	fmt.Printf("State: %s\n", w.state)
}

func main() {
	numCPUs := runtime.NumCPU()
	pool := tunny.New(numCPUs, func() tunny.Worker {
		return myWorker{}
	})
	defer pool.Close()

	for i := 0; i < 10; i++ {
		go func() {
			var data interface{} = i
			result := pool.Process(data)
			if result == nil {
				fmt.Println("success!")
			} else {
				fmt.Println("failure!")
			}

		}()
	}

	time.Sleep(time.Second)
	for {
		qLen := pool.QueueLength()
		fmt.Printf(" --------------------------------------- Queue length: %d\n", qLen)
		if qLen == 0 {
			break
		}
		time.Sleep(time.Second)
	}

	time.Sleep(5 * time.Second)

	fmt.Println("Done!")
}
