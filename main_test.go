package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MMorsell/tasks/task"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestExecuteSimple(t *testing.T) {
	var counter int
	mk := func(expect int, label string) task.Func {
		return func(ctx context.Context, appendTask func(task.Name)) error {
			t.Helper()
			if counter != expect {
				t.Fatalf("%s: counter=%d, want=%d", label, counter, expect)
			}
			counter++
			return nil
		}
	}

	tasks := []task.Func{
		mk(0, "task1"),
		mk(1, "task2"),
		mk(2, "task3"),
	}
	Execute(t.Context(), map[task.Name]task.Func{}, tasks, alwaysNilstoreOutcomeMock)
}

func TestExecuteAppending(t *testing.T) {
	var majorCounter int
	var minorCounter int

	majorFunc := func(expect int, label string) task.Func {
		return func(ctx context.Context, appendTask func(task.Name)) error {
			minorCounter = 0
			t.Helper()
			if majorCounter != expect {
				t.Fatalf("%s: majorCounter=%d, want=%d", label, majorCounter, expect)
			}
			majorCounter++

			appendTask("minorTask0")
			return nil
		}
	}

	minorFunc := func(expect int, label string) task.Func {
		return func(ctx context.Context, appendTask func(task.Name)) error {
			t.Helper()
			if minorCounter != expect {
				t.Fatalf("%s: minorCounter=%d, want=%d", label, majorCounter, expect)
			}
			minorCounter++

			if minorCounter == 1 {
				appendTask("minorTask1")
			}

			return nil
		}
	}

	taskRegistry := map[task.Name]task.Func{
		"minorTask0": minorFunc(0, "minor_0"),
		"minorTask1": minorFunc(1, "minor_1"),
	}

	tasks := []task.Func{
		majorFunc(0, "task1"),
		majorFunc(1, "task2"),
	}
	Execute(t.Context(), taskRegistry, tasks, alwaysNilstoreOutcomeMock)
}

var alwaysNilstoreOutcomeMock = func(u uuid.UUID, t1, t2 time.Time, tn task.Name, err error) error { return nil }

func TestGetUCTask(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		registry map[task.Name]task.Func
		err      error
	}{
		"OK_UC_DATA": {
			registry: map[task.Name]task.Func{
				"GetUC": task.EnsureValidUC(func() (any, error) { return "", nil }),
			},
		},
		"OUTDATED_UC_DATA": {
			registry: map[task.Name]task.Func{
				"GetUC":           task.EnsureValidUC(func() (any, error) { return "", task.UCDataOutdatedErr }),
				"ReRequestUCData": task.ReRequestUCData(),
			},
		},
		"NOK_UC_DATA": {
			registry: map[task.Name]task.Func{
				"GetUC":           task.EnsureValidUC(func() (any, error) { return "", errors.New("a") }),
				"ReRequestUCData": task.ReRequestUCData(),
			},
			err: errors.New("aoeu"),
		},
	}

	var registry = map[task.Name]task.Func{
		"GetUC":           task.EnsureValidUC(func() (any, error) { return "", nil }),
		"ReRequestUCData": task.ReRequestUCData(),
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tasks := []task.Func{
				test.registry["GetUC"],
			}
			err := Execute(t.Context(), registry, tasks, alwaysNilstoreOutcomeMock)
			if test.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
