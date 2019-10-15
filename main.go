package main

import "os"

const (
	telegramBotTokenKey = "API_TOKEN"
	dbLogin             = "DB_LOGIN"
	dbPassword          = "DB_PASSWORD"
	dbName              = "DB_NAME"
	dbHost              = "DB_HOST"
	dbPort              = "DB_PORT"
)

func main() {
	token := os.Getenv(telegramBotTokenKey)
	login := os.Getenv(dbLogin)
	password := os.Getenv(dbPassword)
	dbName := os.Getenv(dbName)
	dbHost := os.Getenv(dbHost)
	dbPort := os.Getenv(dbPort)
	storage := NewStorage(login, password, dbName, dbHost, dbPort)
	messageHandler := NewHandler(token, storage)
	messageHandler.Start()
}
