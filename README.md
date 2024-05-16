# asynqchainmiddleware
This middleware can be used to create a task chain workflow in with [asynq](github.com/hibiken/asynq)  

## Usage
```go
import (
	"github.com/hibiken/asynq"
  "asynqchainmiddleware"
)

redis_config := asynq.RedisClientOpt{
  Addr: "127.0.0.1:6379",
}
client := asynq.NewClient(redis_config)
server := asynq.NewServer(redis_config, asynq.Config{})

mux := asynq.NewServeMux()
mux.Use(ChainTasksMiddlewareFactory(client))
mux.HandleFunc("task1", func(ctx context.Context, t *asynq.Task) error {
  task := asynq.NewTask("task2", nil)
  return &ErrChainTask{
    Task: task,
  }
})
```