package main

import (
	_ "github.com/joho/godotenv/autoload"
)

import (
	"context"
	"github.com/HydroProtocol/hydro-box-dex/backend/cli"
	"github.com/HydroProtocol/hydro-box-dex/backend/connection"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/config"
	"github.com/HydroProtocol/hydro-sdk-backend/websocket"
	"os"
	"sync"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := context.WithCancel(context.Background())

	redisClient := connection.NewRedisClient(config.Getenv("HSK_REDIS_URL"))
	redisClient = redisClient.WithContext(ctx)

	go cli.WaitExitSignal(stop)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		websocket.StartConsumer(ctx, &common.RedisQueueConfig{
			Name:   common.HYDRO_WEBSOCKET_MESSAGES_QUEUE_KEY,
			Ctx:    ctx,
			Client: redisClient,
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		websocket.StartSocketServer(ctx)
	}()

	wg.Wait()

	return 0
}
