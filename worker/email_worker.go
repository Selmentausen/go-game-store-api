package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"game-store-api/database"
	"time"
)

func StartEmailWorker() {
	fmt.Println("Email Worker started...")

	for {
		ctx := context.Background()
		result, err := database.RedisClient.BLPop(ctx, 0, "send_email_queue").Result()

		if err != nil {
			fmt.Println("Worker error:", err)
			continue
		}

		taskJSON := result[1]

		var task map[string]string
		json.Unmarshal([]byte(taskJSON), &task)
		fmt.Printf("Processing email for: %s\n", task["email"])

		time.Sleep(3 * time.Second)
		fmt.Printf("Email sent to %s!\n", task["email"])
	}
}
