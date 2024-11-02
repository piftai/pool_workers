package worker_pool

import (
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	TaskChan chan string
	mu       sync.Mutex
	workers  []*Worker
}

type Worker struct {
	id       int
	active   bool
	quitChan chan struct{}
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		TaskChan: make(chan string),
		workers:  make([]*Worker, 0),
	}
}

func (p *WorkerPool) StartWorker(id int) {
	worker := &Worker{
		id:       id,
		active:   true,
		quitChan: make(chan struct{}),
	}
	p.mu.Lock()
	p.workers = append(p.workers, worker)
	p.mu.Unlock()
	go func() {
		fmt.Printf("Worker №%d started.\n", id)
		for {
			select {
			case task := <-p.TaskChan:

				fmt.Printf("Worker №%d, data: %s\n", id, task) // add checker if chan closed
				time.Sleep(time.Millisecond * 100)
			case <-worker.quitChan:

				fmt.Printf("Worker №%d stopped\n", id)
				return
			}
		}
	}()
}

func (p *WorkerPool) StopWorker(id int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, worker := range p.workers {
		if worker.id == id && worker.active {
			close(worker.quitChan)
			worker.active = false
			break
		}
	}
}
