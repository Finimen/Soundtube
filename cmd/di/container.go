package di

import (
	"log/slog"
	"net/http"
	"soundtube/internal/domain/auth"
	"soundtube/internal/handlers"
	"soundtube/internal/repositories"
	"soundtube/internal/services"
	"soundtube/pkg"
	"soundtube/pkg/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Container struct {
	Config *config.Config

	Logger *pkg.CustomLogger

	Engine *gin.Engine
	Redis  *redis.Client

	Server *http.Server

	Repository *repositories.RepositoryAdapter

	RegisterHandler *handlers.RegisterHandler
	LoginHandler    *handlers.LoginHandler

	Email           auth.IEmailSener
	RegisterService *services.RegisterService
	LoginService    *services.LoginService
}

func NewContainer() (*Container, error) {
	var container = Container{}

	if err := container.init(); err != nil {
		return nil, err
	}

	return &container, nil
}

func (c *Container) init() error {
	if err := c.initCore(); err != nil {
		return err
	}

	c.initProdFeatures()

	return nil
}

func (c *Container) initCore() error {
	var err error
	c.Config, err = config.LoadConfig()
	if err != nil {
		return err
	}

	c.Logger = pkg.NewLogger(slog.Default(), c.Config.Traycing.Enabled)

	if err = c.initRepositories(); err != nil {
		return err
	}

	c.initServices()
	c.initHandlers()

	c.initRedis()
	c.initGinEngine()
	c.initServer()

	c.Logger.Info("core initialization was successful")
	return nil
}

func (c *Container) initProdFeatures() {
	c.initHealthCheck()
	c.initTraycing()

	c.Logger.Info("prod features initialization were successful")
}

func (c *Container) initRepositories() error {
	var err error
	c.Repository, err = repositories.NewRepositoryAdapter(&c.Config.Repository, c.Logger)
	if err != nil {
		return err
	}

	return nil
}

func (c *Container) initServices() {
	c.Email = services.NewEmailService(nil, c.Server.Addr, c.Config.Email.From, c.Logger)
	c.RegisterService = services.NewRegisterService(c.Repository, c.Email, c.Logger)
	c.LoginService = services.NewLoginService(c.Config.Token, c.Repository.UserRepository, c.Repository.TokenBlacklist, c.Logger)
}

func (c *Container) initHandlers() {
	c.RegisterHandler = handlers.NewRegisterHandler(c.RegisterService, c.Logger)
	c.LoginHandler = handlers.NewLoginHandler(c.LoginService, c.Logger)
}

func (c *Container) initGinEngine() {
	c.Engine = gin.Default()

	var api = c.Engine.Group("/api", nil)
	{
		var auth = api.Group("/auth", nil)
		{
			auth.POST("/register", c.RegisterHandler.Register)
			auth.POST("/login", c.LoginHandler.Login)
			auth.POST("/logout", c.LoginHandler.Logout)
		}
	}
}

func (c *Container) initRedis() {
	c.Redis = redis.NewClient(&redis.Options{
		Addr:     c.Config.Redis.Addr,
		Password: c.Config.Redis.Password,
		DB:       c.Config.Redis.DB,
	})
}

func (c *Container) initServer() {
	c.Server = &http.Server{
		Addr:         c.Config.Server.Port,
		ReadTimeout:  time.Duration(c.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(c.Config.Server.IdleTimeout) * time.Second,
	}
}

func (c *Container) initTraycing() {

}

func (c *Container) initHealthCheck() {

}

func (c *Container) Close() {

}
