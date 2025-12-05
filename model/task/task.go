package task

import (
	"context"
)

type TaskInterface interface {
	PreEvent(ctx context.Context) (err error)
	Run(ctx context.Context) (err error)
	SuccessEvent(ctx context.Context) (err error)
	FailEvent(ctx context.Context) (err error)
	CompleteEvent(ctx context.Context) (err error)
}
