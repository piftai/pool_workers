package worker_pool

import (
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	taskChan   chan string
	mu         sync.Mutex
	workers    map[int]*Worker
	wg         sync.WaitGroup
	isShutdown bool
}

type Worker struct {
	id       int
	active   bool
	quitChan chan struct{}
}

func NewWorkerPool(capChan int) *WorkerPool {
	if capChan <= 0 {
		capChan = 1000 // default size of channel
	}
	return &WorkerPool{
		taskChan:   make(chan string, capChan),
		workers:    make(map[int]*Worker), // edited in map to optimize search to O(1)
		isShutdown: false,
	}
}

func (p *WorkerPool) AddWorker() {
	worker := &Worker{
		id:       len(p.workers) + 1,
		active:   true,
		quitChan: make(chan struct{}),
	}
	p.mu.Lock()
	p.workers[worker.id] = worker
	p.mu.Unlock()
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		fmt.Printf("Worker №%d started.\n", worker.id)
		for {
			select {
			case task, ok := <-p.taskChan:
				if !ok {
					fmt.Printf("Worker №%d exiting due to shutdown.\n", worker.id)
					return
				}
				fmt.Printf("Worker №%d, data: %s\n", worker.id, task) // add checker if chan closed
				time.Sleep(time.Millisecond * 100)
			case <-worker.quitChan:

				fmt.Printf("Worker №%d stopped\n", worker.id)
				return
			}
		}
	}()
}

func (p *WorkerPool) RemoveWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for id, worker := range p.workers {
		if worker.active {
			worker.quitChan <- struct{}{}
			worker.active = false
			delete(p.workers, id)
			break
		}
	}
}

func (p *WorkerPool) AddTask(task string) {
	if !p.isShutdown {
		p.taskChan <- task
	}
}

func (p *WorkerPool) SetWorkersCount(n int) {
	p.mu.Lock()
	currentCount := len(p.workers)
	p.mu.Unlock()
	if currentCount < n { // if count of workers less than n, and we need to add
		for i := currentCount; i < n; i++ {
			p.AddWorker()
		}
	} else if currentCount > n { // if count of workers more than n
		for toRemove := currentCount - n; toRemove >= 0; toRemove-- {
			p.RemoveWorker()
		}
	}
}

func (p *WorkerPool) Shutdown() {
	p.mu.Lock()

	if p.isShutdown {
		p.mu.Unlock()
		return
	}

	p.isShutdown = true
	p.mu.Unlock()
	close(p.taskChan)
	p.wg.Wait()
	p.mu.Lock()
	for id, worker := range p.workers {
		if worker.active {
			close(worker.quitChan)
			worker.active = false
		}
		delete(p.workers, id)
	}
	p.mu.Unlock()
	fmt.Println("Worker pool shut down completed")
}
