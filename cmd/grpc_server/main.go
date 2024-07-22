package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mikhailsoldatkin/chat-server/internal/config"
	"github.com/mikhailsoldatkin/chat-server/internal/logger"
	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// tables and fields names
const (
	ID        = "id"
	ChatUsers = "chat_users"
	Chats     = "chats"
	Messages  = "messages"
	ChatID    = "chat_id"
	UserID    = "user_id"
	CreatedAt = "created_at"
	FromUser  = "from_user"
	Timestamp = "timestamp"
	Text      = "text"
)

type server struct {
	pb.UnimplementedChatV1Server
	pool *pgxpool.Pool
}

// recordExists checks if a record with given ID exists in a database table.
func (s *server) recordExists(ctx context.Context, ID int64, table string) error {
	var exists bool
	notFoundMsg := fmt.Sprintf("record with ID %d not found in %s table", ID, table)
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=%d)", table, ID)
	err := s.pool.QueryRow(ctx, query).Scan(&exists)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to check existence: %v", err)
	}
	if !exists {
		return status.Errorf(codes.NotFound, notFoundMsg)
	}
	return nil
}

// isUserInChat checks if a user is a member of a chat.
func (s *server) isUserInChat(ctx context.Context, userID int64, chatID int64) error {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE chat_id=$1 AND user_id=$2)", ChatUsers)
	err := s.pool.QueryRow(ctx, query, chatID, userID).Scan(&exists)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to check if user is in chat: %v", err)
	}
	if !exists {
		return status.Errorf(codes.InvalidArgument, "user %d is not a participant of chat %d", userID, chatID)
	}
	return nil
}

// Create handles the creation of a new chat with provided users.
func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	if len(req.Users) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no users provided for the chat")
	}

	tr, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tr, ctx)

	chatBuilder := sq.Insert(Chats).
		PlaceholderFormat(sq.Dollar).
		Columns(CreatedAt).
		Values("NOW()").
		Suffix("RETURNING id")

	chatQuery, chatArgs, err := chatBuilder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build chat insert query: %v", err)
	}

	var chatID int
	err = tr.QueryRow(ctx, chatQuery, chatArgs...).Scan(&chatID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert chat: %v", err)
	}

	chatUserBuilder := sq.Insert(ChatUsers).
		PlaceholderFormat(sq.Dollar).
		Columns(ChatID, UserID)

	for _, userID := range req.Users {
		chatUserBuilder = chatUserBuilder.Values(chatID, userID)
	}

	chatUserQuery, userArgs, err := chatUserBuilder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build chat_users insert query: %v", err)
	}

	_, err = tr.Exec(ctx, chatUserQuery, userArgs...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert users into chat_users: %v", err)
	}

	err = tr.Commit(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	logger.Info("chat %d created with users: %v", chatID, req.Users)

	return &pb.CreateResponse{Id: int64(chatID)}, nil
}

// Delete removes a chat by ID.
func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	tr, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer func(tr pgx.Tx, ctx context.Context) {
		_ = tr.Rollback(ctx)
	}(tr, ctx)

	chatID := req.GetId()
	chatUserDeleteBuilder := sq.Delete(ChatUsers).
		Where(sq.Eq{ChatID: chatID}).
		PlaceholderFormat(sq.Dollar)

	chatUserDeleteQuery, chatUserDeleteArgs, err := chatUserDeleteBuilder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build chat_users delete query: %v", err)
	}

	_, err = tr.Exec(ctx, chatUserDeleteQuery, chatUserDeleteArgs...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete users from chat_users: %v", err)
	}

	logger.Info("users deleted from chat %d", chatID)

	chatDeleteBuilder := sq.Delete(Chats).
		Where(sq.Eq{ID: chatID}).
		PlaceholderFormat(sq.Dollar)

	chatDeleteQuery, chatDeleteArgs, err := chatDeleteBuilder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build chats delete query: %v", err)
	}

	_, err = tr.Exec(ctx, chatDeleteQuery, chatDeleteArgs...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete chat: %v", err)
	}

	err = tr.Commit(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	logger.Info("chat %d deleted", chatID)

	return &emptypb.Empty{}, nil
}

// SendMessage handles sending a message to a chat from a user.
func (s *server) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*emptypb.Empty, error) {
	chatID := req.GetChatId()
	userID := req.GetFromUser()
	messageText := req.GetText()

	if err := s.recordExists(ctx, chatID, Chats); err != nil {
		return nil, err
	}

	// check user existence

	if err := s.isUserInChat(ctx, userID, chatID); err != nil {
		return nil, err
	}

	messageBuilder := sq.Insert(Messages).
		PlaceholderFormat(sq.Dollar).
		Columns(ChatID, FromUser, Text, Timestamp).
		Values(chatID, userID, messageText, time.Now().UTC()).
		Suffix("RETURNING id")

	query, args, err := messageBuilder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build insert query: %v", err)
	}

	var messageID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&messageID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send message: %v", err)
	}

	logger.Info("user %d sent message %d to chat %d", userID, messageID, chatID)

	return &emptypb.Empty{}, nil
}

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, cfg.Database.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.GRPCPort))
	if err != nil {
		logger.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterChatV1Server(s, &server{pool: pool})

	logger.Info("%v server listening at %v", cfg.AppName, lis.Addr())

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", err)
	}
}
