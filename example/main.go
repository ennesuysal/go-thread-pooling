package main

import (
	"fmt"
	"log"
	"sync"

	pool "github.com/ennesuysal/go-thread-pooling"
)

var share int

func (t *Task) Execute(i interface{}) error {
	if t.mutex != nil {
		t.mutex.Lock()
		defer t.mutex.Unlock()
	}
	return t.executeFunc(i)
}

func (t *Task) OnFailure(e error) {
	log.Printf("adding error: %v", e)
}

type Task struct {
	executeFunc func(interface{}) error
	mutex       *sync.Mutex
}

func newTask(executeFunc func(interface{}) error, m *sync.Mutex, parameters interface{}) pool.Task {
	return pool.Task{
		Exec: &Task{
			executeFunc: executeFunc,
			mutex:       m,
		},
		Parameters: parameters,
	}
}

func main() {
	share = 0
	var m sync.Mutex

	p, _ := pool.NewPool(10, 10)
	p.Start()

	for i := 0; i < 1000; i++ {
		p.AddWork(newTask(func(param interface{}) error {
			share += 1
			fmt.Printf("%d. Job\n", param)
			return nil
		}, &m, i))
	}

	p.Stop()

	fmt.Printf("Result: %d", share)
}
