package spider

import "context"

type Task interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}
