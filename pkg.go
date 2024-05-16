package asynqchainmiddleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
)

type ErrChainTask struct {
	Task *asynq.Task
}

func (err *ErrChainTask) Error() string {
	return fmt.Sprintf("chaining task %q", err.Task.Type())
}

func ChainTasksMiddlewareFactory(client *asynq.Client) func(asynq.Handler) asynq.Handler {
	mw := func(h asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
			err := h.ProcessTask(ctx, t)

			var chain_task_err *ErrChainTask
			if errors.As(err, &chain_task_err) {
				_, err = client.Enqueue(chain_task_err.Task)
				if err != nil {
					return fmt.Errorf("failed to schedule task: %w", err)
				}
				return nil
			}

			if err != nil {
				return err
			}
			return nil
		})
	}
	return mw
}
