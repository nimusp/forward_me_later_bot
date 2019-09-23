package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

const telegramBotTokenKey = "API_TOKEN"

func main() {
	botToken := os.Getenv(telegramBotTokenKey)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Println(err)
	}

	bot.Debug = true

	log.Println(bot.Self.UserName)
}
