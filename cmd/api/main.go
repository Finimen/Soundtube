package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"soundtube/cmd/di"
	"syscall"
	"time"
)

// @title SoundTube API
// @version 1.0
// @description Audio sharing and streaming platform API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

func getContainer() *di.Container {
	var container, err = di.NewContainer()
	if err != nil {
		panic(err)
	}

	return container
}

// main is the entry point of the SoundTube application
// @Summary SoundTube Server
// @Description Main server for SoundTube audio platform handling authentication, audio uploads, and social features
func main() {
	var container = getContainer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	container.Logger.Info("initialization completed")

	go func() {
		info := fmt.Sprintf("Server is starting on http://localhost%s\n", container.Server.Addr)
		container.Logger.Info(info)
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
		container.Logger.Error("Server forced to shutdown", err)
	}

	container.Logger.Info("OK ")
}
