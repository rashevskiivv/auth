package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	env "tax-auth/internal"
	"tax-auth/internal/handler"
	handlerAuth "tax-auth/internal/handler/auth"
	handlerUser "tax-auth/internal/handler/user"
	"tax-auth/internal/repository"
	repositoryAuth "tax-auth/internal/repository/auth"
	repositoryUser "tax-auth/internal/repository/user"

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
		pg.Close()
	}()

	router = registerHandlers(router, pg)

	// Running
	err = router.Run(fmt.Sprintf(":%v", strconv.Itoa(appPort)))
	if err != nil {
		log.Fatalf("got error while running: %v", err)
	}
}

func registerHandlers(router *gin.Engine, pg *repository.Postgres) *gin.Engine {
	// Repo
	userRepo := repositoryUser.NewUserRepo(*pg)
	authRepo := repositoryAuth.NewAuthRepo(*pg)
	// Handler
	userHandler := handlerUser.NewUserHandler(userRepo)
	authHandler := handlerAuth.NewAuthHandler(authRepo)

	// Routing
	router.NoRoute(handler.NotFound)
	router.GET("/_hc", handler.HealthCheck)
	// Auth
	router.POST("register", authHandler.RegisterUserHandle)
	router.POST("login", authHandler.AuthenticateUserHandle)
	// User
	router.POST("users", userHandler.InsertUserHandle)
	router.GET("users", userHandler.ReadUsersHandle)

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
	return pg, nil
}
