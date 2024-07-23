-- +goose Up
CREATE TABLE chats
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE chat_users
(
    chat_id BIGINT NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE messages
(
    id        BIGSERIAL PRIMARY KEY,
    chat_id   BIGINT      NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
    from_user BIGINT      NOT NULL,
    text      TEXT        NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- +goose Down
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chat_users;
DROP TABLE IF EXISTS chats;

