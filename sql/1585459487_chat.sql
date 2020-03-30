BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;

CREATE TABLE chat_messages(
  id SERIAL NOT NULL,

  sender_id INTEGER NOT NULL,
  recipient_id INTEGER NOT NULL,
  message_content JSONB NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),

  PRIMARY KEY (id),
  FOREIGN KEY (sender_id) REFERENCES users(id),
  FOREIGN KEY (recipient_id) REFERENCES users(id)

);

CREATE INDEX chat_messages__idx ON chat_messages(sender_id, recipient_id);

END TRANSACTION;