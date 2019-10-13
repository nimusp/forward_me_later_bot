CREATE TABLE IF NOT EXISTS chats(
    id INT NOT NULL,
    time_to_forward TIMESTAMP NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS messages(
    id SERIAL,
    chat_id INT NOT NULL,
    message_id INT NOT NULL,
    added_at TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(chat_id) REFERENCES chats(id)
);