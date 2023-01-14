package main

import (
	"fmt"
	"log"
	"sync"

	pool "github.com/ennesuysal/go-thread-pooling"
)

var share int

func (t *Task) Execute() error {
	if t.mutex != nil {
		t.mutex.Lock()
		defer t.mutex.Unlock()
	}
	return t.executeFunc()
}

func (t *Task) OnFailure(e error) {
	log.Printf("adding error: %v", e)
}

type Task struct {
	executeFunc func() error
	mutex       *sync.Mutex
}

func newTask(executeFunc func() error, m *sync.Mutex) *Task {
	return &Task{
		executeFunc: executeFunc,
		mutex:       m,
	}
}

func main() {
	share = 0
	var m sync.Mutex

	p, _ := pool.NewPool(10, 10)
	p.Start()

	for i := 0; i < 1000; i++ {
		p.AddWork(newTask(func() error {
			share += 1
			return nil
		}, &m))
	}

	p.Stop()

	fmt.Printf("Result: %d", share)
}
