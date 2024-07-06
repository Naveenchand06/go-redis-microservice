package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	app := &App{
		router: loadRoutes(),
		rdb:    redis.NewClient(&redis.Options{}),
	}

	return app
}

func (a *App) Start(ctx context.Context) error {

	// * Check for connection state with REDIS
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %v", err)
	}

	// * Closing DB connection
	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis connection ", err)
		}
	}()

	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	fmt.Println("Starting Server...")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %v", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}
