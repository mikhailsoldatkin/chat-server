package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/mikhailsoldatkin/chat-server/internal/repository"
	repoMocks "github.com/mikhailsoldatkin/chat-server/internal/repository/mocks"
	"github.com/mikhailsoldatkin/chat-server/internal/service/chat"
	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()
	type chatRepoMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx context.Context
		req *pb.SendMessageRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		chatID = gofakeit.Int64()
		userID = gofakeit.Int64()
		msg    = gofakeit.BeerName()

		req = &pb.SendMessageRequest{
			ChatId:   chatID,
			FromUser: userID,
			Text:     msg,
		}
		wantErr = fmt.Errorf("repository error")
	)

	tests := []struct {
		name         string
		args         args
		err          error
		chatRepoMock chatRepoMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			chatRepoMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, req).Return(nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: wantErr,
			chatRepoMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, req).Return(wantErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatRepoMock := tt.chatRepoMock(mc)
			service := chat.NewMockService(chatRepoMock)

			repoErr := service.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, repoErr)
		})
	}
}
