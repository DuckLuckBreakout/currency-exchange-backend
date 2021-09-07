package server

import (
	"fmt"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/middlewares"
	"log"

	session_handler "github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/session/handler"
	session_repo "github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/session/repository"
	session_usecase "github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/session/usecase"
	user_handler "github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/user/handler"
	user_repo "github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/user/repository"
	user_usecase "github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/user/usecase"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/configer"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/s3_storage"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Start() {
	// Init config
	configer.InitConfig("configs/app/api_server.yaml")

	// Connect to postgreSql db
	postgreSqlConn, err := sqlx.Open(
		"postgres",
		fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
			configer.AppConfig.Postgresql.User,
			configer.AppConfig.Postgresql.Password,
			configer.AppConfig.Postgresql.DBName,
			configer.AppConfig.Postgresql.Host,
			configer.AppConfig.Postgresql.Port,
			configer.AppConfig.Postgresql.Sslmode,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer postgreSqlConn.Close()
	if err := postgreSqlConn.Ping(); err != nil {
		log.Fatal(err)
	}

	// Connect to redis db #0
	redisConnDB0 := redis.NewClient(&redis.Options{
		Addr:     configer.AppConfig.Redis.Addr,
		Password: configer.AppConfig.Redis.Password,
		DB:       0,
	})
	if redisConnDB0 == nil {
		log.Fatal(err)
	}
	defer redisConnDB0.Close()

	// Connect to redis db #1
	redisConnDB1 := redis.NewClient(&redis.Options{
		Addr:     configer.AppConfig.Redis.Addr,
		Password: configer.AppConfig.Redis.Password,
		DB:       1,
	})
	if redisConnDB1 == nil {
		log.Fatal(err)
	}
	defer redisConnDB1.Close()

	// Init s3
	s3_storage.InitNewConnection(&configer.AppConfig.S3)

	// Create logger
	userRepo := user_repo.NewRepository(postgreSqlConn)
	sessionRepo := session_repo.NewSessionRedisRepository(redisConnDB0, redisConnDB1)

	userUseCase := user_usecase.NewUseCase(userRepo)
	sessionUseCase := session_usecase.NewUseCase(sessionRepo)

	userHandler := user_handler.NewHandler(userUseCase, sessionUseCase)
	sessHandler := session_handler.NewHandler(sessionUseCase)

	mw := middlewares.NewMiddlewareManager(sessionUseCase)

	mainRouter := gin.New()
	mainRouter.Use(gin.Logger())
	mainRouter.Use(gin.Recovery())
	mainRouter.Use(mw.CORSMiddleware())

	// User
	mainRouter.POST("/api/v1/user/auth", userHandler.SignUp)
	mainRouter.GET("/api/v1/user/me/email", userHandler.CheckUniqEmail)
	mainRouter.PUT("/api/v1/user/me/login", mw.TokenAuth(), userHandler.UpdateSelfLogin)
	mainRouter.PUT("/api/v1/user/me/avatar", mw.TokenAuth(), userHandler.UpdateSelfAvatar)
	mainRouter.PUT("/api/v1/user/auth", userHandler.LogIn)
	mainRouter.GET("/api/v1/user/me", mw.TokenAuth(), userHandler.GetSelfProfile)
	mainRouter.DELETE("/api/v1/user/me", mw.TokenAuth(), userHandler.DeleteSelfProfile)

	// Session
	mainRouter.DELETE("/api/v1/session", mw.TokenAuth(), sessHandler.DestroySession)
	mainRouter.PUT("/api/v1/session", sessHandler.RefreshSession)

	mainRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(mainRouter.Run(fmt.Sprintf(
		"%s:%s",
		configer.AppConfig.Server.Host,
		configer.AppConfig.Server.Port,
	)))
}
