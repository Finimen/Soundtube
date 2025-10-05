package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"soundtube/cmd/di"
	"syscall"
	"time"
)

func getContainer() *di.Container {
	var container, err = di.NewContainer()
	if err != nil {
		panic(err)
	}

	return container
}

func main() {
	var container = getContainer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		container.Logger.Info("Server is starting on http://localhost%s\n", container.Server.Addr)
		if err := container.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			container.Logger.Info("listen: http://localhost%s\n", err)
		}
	}()

	<-quit

	container.Logger.Info("Shutting down server")
	container.Close()

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	if err := container.Server.Shutdown(ctx); err != nil {
		container.Logger.Error(ctx, "Server forced to shutdown", err)
	}

	container.Logger.Info("initialization completed")
}
