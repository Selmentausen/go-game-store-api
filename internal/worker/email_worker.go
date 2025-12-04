package worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func StartEmailWorker(redisClient *redis.Client) {
	slog.Info("Email Worker started")

	for {
		ctx := context.Background()
		result, err := redisClient.BLPop(ctx, 0, "send_email_queue").Result()

		if err != nil {
			slog.Error("Worker failed to pop task", "error", err)
			continue
		}

		taskJSON := result[1]

		var task map[string]string
		if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
			slog.Error("Worker failed to parse task", "error", err, "raw_data", taskJSON)
			continue
		}
		slog.Info("Processing email",
			"email", task["email"],
			"type", task["type"],
		)

		// Simulate work
		time.Sleep(2 * time.Second)

		slog.Info("Email sent successfully", "email", task["email"])
	}
}
