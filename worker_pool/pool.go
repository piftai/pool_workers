package worker_pool

import (
	"fmt"
	"strconv"
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

func main() {
	workerPool := NewWorkerPool()

	// Стартуем несколько воркеров
	workerPool.StartWorker(1)
	time.Sleep(time.Second * 2)
	workerPool.StartWorker(2)
	workerPool.StartWorker(3)
	workerPool.StartWorker(4)
	workerPool.StartWorker(5)

	// Добавляем задачи
	for i := 0; i < 100; i++ {
		workerPool.TaskChan <- "task №" + strconv.Itoa(i)
	}
	time.Sleep(time.Second)

	// Динамически добавляем воркера
	workerPool.StartWorker(3)
	for i := 100; i < 150; i++ {
		workerPool.TaskChan <- "task №" + strconv.Itoa(i)
	}

	// Динамически удаляем воркера
	workerPool.StopWorker(1)
	fmt.Println("Worker №1 was deleted.")

	time.Sleep(time.Second * 10)

	// Завершаем оставшиеся задачи
	for i := 150; i < 1000; i++ {
		workerPool.TaskChan <- "task №" + strconv.Itoa(i)
	}

	time.Sleep(time.Second)
	close(workerPool.TaskChan)
}
