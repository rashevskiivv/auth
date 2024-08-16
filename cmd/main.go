package cmd

import (
	"log"
	"strconv"
	env "tax-auth/internal"

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

	router = registerHandlers(router)

	// Running
	err = router.Run(":" + strconv.Itoa(appPort))
	if err != nil {
		log.Fatalf("got error while running: %v", err)
	}
}
