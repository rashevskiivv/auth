package env

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	envAppPort          = "APP_PORT"
	envPostgresDriver   = "POSTGRES_DRIVER"
	envPostgresUser     = "POSTGRES_USER"
	envPostgresPassword = "POSTGRES_PASSWORD"
	envPostgresHost     = "POSTGRES_HOST"
	envPostgresPort     = "POSTGRES_PORT"
	envPostgresDB       = "POSTGRES_DB"
)

func init() {
	if err := godotenv.Load("deployment/.env"); err != nil {
		log.Fatal("can not find .env file: ", err)
	}
}

func GetAppPortEnv() (int, error) {
	portStr := os.Getenv(envAppPort)
	if portStr == "" {
		return 0, errors.New(fmt.Sprintf("can not found: %v", envAppPort))
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("can not convert to integer: %v", envAppPort))
	}
	return port, nil
}

func GetDBUrlEnv() (string, error) {
	dbDriver := os.Getenv(envPostgresDriver)
	if dbDriver == "" {
		return "", fmt.Errorf("can not found: %v", envPostgresDriver)
	}
	dbUser := os.Getenv(envPostgresUser)
	if dbUser == "" {
		return "", fmt.Errorf("can not found: %v", envPostgresUser)
	}
	dbPassword := os.Getenv(envPostgresPassword)
	if dbPassword == "" {
		return "", fmt.Errorf("can not found: %v", envPostgresPassword)
	}
	dbHost := os.Getenv(envPostgresHost)
	if dbHost == "" {
		return "", fmt.Errorf("can not found: %v", envPostgresHost)
	}
	dbPort := os.Getenv(envPostgresPort)
	if dbPort == "" {
		return "", fmt.Errorf("can not found: %v", envPostgresPort)
	}
	dbName := os.Getenv(envPostgresDB)
	if dbName == "" {
		return "", fmt.Errorf("can not found: %v", envPostgresDB)
	}
	return fmt.Sprintf("%v://%v:%v@%v:%v/%v", dbDriver, dbUser, dbPassword, dbHost, dbPort, dbName), nil
}
