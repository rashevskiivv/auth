package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	env "github.com/rashevskiivv/auth/internal"
	"github.com/rashevskiivv/auth/internal/handler"
	handlerAuth "github.com/rashevskiivv/auth/internal/handler/auth"
	handlerUser "github.com/rashevskiivv/auth/internal/handler/user"
	"github.com/rashevskiivv/auth/internal/repository"
	repositoryAuth "github.com/rashevskiivv/auth/internal/repository/auth"
	repositoryUser "github.com/rashevskiivv/auth/internal/repository/user"
	usecaseAuth "github.com/rashevskiivv/auth/internal/usecase/auth"
	usecaseUser "github.com/rashevskiivv/auth/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	appPort, err := env.GetAppPortEnv()
	if err != nil {
		log.Fatal(err)
	}

	pg, err := getPGInstance()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = pg.Close()
		if err != nil {
			log.Fatalf("pg instance close error: %v", err)
		}
		log.Println("pg instance closed")
	}()

	authHandler, userHandler := createHandlers(pg)
	router = registerHandlers(router, authHandler, userHandler)

	// Running
	err = router.Run(fmt.Sprintf(":%v", strconv.Itoa(appPort)))
	if err != nil {
		log.Fatalf("got error while running: %v", err)
	}
}

func createHandlers(pg *repository.Postgres) (handlerAuth.HandlerI, handlerUser.HandlerI) {
	// Repo
	userRepo := repositoryUser.NewUserRepo(*pg)
	authRepo := repositoryAuth.NewAuthRepo(*pg)
	// UseCase
	authUC := usecaseAuth.NewAuthUseCase(authRepo, userRepo)
	userUC := usecaseUser.NewUserUseCase(userRepo)
	// Handler
	authHandler := handlerAuth.NewAuthHandler(authUC)
	userHandler := handlerUser.NewUserHandler(userUC)

	log.Println("handlers created")

	return authHandler, userHandler
}

func registerHandlers(router *gin.Engine, authHandler handlerAuth.HandlerI, userHandler handlerUser.HandlerI) *gin.Engine {
	// Routing
	router.NoRoute(handler.NotFound)
	router.GET("/_hc", handler.HealthCheck)
	// Auth
	router.POST("register", authHandler.RegisterUserHandle)
	router.POST("login", authHandler.AuthenticateUserHandle)
	router.GET("check", authHandler.CheckTokenHandle)
	// User
	group := router.Group("users")
	group.Use(handler.TokenAuthMiddleware(authHandler))
	group.POST("", userHandler.UpsertUserHandle)
	group.GET("", userHandler.ReadUsersHandle)
	group.DELETE("", userHandler.DeleteUsersHandle)

	log.Println("handlers registered")

	return router
}

func getPGInstance() (*repository.Postgres, error) {
	url, err := env.GetDBUrlEnv()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	pg, err := repository.NewPG(ctx, url)
	if err != nil {
		return nil, err
	}
	log.Println("pg instance created")
	return pg, nil
}
