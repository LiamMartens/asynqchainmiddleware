package asynqchainmiddleware

import (
	"context"
	"log"
	"testing"

	"github.com/hibiken/asynq"
)

func TestMiddleware(t *testing.T) {
	redis_config := asynq.RedisClientOpt{
		Addr: "127.0.0.1:6379",
	}
	client := asynq.NewClient(redis_config)
	server := asynq.NewServer(redis_config, asynq.Config{})

	task2_chan := make(chan string)

	mux := asynq.NewServeMux()
	mux.Use(ChainTasksMiddlewareFactory(client))
	mux.HandleFunc("task1", func(ctx context.Context, t *asynq.Task) error {
		task := asynq.NewTask("task2", nil)
		return &ErrChainTask{
			Task: task,
		}
	})
	mux.HandleFunc("task2", func(ctx context.Context, t *asynq.Task) error {
		task2_chan <- "OK"
		return nil
	})

	go func() {
		if err := server.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	task1 := asynq.NewTask("task1", nil)
	client.Enqueue(task1)

	result := <-task2_chan
	if result != "OK" {
		t.Fatalf("Did not receive %q from task2 channel", "OK")
	}

	server.Stop()
}
