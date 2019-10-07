package main

import (
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const timeRegexpPattern = "^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$"

// command list
const (
	startCommand   = "start"
	setTimeCommand = "set_time_to_forward"
	giveItCommand  = "give_it_all_right_now"
)

// command handler message
const (
	setTimeCommandMessage   = "Enter time in which you want to receive all daily messages. \nFormat: HH:mm"
	startCommandMessage     = "Received /start"
	giveItAllCommandMessage = "Received /give_it_all_right_now"
	wrongCommandMessage     = "I don't know that command"
)

type MessageHandler struct {
	bot            *tgbotapi.BotAPI
	storage        *MessageStorage
	chatToSettings map[int64]bool
}

func NewHandler(token string, storage *MessageStorage) *MessageHandler {
	telegramBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	return &MessageHandler{
		bot:            telegramBot,
		storage:        storage,
		chatToSettings: make(map[int64]bool, 0),
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
		isCommand := event.Message.IsCommand()
		if isCommand {
			commandAnswer := h.handleCommandMessage(event, chatID)
			messageToUser.Text = commandAnswer
		}

		isUserTunned := h.storage.isUserTunned(chatID)
		if isUserTunned && !isCommand {
			h.storage.AddMessage(chatID, messageText)
		}

		if timePattern.MatchString(event.Message.Text) && h.chatToSettings[chatID] {
			h.chatToSettings[chatID] = false
			h.storage.UpdateUserSettings(chatID, messageText)
			messageToUser.Text = "Received time: " + event.Message.Text
		}

		if messageToUser.Text != "" {
			if _, err := h.bot.Send(messageToUser); err != nil {
				log.Println(err)
			}
		}
	}
}

func (h *MessageHandler) handleCommandMessage(update tgbotapi.Update, chatID int64) string {
	var answer string

	switch update.Message.Command() {
	case startCommand:
		answer = startCommandMessage
	case setTimeCommand:
		h.chatToSettings[chatID] = true
		answer = setTimeCommandMessage
	case giveItCommand:
		answer = giveItAllCommandMessage
	default:
		answer = wrongCommandMessage
	}

	return answer
}
