package main

import (
	"log"
	"regexp"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const timeRegexpPattern = "^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$"

// command list
const (
	startCommand   = "start"
	setTimeCommand = "set_time_to_forward"
)

// command handler message
const (
	setTimeCommandMessage = "Enter time in which you want to receive all daily messages. \nFormat: HH:mm"
	startCommandMessage   = "Received /start"
	wrongCommandMessage   = "I don't know that command"
)

type messageForwarder struct {
	bot     *tgbotapi.BotAPI
	storage *MessageStorage
}

type MessageHandler struct {
	bot               *tgbotapi.BotAPI
	storage           *MessageStorage
	isNowConfigurable map[int64]bool
	forwarder         *messageForwarder
}

func NewHandler(token string, storage *MessageStorage) *MessageHandler {
	telegramBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	forwarder := &messageForwarder{
		bot:     telegramBot,
		storage: storage,
	}

	return &MessageHandler{
		bot:               telegramBot,
		storage:           storage,
		isNowConfigurable: make(map[int64]bool, 0),
		forwarder:         forwarder,
	}
}

func (h *MessageHandler) Start() {
	log.Println("Started " + h.bot.Self.UserName)
	timePattern := regexp.MustCompile(timeRegexpPattern)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	go h.forwarder.start()

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
		if isUserTunned && !isCommand && !h.isNowConfigurable[chatID] {
			h.storage.AddMessage(chatID, event.Message.MessageID)
			continue
		}

		if timePattern.MatchString(event.Message.Text) && h.isNowConfigurable[chatID] {
			h.isNowConfigurable[chatID] = false
			h.storage.UpdateUserSettings(chatID, messageText)
			messageToUser.Text = "Received time: " + event.Message.Text
		}

		if _, err := h.bot.Send(messageToUser); err != nil {
			log.Println(err)
		}
	}
}

func (h *MessageHandler) handleCommandMessage(update tgbotapi.Update, chatID int64) string {
	var answer string

	switch update.Message.Command() {
	case startCommand:
		answer = startCommandMessage
	case setTimeCommand:
		h.isNowConfigurable[chatID] = true
		answer = setTimeCommandMessage
	default:
		answer = wrongCommandMessage
	}

	return answer
}

func (m *messageForwarder) start() {
	mutex := &sync.Mutex{}
	ticker := time.NewTicker(500 * time.Millisecond)
	updateTicker := time.NewTicker(5 * time.Second)

	chatToMessageList, chatToTime := m.storage.getAllSheduledJobs()

	for {
		select {
		case <-ticker.C:
			for chat, messageList := range chatToMessageList {
				if len(messageList) > 0 {
					chatTime := chatToTime[chat]

					for _, message := range messageList {
						if isReadyToSend(chatTime, message.AddedAtTime) {
							messageToForward := tgbotapi.NewMessage(chat, "received today")
							messageToForward.ReplyToMessageID = message.MessageID
							if _, err := m.bot.Send(messageToForward); err != nil {
								log.Println(err)
							}
							m.storage.DeleteMessageByID(message.MessageID)
						}
					}
					mutex.Lock()
					chatToMessageList, chatToTime = m.storage.getAllSheduledJobs()
					mutex.Unlock()
				}
			}
		case <-updateTicker.C:
			mutex.Lock()
			chatToMessageList, chatToTime = m.storage.getAllSheduledJobs()
			mutex.Unlock()
		default:
		}
	}
}

func isReadyToSend(chatTime, messageTime time.Time) bool {
	currentTime := time.Now()

	currentDay := currentTime.Day()
	currentMonth := currentTime.Month()
	currentYear := currentTime.Year()

	isSetTomorowOrBefore := currentYear > chatTime.Year() ||
		currentYear == chatTime.Year() && currentMonth > chatTime.Month() ||
		currentYear == chatTime.Year() && currentMonth == chatTime.Month() && currentDay > chatTime.Day()

	isSendToday := currentYear == chatTime.Year() && currentMonth == chatTime.Month() && currentDay == chatTime.Day()

	isSendBeforeTrigger := currentTime.Hour() > chatTime.Hour() && messageTime.Hour() < chatTime.Hour() ||
		currentTime.Hour() == chatTime.Hour() && messageTime.Hour() == chatTime.Hour() &&
			currentTime.Minute() >= chatTime.Minute() && messageTime.Minute() <= chatTime.Minute()

	return isSetTomorowOrBefore || isSendToday && isSendBeforeTrigger
}
