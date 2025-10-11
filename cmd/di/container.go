package di

import (
	"log/slog"
	"net/http"
	"soundtube/internal/domain"
	"soundtube/internal/domain/auth"
	"soundtube/internal/handlers"
	"soundtube/internal/repositories"
	"soundtube/internal/services"
	"soundtube/pkg"
	"soundtube/pkg/config"
	"soundtube/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type Container struct {
	isShuttingDown bool

	Config *config.Config

	TraceProvider *tracesdk.TracerProvider
	Tracer        trace.Tracer
	Logger        *pkg.CustomLogger

	Engine *gin.Engine
	Redis  *redis.Client
	Cache  domain.ICache

	Server *http.Server

	RateLimiter *pkg.RateLimiter

	TokenBlackList auth.ITokenBlacklist

	Repository *repositories.RepositoryAdapter

	RegisterHandler  *handlers.RegisterHandler
	LoginHandler     *handlers.LoginHandler
	VerifyHandler    *handlers.EmailHandler
	SoundHandler     *handlers.SoundHandler
	CommentHandler   *handlers.CommentHandler
	UploadHandler    *handlers.UploadHandler
	ReactionsHandler *handlers.ReactionHandler

	Email           *services.EmailService
	RegisterService *services.RegisterService
	LoginService    *services.LoginService
	SoundService    *services.SoundService
	ReactionService *services.ReactionService
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

	c.initRateLimiter()
	c.initRedis()

	if err = c.initRepositories(); err != nil {
		return err
	}

	c.initServices()
	c.initHandlers()

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
	c.Repository, err = repositories.NewRepositoryAdapter(&c.Config.Database, &c.Config.DatabaseConnections, c.Logger)
	if err != nil {
		return err
	}

	c.TokenBlackList = repositories.NewTokenBlacklist(c.Redis, c.Logger)
	c.Cache = repositories.NewRedisCache(c.Redis)

	return nil
}

func (c *Container) initServices() {
	c.Email = services.NewEmailService(c.Repository.UserRepository, c.Config.Server.Host+c.Config.Server.Port, &c.Config.Email, c.Logger)
	c.RegisterService = services.NewRegisterService(c.Repository, c.Email, c.Logger)
	c.LoginService = services.NewLoginService(c.Config.Token, c.Repository.UserRepository, c.TokenBlackList, c.Logger)
	c.SoundService = services.NewSoundService(c.Repository.SoundRepository, c.Repository.UserRepository, c.Logger)
	c.ReactionService = services.NewRactionService(c.Repository.SoundReactionRepository, c.Repository.SoundPartisipantsRepository, c.Cache, c.Logger)
}

func (c *Container) initHandlers() {
	c.RegisterHandler = handlers.NewRegisterHandler(c.RegisterService, c.Logger)
	c.LoginHandler = handlers.NewLoginHandler(c.LoginService, c.Logger)
	c.SoundHandler = handlers.NewSoundHandler(c.SoundService, c.Logger)
	c.VerifyHandler = handlers.NewEmailHandler(c.Email, c.Logger)
	c.UploadHandler = handlers.NewUploadHandler(c.SoundService, c.Logger)
	c.ReactionsHandler = handlers.NewReactionHandler(c.ReactionService, c.Logger)
}

func (c *Container) initGinEngine() {
	c.Engine = gin.Default()

	c.Engine.Use(middleware.SecurityMiddleware())
	c.Engine.Use(middleware.RequsetIDMiddleware())
	c.Engine.Use(middleware.RateLimiterMiddleware(c.RateLimiter))

	c.Engine.Static("/static", "../../static")
	c.Engine.LoadHTMLGlob("../../static/*.html")

	var api = c.Engine.Group("/api")
	{
		var auth = api.Group("/auth")
		{
			auth.POST("/register", c.RegisterHandler.Register)
			auth.POST("/login", c.LoginHandler.Login)
			auth.POST("/logout", c.LoginHandler.Logout)
			auth.GET("/verify-email", c.VerifyHandler.VerifyEmail)
		}

		var authRequered = api.Group("")
		authRequered.Use(middleware.AuthMiddleware(c.LoginService, c.Logger))

		var sounds = authRequered.Group("/sounds")
		{
			sounds.GET("/", c.SoundHandler.GetSounds)
			sounds.POST("/", c.SoundHandler.CreateSound)
			sounds.POST("/upload", c.UploadHandler.UploadSoundFile)
			sounds.PATCH("/:id", c.SoundHandler.UpdateSound)
			sounds.DELETE("/:id", c.SoundHandler.DeleteSound)

			sounds.GET("/:id/comments", c.CommentHandler.GetComments)
			sounds.POST("/:id/comments", c.CommentHandler.CreateComment)

			sounds.PUT("/:id/reactions", c.ReactionsHandler.SetReactionSound)
			sounds.DELETE("/:id/reactions", c.ReactionsHandler.DeleteReactionSound)
			sounds.GET("/:id/reactions", c.ReactionsHandler.GetReactionSound)
		}

		var comments = authRequered.Group("/comments")
		{
			comments.PATCH("/:id", c.CommentHandler.UpdateComment)
			comments.DELETE("/:id", c.CommentHandler.DeleteComment)

			comments.PUT("/:id/reactions", c.ReactionsHandler.SetReactionComment)
			comments.DELETE("/:id/reactions", c.ReactionsHandler.DeleteReactionComment)
			comments.GET("/:id/reactions", c.ReactionsHandler.GetReactionComment)
		}
	}

	c.Engine.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
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
		Handler:      c.Engine,
		ReadTimeout:  time.Duration(c.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(c.Config.Server.IdleTimeout) * time.Second,
	}
}

func (c *Container) initRateLimiter() {
	c.RateLimiter = pkg.NewRateLimiter(&c.Config.RateLimiter)
}

func (c *Container) initTraycing() error {
	if !c.Config.Traycing.Enabled {
		return nil
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(c.Config.Traycing.Endpoint)))
	if err != nil {
		return err
	}

	c.TraceProvider = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(c.Config.Traycing.ServiceName),
			attribute.String("environment", c.Config.Environment.Current),
		)),
	)

	otel.SetTracerProvider(c.TraceProvider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	c.Tracer = c.TraceProvider.Tracer("app")

	c.Logger.SetTracer(c.Tracer)

	c.Logger.Info("Initializing tracing",
		"service_name", c.Config.Traycing.ServiceName,
		"endpoint", c.Config.Traycing.Endpoint,
	)

	return nil
}

func (c *Container) initHealthCheck() {
	c.Engine.GET("/health", func(ctx *gin.Context) {
		health := map[string]string{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}

		if err := c.Repository.HealthCheck(ctx.Request.Context()); err != nil {
			health["database"] = "unhealthy"
			health["status"] = "degraded"
			ctx.JSON(http.StatusInternalServerError, health)
			return
		}

		if err := c.Redis.Ping().Err(); err != nil {
			health["redis"] = "unhealthy"
			health["status"] = "degraded"
			ctx.JSON(http.StatusInternalServerError, health)
			return
		}

		health["database"] = "healthy"
		health["redis"] = "healthy"
		ctx.JSON(http.StatusOK, health)
	})

	c.Engine.GET("/ready", func(ctx *gin.Context) {
		if c.isShuttingDown {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "shutting down"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	c.Engine.GET("/live", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "live"})
	})
}

func (c *Container) Close() error {
	c.isShuttingDown = true

	if err := c.Repository.Close(); err != nil {
		return err
	}

	return nil
}
