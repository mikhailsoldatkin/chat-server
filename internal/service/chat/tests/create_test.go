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
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type chatRepoMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx context.Context
		req []int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id     = gofakeit.Int64()
		userID = gofakeit.Int64()
		req    = []int64{userID}

		wantErr = fmt.Errorf("repository error")
	)

	tests := []struct {
		name         string
		args         args
		want         int64
		err          error
		chatRepoMock chatRepoMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			chatRepoMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  wantErr,
			chatRepoMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, wantErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatRepoMock := tt.chatRepoMock(mc)
			service := chat.NewMockService(chatRepoMock)

			resp, repoErr := service.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, repoErr)
			require.Equal(t, tt.want, resp)
		})
	}
}
