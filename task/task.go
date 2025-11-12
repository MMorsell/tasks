package task

import "context"

type Func func(ctx context.Context, appendTaskFunc func(Name)) error

type Name string
