package main

import (
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const timeRegexpPattern = "^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$"

type MessageHandler struct {
	bot     *tgbotapi.BotAPI
	storage *MessageStorage
}

func NewHandler(token string, storage *MessageStorage) *MessageHandler {
	telegramBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	return &MessageHandler{
		bot:     telegramBot,
		storage: storage,
	}
}

func (h *MessageHandler) Start() {
	log.Println("Started " + h.bot.Self.UserName)
	timePattern := regexp.MustCompile(timeRegexpPattern)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updateEvents, _ := h.bot.GetUpdatesChan(updateConfig)
	for event := range updateEvents {
		if event.Message == nil {
			continue
		}
		chatID := event.Message.Chat.ID
		messageText := event.Message.Text

		messageToUser := tgbotapi.NewMessage(chatID, "")

		if event.Message.IsCommand() {
			commandAnswer := handleCommandMessage(event)
			messageToUser.Text = commandAnswer
		}

		if timePattern.MatchString(event.Message.Text) {
			h.storage.UpdateUserSettings(chatID, messageText)
			messageToUser.Text = "Received time: " + event.Message.Text
			log.Println("Set time " + messageText + " for chat " + string(chatID))
		} else {
			h.storage.AddMessage(chatID, messageText)
		}

		if _, err := h.bot.Send(messageToUser); err != nil {
			log.Println(err)
		}
	}
}

func handleCommandMessage(update tgbotapi.Update) string {
	var answer string

	switch update.Message.Command() {
	case "start":
		answer = "Received /start"
	case "set_time_to_forward":
		answer = "Enter time in which you want to receive all daily messages. \nFormat: HH:mm"
	case "give_it_all_right_now":
		answer = "Received /give_it_all_right_now"
	default:
		answer = "I don't know that command"
	}

	return answer
}
