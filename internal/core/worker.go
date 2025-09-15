package core

import (
	"encoding/json"
	"fmt"
	"time"
)

type Worker struct {
	ID            string
	TaskChannel   <-chan Task
	ResultChannel chan<- TaskResult
	QuitChannel   <-chan bool
}

type WorkerPool struct {
	Workers       []*Worker
	TaskChannel   <-chan Task
	ResultChannel chan<- TaskResult
	QuitChannel   chan<- bool
	Size          int
}

func NewWorker(id string, taskChan <-chan Task, resultChan chan<- TaskResult, quitChan <-chan bool) *Worker {
	return &Worker{ID: id, TaskChannel: taskChan, ResultChannel: resultChan, QuitChannel: quitChan}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case task := <-w.TaskChannel:
				result := w.processTask(task)
				w.ResultChannel <- result
			case <-w.QuitChannel:
				return
			}
		}
	}()
}

func (w *Worker) processTask(task Task) TaskResult {
	startTime := time.Now()

	switch task.Type {
	case CPU_INTENSIVE:
		return w.handleCPUTask(task, startTime)
	case IO_BOUND:
		return w.handleIOTask(task, startTime)
	case TIME_BASED:
		return w.handleTimeBasedTask(task, startTime)
	default:
		return TaskResult{TaskID: task.ID, Status: FAILED, Data: task.Payload, Duration: time.Since(startTime), WorkerID: w.ID, CompletedAt: time.Now(), Error: "unknown task type"}

	}
}

func NewWorkerPool(size int, taskChan <-chan Task) *WorkerPool {
	resultChan := make(chan TaskResult)
	quitChan := make(chan bool)
	workerPool := WorkerPool{
		TaskChannel:   taskChan,
		ResultChannel: resultChan,
		QuitChannel:   quitChan,
		Size:          size,
	}

	for i := 1; i <= size; i++ {
		workerID := fmt.Sprintf("worker_%d", i)
		worker := NewWorker(workerID, taskChan, resultChan, quitChan)
		workerPool.Workers = append(workerPool.Workers, worker)
	}

	return &workerPool
}

func (wp *WorkerPool) Start() {
	for _, worker := range wp.Workers {
		worker.Start()
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.QuitChannel)
}

// BELOW THIS POINT, PLACEHOLDER FUNCTIONS
// TODO: MAKE REAL FUNCTIONS WHEN TIME IS RIGHT
func (w *Worker) handleCPUTask(task Task, startTime time.Time) TaskResult {
	// Simulate CPU work with a brief calculation
	result := 0
	for i := 0; i < 1000000; i++ {
		result += i // Just enough work to see CPU usage
	}

	return TaskResult{
		TaskID:      task.ID,
		Status:      COMPLETED,
		WorkerID:    w.ID,
		CompletedAt: time.Now(),
		Duration:    time.Since(startTime),
		Data:        json.RawMessage(fmt.Sprintf(`{"result": %d}`, result)),
	}
}

func (w *Worker) handleIOTask(task Task, startTime time.Time) TaskResult {
	// Simulate I/O delay
	time.Sleep(100 * time.Millisecond) // Pretend we're reading a file

	return TaskResult{
		TaskID:      task.ID,
		Status:      COMPLETED,
		WorkerID:    w.ID,
		CompletedAt: time.Now(),
		Duration:    time.Since(startTime),
		Data:        json.RawMessage(`{"message": "I/O operation completed"}`),
	}
}

func (w *Worker) handleTimeBasedTask(task Task, startTime time.Time) TaskResult {
	// Simulate time-based processing
	time.Sleep(500 * time.Millisecond) // Pretend we're processing something

	return TaskResult{
		TaskID:      task.ID,
		Status:      COMPLETED,
		WorkerID:    w.ID,
		CompletedAt: time.Now(),
		Duration:    time.Since(startTime),
		Data:        json.RawMessage(`{"status": "time-based processing done"}`),
	}
}
