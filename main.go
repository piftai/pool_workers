package main

import (
	"fmt"
	"pool_workers/worker_pool"
	"strconv"
	"time"
)

func main() {
	workerPool := worker_pool.NewWorkerPool()

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
