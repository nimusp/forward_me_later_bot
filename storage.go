package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"database/sql"
	_ "github.com/lib/pq"
)

type MessageStorage struct {
	db       *sql.DB
	userList map[int64]bool
}

type Message struct {
	MessageID   int
	AddedAtTime time.Time
}

func NewStorage(dbLogin, dbPassword, dbName string) *MessageStorage {
	dataSourseName := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbLogin, dbPassword, dbName)
	db, err := sql.Open("postgres", dataSourseName)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &MessageStorage{
		db:       db,
		userList: make(map[int64]bool, 0),
	}
}

func (s *MessageStorage) AddMessage(chatID int64, messageID int) {
	if !s.isUserTunned(chatID) {
		return
	}

	stmt, err := s.db.Prepare(
		`INSERT INTO messages (chat_id, message_id, added_at)
		 VALUES ($1, $2, $3)`,
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(chatID, messageID, time.Now())
	if err != nil {
		log.Println(err)
	}

	log.Println("Chat: " + strconv.FormatInt(chatID, 10) + " message: " + strconv.Itoa(messageID))
}

func (s *MessageStorage) UpdateUserSettings(chatID int64, time string) {
	parsedTime := parseTime(time)

	stmt, err := s.db.Prepare(
		`INSERT INTO chats (id, time_to_forward)
		 VALUES ($1, $2)
		 ON CONFLICT (id) DO UPDATE SET time_to_forward=$2`,
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(chatID, parsedTime)
	if err != nil {
		log.Println(err)
		return
	}
	s.userList[chatID] = true
	log.Println("Set " + time + " for chat " + strconv.FormatInt(chatID, 10))
}

func (s *MessageStorage) isUserTunned(chatID int64) bool {
	_, isExist := s.userList[chatID]
	if !isExist {
		stmt, err := s.db.Prepare(
			`SELECT COUNT(*)
			 FROM chats
			 WHERE chats.id = $1`,
		)
		if err != nil {
			log.Println(err)
			return false
		}
		row := stmt.QueryRow(chatID)
		var count int
		row.Scan(&count)
		if count > 0 {
			s.userList[chatID] = true
		}

		isExist = count > 0
	}

	return isExist
}

func (s *MessageStorage) DeleteMessageByID(messageID int) {
	stmt, err := s.db.Prepare(
		`DELETE FROM messages
		 WHERE messages.message_id = $1`,
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(messageID)
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

func (s *MessageStorage) getAllSheduledJobs() (map[int64][]Message, map[int64]time.Time) {
	size := s.db.QueryRow(
		`SELECT COUNT(messages.message_id)
		 FROM messages
		 GROUP BY messages.message_id
		 ORDER BY DESC
		 LIMIT 1`,
	)
	var dataSize int
	size.Scan(&dataSize)

	chatToMessage := make(map[int64][]Message, dataSize)
	chatToTime := make(map[int64]time.Time, dataSize)

	rows, err := s.db.Query(
		`SELECT messages.chat_id, messages.message_id, messages.added_at
		 FROM messages`,
	)
	if err != nil {
		log.Println(err)
		first := make(map[int64][]Message)
		second := make(map[int64]time.Time)
		return first, second
	}
	defer rows.Close()

	for rows.Next() {
		var chatID int64
		var messageID int
		var addedAt time.Time
		rows.Scan(&chatID, &messageID, &addedAt)
		message := Message{MessageID: messageID, AddedAtTime: addedAt}
		chatToMessage[chatID] = append(chatToMessage[chatID], message)
	}

	userRows, err := s.db.Query(
		`SELECT chats.id, chats.time_to_forward
		 FROM chats`,
	)
	if err != nil {
		log.Println(err)
		first := make(map[int64][]Message)
		second := make(map[int64]time.Time)
		return first, second
	}
	defer userRows.Close()

	for userRows.Next() {
		var settedTime time.Time
		var chatID int64
		userRows.Scan(&chatID, &settedTime)
		chatToTime[chatID] = settedTime
	}

	return chatToMessage, chatToTime
}
