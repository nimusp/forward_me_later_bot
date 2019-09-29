package main

type MessageStorage struct {
	data map[int64][]string
}

type user struct {
	chatID int64
	time   string
}

func NewStorage() *MessageStorage {
	storage := MessageStorage{}
	storage.start()
	return &storage
}

func (s *MessageStorage) AddUser(chatID int64, time string) {
	_ = user{
		chatID: chatID,
		time:   time,
	}
	s.data[chatID] = make([]string, 0)
}

func (s *MessageStorage) AddMessage(chatID int64, message string) {
	_, isUserExist := s.data[chatID]
	if !isUserExist {
		return
	}
	s.data[chatID] = append(s.data[chatID], message)
}

func (s *MessageStorage) UpdateUserSettings(chatID int64, time string) {
	//update user time settings
}

func (s *MessageStorage) start() {
	//start cron
}
