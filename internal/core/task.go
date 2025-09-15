package core

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type TaskType string

const (
	CPU_INTENSIVE TaskType = "cpu_intensive"
	IO_BOUND      TaskType = "io_bound"
	TIME_BASED    TaskType = "time_based"
)

type TaskStatus string

const (
	PENDING    TaskStatus = "pending"
	PROCESSING TaskStatus = "processing"
	COMPLETED  TaskStatus = "completed"
	FAILED     TaskStatus = "failed"
	TIMEOUT    TaskStatus = "timeout"
)

type Task struct {
	ID        string          `json:"id"`
	Type      TaskType        `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Priority  int             `json:"priority"`
	CreatedAt time.Time       `json:"created_at"`
	Timeout   time.Duration   `json:"timeout"`
	Status    TaskStatus      `json:"status"`
}

type TaskResult struct {
	TaskID      string          `json:"task_id"`
	Status      TaskStatus      `json:"status"`
	Data        json.RawMessage `json:"data,omitempty"`
	Duration    time.Duration   `json:"duration"`
	WorkerID    string          `json:"worker_id"`
	CompletedAt time.Time       `json:"completed_at"`
	Error       string          `json:"error,omitempty"`
}

func NewTask(taskType TaskType, payload interface{}) (*Task, error) {
	id := fmt.Sprintf("task_%d_%d", time.Now().Unix(), time.Now().Nanosecond())

	p, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	raw := json.RawMessage(p)

	var timeout time.Duration

	switch taskType {
	case CPU_INTENSIVE:
		timeout = 30 * time.Second
	case IO_BOUND:
		timeout = 60 * time.Second
	default:
		timeout = 20 * time.Second

	}

	task := Task{ID: id, Type: taskType, CreatedAt: time.Now(), Payload: raw, Timeout: timeout, Status: PENDING}

	if err := task.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid task created: %w", err)
	}

	return &task, nil
}

func NewTaskResult(taskID string, workerID string, status TaskStatus) *TaskResult {
	return &TaskResult{TaskID: taskID, Status: status, WorkerID: workerID, CompletedAt: time.Now()}
}

func (t *Task) IsValid() error {
	if t.ID == "" {
		return fmt.Errorf("id must be supplied")
	}

	if !t.Type.IsValid() {
		return fmt.Errorf("task type '%s' is not valid", t.Type)
	}

	if t.CreatedAt.IsZero() {
		return fmt.Errorf("created at time is zero time")
	}

	if t.Timeout <= 0 {
		return fmt.Errorf("timeout must be a positive number")
	}

	return nil
}

func (r *TaskResult) SetData(data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	raw := json.RawMessage(d)

	r.Data = raw

	return nil
}

func (r *TaskResult) SetError(err error) {
	r.Status = FAILED
	r.Error = err.Error()
	log.Printf("Task %s failed: %v", r.TaskID, err)
}

func (r *TaskResult) String() string {
	s := fmt.Sprintf("TaskID: %s, Status: %s, Duration: %v, WorkerID: %s\n", r.TaskID, r.Status, r.Duration, r.WorkerID)
	return s
}

func (tt TaskType) IsValid() bool {
	switch tt {
	case CPU_INTENSIVE, IO_BOUND, TIME_BASED:
		return true
	default:
		return false
	}
}

func (ts TaskStatus) IsValid() bool {
	switch ts {
	case PENDING, PROCESSING, COMPLETED, FAILED, TIMEOUT:
		return true
	default:
		return false
	}
}

func (t *Task) IsExpired() bool {
	elapsed := time.Since(t.CreatedAt)

	return elapsed > t.Timeout
}

func (t *Task) UnmarshalPayload(target interface{}) error {
	if err := json.Unmarshal(t.Payload, &target); err != nil {
		return err
	}

	return nil
}

func (t *Task) String() string {
	s := fmt.Sprintf("ID: %s, Type: %s, Status: %s, CreatedAt: %v\n", t.ID, t.Type, t.Status, t.CreatedAt)
	return s
}
