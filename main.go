package main

import (
	"pool_workers/worker_pool"
	"strconv"
	"time"
)

func main() { // time.sleep for comfort view of output, if you don't like it remove it
	workerPool := worker_pool.NewWorkerPool(0)

	workerPool.SetWorkersCount(6) // start with 6 workers
	time.Sleep(time.Second)

	for i := 0; i < 100; i++ {
		workerPool.AddTask("task №" + strconv.Itoa(i)) // add task
	}

	workerPool.AddWorker() // dynamically add worker

	time.Sleep(time.Second)
	for i := 100; i < 150; i++ {
		workerPool.AddTask("task №" + strconv.Itoa(i)) // more tasks :)
	}

	workerPool.RemoveWorker() // dynamically delete worker
	time.Sleep(time.Second * 2)

	for i := 150; i < 200; i++ {
		workerPool.AddTask("task №" + strconv.Itoa(i))
	}

}
