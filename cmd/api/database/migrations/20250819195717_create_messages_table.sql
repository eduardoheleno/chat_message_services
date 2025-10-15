-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
	id bigint(20) unsigned auto_increment NOT NULL,
	user_id bigint(20) unsigned NOT NULL,
	chat_id bigint(20) unsigned NOT NULL,
	content longblob NOT NULL,
	nonce longblob NOT NULL,
	created_at datetime(3) NULL,
	CONSTRAINT messages_pk PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
-- +goose StatementEnd
