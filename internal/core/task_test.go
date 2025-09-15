package core

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTask_ValidInputs(t *testing.T) {
	testCases := []struct {
		name     string
		taskType TaskType
		payload  map[string]interface{}
	}{
		{
			name:     "Test 1: Cpu intensive, no error",
			taskType: CPU_INTENSIVE,
			payload:  map[string]interface{}{"name": "Marine", "race": "Terran", "health": 45, "damage": 6, "is_upgraded": true},
		},
		{
			name:     "Test 2: IO bound, no error",
			taskType: IO_BOUND,
			payload:  map[string]interface{}{"name": "Marine", "race": "Terran", "health": 45, "damage": 6, "is_upgraded": true},
		},
		{
			name:     "Test 3: Time based, no error",
			taskType: TIME_BASED,
			payload:  map[string]interface{}{"name": "Marine", "race": "Terran", "health": 45, "damage": 6, "is_upgraded": true},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewTask(tc.taskType, tc.payload)
			if err != nil {
				t.Errorf("Test %s: No error expected but error received: %v", tc.name, err)
			}
			if result == nil {
				t.Errorf("Test %s: Expected valid TaskType, received nil", tc.name)
			}
		})
	}

}

func TestNewTask_InvalidInputs(t *testing.T) {
	testCases := []struct {
		name        string
		taskType    TaskType
		payload     map[string]interface{}
		expectedErr bool
	}{
		{
			name:     "Test 1: Unexpected Task Type, invalid task type error expected",
			taskType: "gpu_bound",
			payload:  map[string]interface{}{"name": "Marine", "race": "Terran", "health": 45, "damage": 6, "is_upgraded": true},
		},
		{
			name:     "Test 2: Unmarshalable JSON, marshaling error expected",
			taskType: CPU_INTENSIVE,
			payload:  map[string]interface{}{"name": "Marine", "race": "Terran", "health": 45, "damage": 6, "is_upgraded": true, "attack": func() string { return "C-14 Gauss Rifle attack!" }},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewTask(tc.taskType, tc.payload)
			if err == nil {
				t.Errorf("Test %s: Expected error but no error received", tc.name)
			}
		})
	}

}

func TestTask_IsValid(t *testing.T) {
	testCases := []struct {
		name        string
		task        Task
		expectedErr bool
	}{
		{
			name:        "Test 1: Valid Task",
			task:        Task{ID: "valid_id", Type: CPU_INTENSIVE, CreatedAt: time.Now(), Timeout: 20 * time.Second, Status: "pending"},
			expectedErr: false,
		},
		{
			name:        "Test 2: Invalid Task - Missing ID",
			task:        Task{ID: "", Type: CPU_INTENSIVE, CreatedAt: time.Now(), Timeout: 20 * time.Second, Status: "pending"},
			expectedErr: true,
		},
		{
			name:        "Test 3: Invalid Task - Invalid Type",
			task:        Task{ID: "valid_id", Type: "gpu_intensive", CreatedAt: time.Now(), Timeout: 20 * time.Second, Status: "pending"},
			expectedErr: true,
		},
		{
			name:        "Test 4: Invalid Task - CreatedAt is Zero",
			task:        Task{ID: "valid_id", Type: CPU_INTENSIVE, CreatedAt: time.Time{}, Timeout: 20 * time.Second, Status: "pending"},
			expectedErr: true,
		},
		{
			name:        "Test 5: Invalid Task - Timeout is Zero",
			task:        Task{ID: "valid_id", Type: CPU_INTENSIVE, CreatedAt: time.Now(), Timeout: time.Duration(0), Status: "pending"},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.task.IsValid()
			if !tc.expectedErr && err != nil {
				t.Errorf("Test %s expected no error but received %v", tc.name, err)
			}
			if tc.expectedErr && err == nil {
				t.Errorf("Test %s expected error but received none", tc.name)
			}
		})
	}

}

func TestTask_JSONSerialization_RoundTrip(t *testing.T) {
	originalPayload := map[string]interface{}{"name": "Marine", "race": "Terran", "health": 45, "damage": 6, "is_upgraded": true}

	task, err := NewTask(CPU_INTENSIVE, originalPayload)
	if err != nil {
		t.Fatalf("Failed to creating task: %v", err)
	}

	var unmarshaled map[string]interface{}

	err = task.UnmarshalPayload(&unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal payload: %v", err)
	}

	if unmarshaled["name"] != originalPayload["name"] {
		t.Errorf("Name mismatch: expected %v, got %v", originalPayload["name"], unmarshaled["name"])
	}

	if unmarshaled["health"] != float64(45) {
		t.Errorf("Health mismatch: expected 45, got %v", unmarshaled["health"])
	}
}

func TestTask_IsExpired(t *testing.T) {
	task, _ := NewTask(CPU_INTENSIVE, map[string]string{"test": "data"})
	task.Timeout = 1 * time.Millisecond

	time.Sleep(2 * time.Millisecond)

	if !task.IsExpired() {
		t.Error("Task should be expired but isn't")
	}
}

func TestTaskResult_SetError(t *testing.T) {
	result := NewTaskResult("test-123", "worker-1", COMPLETED)

	testErr := fmt.Errorf("reactor meltdown")
	result.SetError(testErr)

	if result.Status != FAILED {
		t.Errorf("Expected status FAILED, got %v", result.Status)
	}
}
