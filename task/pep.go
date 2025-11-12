package task

import (
	"context"
)

func RequireNonPEP() Func {
	return func(ctx context.Context, appendTaskFunc func(Name)) error {
		return nil
	}
}
