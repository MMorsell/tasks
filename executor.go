package main

import (
	"context"
	"slices"
	"time"

	"github.com/MMorsell/tasks/internal"
	"github.com/MMorsell/tasks/task"
	"github.com/google/uuid"
)

type storeTaskOutcomeFunc func(uuid.UUID, time.Time, time.Time, task.Name, error) error

func Execute(
	ctx context.Context,
	taskRegistry map[task.Name]task.Func,
	tasks []task.Func,
	storeOutcome storeTaskOutcomeFunc,
) error {
	slices.Reverse(tasks)

	stack := internal.NewStack(tasks)

	for {
		if stack.IsEmpty() {
			return nil
		}

		f, _ := stack.Top()
		var (
			taskID = uuid.New()
			start  = time.Now()
		)
		err := f(ctx, Append(stack, taskRegistry))
		stop := time.Now()
		if err := storeOutcome(taskID, start, stop, "", err); err != nil {
			// todo(mm): how to handle failed to store task outcome....
		}
		if err != nil {
			return err
		}
	}
}

func Append(stack *internal.Stack[task.Func], taskRegistry map[task.Name]task.Func) func(task.Name) {
	return func(tn task.Name) {
		stack.Push(taskRegistry[tn])
	}
}
