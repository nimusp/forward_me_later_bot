package main

import (
	"errors"
	"log"
	"strconv"
	"time"
)

type MessageStorage struct {
	userList map[int64]*user
	data     map[user][]int
}

type user struct {
	chatID int64
	time   time.Time
}

func NewStorage() *MessageStorage {
	return &MessageStorage{
		userList: make(map[int64]*user, 0),
		data:     make(map[user][]int),
	}
}

func (s *MessageStorage) AddUser(chatID int64, stringTime string) {
	castedTime := parseTime(stringTime)

	user := user{
		chatID: chatID,
		time:   castedTime,
	}
	s.userList[chatID] = &user
}

func (s *MessageStorage) AddMessage(chatID int64, messageID int) {
	user, err := s.findUser(chatID)
	if err != nil {
		return
	}
	s.data[*user] = append(s.data[*user], messageID)
	log.Println("Chat: " + strconv.FormatInt(chatID, 10) + " message: " + strconv.Itoa(messageID))
}

func (s *MessageStorage) UpdateUserSettings(chatID int64, time string) {
	usr, err := s.findUser(chatID)
	if err != nil {
		usr = &user{
			chatID: chatID,
		}
		s.userList[chatID] = usr
	}
	usr.time = parseTime(time)
	log.Println("Set " + time + " for chat " + strconv.FormatInt(chatID, 10))
}

func (s *MessageStorage) isUserTunned(chatID int64) bool {
	_, isExist := s.userList[chatID]
	return isExist
}

func (s *MessageStorage) findUser(chatID int64) (*user, error) {
	user, isExist := s.userList[chatID]
	if !isExist {
		return nil, errors.New("user by chatID " + strconv.FormatInt(chatID, 10) + " not found")
	}
	return user, nil
}

func (s *MessageStorage) DeleteMessageForChat(chatID int64) {
	user, err := s.findUser(chatID)
	if err != nil {
		log.Println(err)
		return
	}
	s.data[*user] = make([]int, 0)
}

func parseTime(stringTime string) time.Time {
	timeFormat := "15:04"
	castedTime, err := time.Parse(timeFormat, stringTime)
	if err != nil {
		log.Println(err)
	}

	now := time.Now()
	currentDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	correctTime := currentDay.Add(
		time.Hour*time.Duration(castedTime.Hour()) + time.Minute*time.Duration(castedTime.Minute()),
	)

	return correctTime
}

func (s *MessageStorage) getAllSheduledJobs() (map[int64][]int, map[int64]time.Time) {
	chatToMessage := make(map[int64][]int, len(s.data))
	chatToTime := make(map[int64]time.Time, len(s.data))

	for user, messageList := range s.data {
		chatToMessage[user.chatID] = messageList
		chatToTime[user.chatID] = user.time
	}

	return chatToMessage, chatToTime
}
