package core

import "fmt"

type TaskQueue struct {
	Tasks   chan Task
	MaxSize int
}

func NewTaskQueue(capacity int) *TaskQueue {
	tasksChan := make(chan Task, capacity)

	taskQueue := TaskQueue{Tasks: tasksChan, MaxSize: capacity}

	return &taskQueue
}

func (tq *TaskQueue) Submit(task *Task) error {
	select {
	case tq.Tasks <- *task:
		return nil
	default:
		return fmt.Errorf("queue full, please try again shortly")
	}
}

func (tq *TaskQueue) GetTaskChannel() <-chan Task {
	return tq.Tasks
}

func (tq *TaskQueue) Close() {
	close(tq.Tasks)
}

func (tq *TaskQueue) Status() (current, max int) {
	return len(tq.Tasks), tq.MaxSize
}
