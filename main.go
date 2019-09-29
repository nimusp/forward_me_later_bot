package main

import (
	"os"
)

const telegramBotTokenKey = "API_TOKEN"

func main() {
	token := os.Getenv(telegramBotTokenKey)
	storage := NewStorage()
	messageHandler := NewHandler(token, storage)
	messageHandler.Start()
}
