package app

import (
	"context"
	"log"

	accessProto "github.com/mikhailsoldatkin/auth/pkg/access_v1"
	"github.com/mikhailsoldatkin/chat-server/internal/api/chat"
	"github.com/mikhailsoldatkin/chat-server/internal/config"
	"github.com/mikhailsoldatkin/chat-server/internal/repository"
	chatRepository "github.com/mikhailsoldatkin/chat-server/internal/repository/chat"
	"github.com/mikhailsoldatkin/chat-server/internal/service"
	chatService "github.com/mikhailsoldatkin/chat-server/internal/service/chat"
	"github.com/mikhailsoldatkin/platform_common/pkg/closer"
	"github.com/mikhailsoldatkin/platform_common/pkg/db"
	"github.com/mikhailsoldatkin/platform_common/pkg/db/pg"
	"github.com/mikhailsoldatkin/platform_common/pkg/db/transaction"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type serviceProvider struct {
	config         *config.Config
	dbClient       db.Client
	txManager      db.TxManager
	chatRepository repository.ChatRepository
	chatService    service.ChatService
	//accessClient   accessProto.AccessV1Client

	chatImplementation *chat.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) Config() *config.Config {
	if s.config == nil {
		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}
		s.config = cfg
	}

	return s.config
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.Config().Database.PostgresDSN)
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("db ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.ChatRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.chatService
}

func (s *serviceProvider) AccessClient() accessProto.AccessV1Client {
	creds, err := credentials.NewClientTLSFromFile("cert/ca.cert", "")
	if err != nil {
		log.Fatalf("failed to process credentials: %v", err)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to create connection: %v", err)
	}
	closer.Add(conn.Close)

	client := accessProto.NewAccessV1Client(conn)

	return client
}

func (s *serviceProvider) ChatImplementation(ctx context.Context) *chat.Implementation {
	if s.chatImplementation == nil {
		s.chatImplementation = chat.NewImplementation(s.ChatService(ctx))
	}

	return s.chatImplementation
}
