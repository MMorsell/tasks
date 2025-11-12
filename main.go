package main

import (
	"context"
	"time"

	"github.com/MMorsell/tasks/task"
	"github.com/google/uuid"
)

func main() {

	registry := map[task.Name]task.Func{
		"GetUC":           task.EnsureValidUC(func() (any, error) { return "", nil }),
		"ReRequestUCData": task.ReRequestUCData(),
	}

	Run(context.TODO(), registry)
}

func Run(ctx context.Context, registry map[task.Name]task.Func) error {
	tasks := []task.Func{
		task.EnsureValidUC(func() (any, error) { return "", nil }),
		task.RequireNonPEP(),
	}

	return Execute(ctx, registry, tasks, func(u uuid.UUID, t1, t2 time.Time, n task.Name, err error) error { return nil })
}


