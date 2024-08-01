package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/mikhailsoldatkin/chat-server/internal/client/db"
	"github.com/mikhailsoldatkin/chat-server/internal/customerrors"
	"github.com/mikhailsoldatkin/chat-server/internal/repository"
	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
)

const (
	tableChatUsers  = "chat_users"
	tableChats      = "chats"
	tableMessages   = "messages"
	columnID        = "id"
	columnCreatedAt = "created_at"
	columnChatID    = "chat_id"
	columnUserID    = "user_id"
	columnFromUser  = "from_user"
	columnTimestamp = "timestamp"
	columnText      = "text"
	chatEntity      = "chat"
)

type repo struct {
	db db.Client
}

// NewRepository creates a new instance of the chat repository.
func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

// chatExists checks that chat exists.
func (r *repo) chatExists(ctx context.Context, id int64) error {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE %s=$1", tableChats, columnID)
	q := db.Query{
		Name:     "chat_repository.chatExists",
		QueryRaw: query,
	}

	var exists int
	err := r.db.DB().QueryRowContext(ctx, q, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return customerrors.NewNotFoundError(chatEntity, id)
		}
		return err
	}

	return nil
}

// isUserInChat checks if a user is a member of a chat.
func (r *repo) isUserInChat(ctx context.Context, userID int64, chatID int64) error {
	query := fmt.Sprintf(
		"SELECT 1 FROM %s WHERE %s=$1 AND %s=$2",
		tableChatUsers, columnChatID, columnUserID,
	)
	q := db.Query{
		Name:     "chat_repository.isUserInChat",
		QueryRaw: query,
	}

	var exists int
	err := r.db.DB().QueryRowContext(ctx, q, chatID, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return customerrors.NewUserNotInChatError(userID, chatID)
		}
		return err
	}

	return nil
}

// Create inserts a new chat into the database.
func (r *repo) Create(ctx context.Context, users []int64) (int64, error) {
	chatBuilder := sq.Insert(tableChats).
		PlaceholderFormat(sq.Dollar).
		Columns(columnCreatedAt).
		Values("NOW()").
		Suffix("RETURNING id")

	chatQuery, chatArgs, err := chatBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	qChat := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: chatQuery,
	}
	var chatID int64
	err = r.db.DB().ScanOneContext(ctx, &chatID, qChat, chatArgs...)
	if err != nil {
		return 0, err
	}

	chatUsersBuilder := sq.Insert(tableChatUsers).
		PlaceholderFormat(sq.Dollar).
		Columns(columnChatID, columnUserID)

	for _, userID := range users {
		chatUsersBuilder = chatUsersBuilder.Values(chatID, userID)
	}

	chatUsersQuery, chatUsersArgs, err := chatUsersBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	qChatUsers := db.Query{
		Name:     "chat_repository.CreateChatUsers",
		QueryRaw: chatUsersQuery,
	}

	_, err = r.db.DB().ExecContext(ctx, qChatUsers, chatUsersArgs...)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

// Delete removes a chat by ID from the database.
func (r *repo) Delete(ctx context.Context, id int64) error {
	chatUsersDeleteBuilder := sq.Delete(tableChatUsers).
		Where(sq.Eq{columnChatID: id}).
		PlaceholderFormat(sq.Dollar)

	chatUsersDeleteQuery, chatUsersDeleteArgs, err := chatUsersDeleteBuilder.ToSql()
	if err != nil {
		return err
	}

	qChatUsers := db.Query{
		Name:     "chat_repository.DeleteChatUsers",
		QueryRaw: chatUsersDeleteQuery,
	}

	_, err = r.db.DB().ExecContext(ctx, qChatUsers, chatUsersDeleteArgs...)
	if err != nil {
		return err
	}

	chatDeleteBuilder := sq.Delete(tableChats).
		Where(sq.Eq{columnID: id}).
		PlaceholderFormat(sq.Dollar)

	chatDeleteQuery, chatDeleteArgs, err := chatDeleteBuilder.ToSql()
	if err != nil {
		return err
	}

	qChat := db.Query{
		Name:     "chat_repository.Delete",
		QueryRaw: chatDeleteQuery,
	}

	_, err = r.db.DB().ExecContext(ctx, qChat, chatDeleteArgs...)
	if err != nil {
		return err
	}

	return nil
}

// SendMessage handles sending a message from a user to a chat, creates a message in database.
func (r *repo) SendMessage(ctx context.Context, req *pb.SendMessageRequest) error {
	if err := r.chatExists(ctx, req.GetChatId()); err != nil {
		return err
	}

	if err := r.isUserInChat(ctx, req.GetFromUser(), req.GetChatId()); err != nil {
		return err
	}

	builder := sq.Insert(tableMessages).
		PlaceholderFormat(sq.Dollar).
		Columns(columnChatID, columnFromUser, columnText, columnTimestamp).
		Values(req.GetChatId(), req.GetFromUser(), req.GetText(), time.Now().UTC())

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_repository.SendMessage",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
