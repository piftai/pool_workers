package worker_pool

import (
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	taskChan chan string
	mu       sync.Mutex
	workers  map[int]*Worker
}

type Worker struct {
	id       int
	active   bool
	quitChan chan struct{}
}

func NewWorkerPool(capChan int) *WorkerPool {
	if capChan == 0 {
		capChan = 1000 // default size of channel
	}
	return &WorkerPool{
		taskChan: make(chan string, capChan),
		workers:  make(map[int]*Worker), // edited in map to optimize search to O(1)
	}
}

func (p *WorkerPool) StartWorker() {
	worker := &Worker{
		id:       len(p.workers) + 1,
		active:   true,
		quitChan: make(chan struct{}),
	}
	p.mu.Lock()
	p.workers[worker.id] = worker
	p.mu.Unlock()
	go func() {
		fmt.Printf("Worker №%d started.\n", worker.id)
		for {
			select {
			case task := <-p.taskChan:

				fmt.Printf("Worker №%d, data: %s\n", worker.id, task) // add checker if chan closed
				time.Sleep(time.Millisecond * 100)
			case <-worker.quitChan:

				fmt.Printf("Worker №%d stopped\n", worker.id)
				return
			}
		}
	}()
}

func (p *WorkerPool) StopWorker(id int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.workers[id].quitChan <- struct{}{} // O(1)
	p.workers[id].active = false
}

func (p *WorkerPool) AddTask(task string) {
	p.taskChan <- task
}

func (p *WorkerPool) SetWorkersCount(n int) {
	p.mu.Lock()
	currentCount := len(p.workers)
	p.mu.Unlock()

	if currentCount < n { // if count of workers less than n, and we need to add
		for i := currentCount; i < n; i++ {
			p.StartWorker()
		}
	} else if currentCount > n { // if count of workers more than n
		toRemove := currentCount - n
		p.mu.Lock()
		for id, worker := range p.workers {
			if toRemove == 0 {
				break
			}
			if worker.active {
				worker.quitChan <- struct{}{}
				worker.active = false
				delete(p.workers, id)
				toRemove--
			}
		}
		p.mu.Unlock()
	}
}
