package test

import (
	"fmt"
	"testing"
)

// 1. Define task
func ProcessNumber(n int) (int, error) {
	return n * n, nil
}

type Result struct {
	workerId int
	n        int
	result   int
	err      error
}

// 2. Create the worker
// Worker: goroutine that processes tasks and send the results through a channel.
type Worker struct {
	id         int
	taskQueue  <-chan int
	resultChan chan<- Result
}

func (w *Worker) Start() {
	go func() {
		fmt.Printf("Starting worker with id %d\n", w.id)
		for n := range w.taskQueue {
			data, err := ProcessNumber(n)
			w.resultChan <- Result{
				workerId: w.id,
				n:        n,
				result:   data,
				err:      err,
			}
		}
	}()
}

// 3. Implement the worker pool
// Worker pool: manages the worker, distributes task, and collects results.
type WorkerPool struct {
	taskQueue   chan int
	resultChan  chan Result
	workerCount int
}

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		taskQueue:   make(chan int, workerCount),
		resultChan:  make(chan Result, workerCount),
		workerCount: workerCount,
	}
}

func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.workerCount; i++ {
		worker := Worker{id: i, taskQueue: wp.taskQueue, resultChan: wp.resultChan}
		worker.Start()
	}
}

func (wp *WorkerPool) Submit(n int) {
	fmt.Println("Submit to task queue...")
	wp.taskQueue <- n
}

func (wp *WorkerPool) GetResult() Result {
	fmt.Println("Getting result...")
	return <-wp.resultChan
}

func TestWorkerSimV1(t *testing.T) {
	// Generate task
	var tasks []int
	for i := 0; i < 15; i++ {
		tasks = append(tasks, i)
	}

	workerPool := NewWorkerPool(5)
	// Deploy worker
	workerPool.Start()

	// Submit tasks to worker from worker pool
	for _, n := range tasks {
		workerPool.Submit(n)
	}

	// Collect results and handle error
	for i := 0; i < len(tasks); i++ {
		result := workerPool.GetResult()
		if result.err != nil {
			fmt.Printf("ERROR by worker with id %d : %+v\n", result.workerId, result.err)
		} else {
			fmt.Printf("SUCCESS: Worker id: %d, N: %d, Result: %d\n", result.workerId, result.n, result.result)
		}
	}
}
