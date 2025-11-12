package task

import (
	"context"
	"errors"
	"fmt"
)

type getUCDataFunc func() (any, error)

var UCDataOutdatedErr = errors.New("UCDataOutdatedErr")

func EnsureValidUC(getUCData getUCDataFunc) Func {
	return func(ctx context.Context, appendTaskFunc func(Name)) error {
		data, err := getUCData()
		if err != nil {
			if errors.Is(err, UCDataOutdatedErr) {
				appendTaskFunc("ReRequestUCData")
				return nil
			}
			return err
		}

		fmt.Print(data)
		return nil
	}
}

func ReRequestUCData() Func {
	return func(ctx context.Context, appendTaskFunc func(Name)) error {
		return nil
	}
}
