package pooling

import (
	"sync"
	"time"
)

type taskTemplate interface {
	Execute(interface{}) error
	OnFailure(error)
}

type Task struct {
	Exec       taskTemplate
	Parameters interface{}
}

type Pool interface {
	Start()
	Stop()
	AddWork(Task)
}

type PoolSt struct {
	workerCount    int
	tasks          chan Task
	start          sync.Once
	stop           sync.Once
	wg             sync.WaitGroup
	wgAddingWorker sync.WaitGroup
	quit           chan struct{}
}

func (p *PoolSt) Start() {
	p.start.Do(func() {
		p.startWorkers()
	})
}

func (p *PoolSt) Stop() {
	p.stop.Do(func() {
		p.wgAddingWorker.Wait()
		for len(p.tasks) > 0 {
			time.Sleep(100 * time.Microsecond)
		}
		close(p.quit)
		p.wg.Wait()
	})
}

func (p *PoolSt) AddWorkHelper(t Task) {
	select {
	case p.tasks <- t:
	case <-p.quit:
	}
	p.wgAddingWorker.Done()
}

func (p *PoolSt) AddWork(t Task) {
	p.wgAddingWorker.Add(1)
	go p.AddWorkHelper(t)
}

func (p *PoolSt) startWorkers() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go func(workerNum int) {
			for {
				select {
				case <-p.quit:
					p.wg.Done()
					return
				case task, ok := <-p.tasks:
					if !ok {
						return
					}

					if err := task.Exec.Execute(task.Parameters); err != nil {
						task.Exec.OnFailure(err)
					}
				}
			}
		}(i)
	}
}

func NewPool(workerCount int, chanSize int) (Pool, error) {
	return &PoolSt{
		workerCount:    workerCount,
		tasks:          make(chan Task, chanSize),
		start:          sync.Once{},
		stop:           sync.Once{},
		wg:             sync.WaitGroup{},
		wgAddingWorker: sync.WaitGroup{},
		quit:           make(chan struct{}),
	}, nil
}
