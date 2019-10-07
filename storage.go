package main

import (
	"errors"
	"log"
	"strconv"
)

type MessageStorage struct {
	userList map[int64]*user
	data     map[user][]string
}

type user struct {
	chatID int64
	time   string
}

func NewStorage() *MessageStorage {
	storage := MessageStorage{
		userList: make(map[int64]*user, 0),
		data:     make(map[user][]string),
	}
	storage.start()
	return &storage
}

func (s *MessageStorage) AddUser(chatID int64, time string) {
	user := user{
		chatID: chatID,
		time:   time,
	}
	s.userList[chatID] = &user
}

func (s *MessageStorage) AddMessage(chatID int64, message string) {
	user, err := s.findUser(chatID)
	if err != nil {
		return
	}
	s.data[*user] = append(s.data[*user], message)
	log.Println("Chat: " + strconv.FormatInt(chatID, 10) + " message: " + message)
}

func (s *MessageStorage) UpdateUserSettings(chatID int64, time string) {
	usr, err := s.findUser(chatID)
	if err != nil {
		usr = &user{
			chatID: chatID,
		}
		s.userList[chatID] = usr
	}
	usr.time = time
	log.Println("Set " + time + " for chat " + strconv.FormatInt(chatID, 10))
}

func (s *MessageStorage) isUserTunned(chatID int64) bool {
	_, isExist := s.userList[chatID]
	return isExist
}

func (s *MessageStorage) start() {
	//start cron
}

func (s *MessageStorage) findUser(chatID int64) (*user, error) {
	user, isExist := s.userList[chatID]
	if !isExist {
		return nil, errors.New("user by chatID " + strconv.FormatInt(chatID, 10) + " not found")
	}
	return user, nil
}
