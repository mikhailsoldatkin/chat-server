package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	chatAPI "github.com/mikhailsoldatkin/chat-server/internal/api/chat"
	"github.com/mikhailsoldatkin/chat-server/internal/customerrors"
	"github.com/mikhailsoldatkin/chat-server/internal/service"
	"github.com/mikhailsoldatkin/chat-server/internal/service/chat/model"
	serviceMocks "github.com/mikhailsoldatkin/chat-server/internal/service/mocks"

	pb "github.com/mikhailsoldatkin/chat-server/pkg/chat_v1"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *pb.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id     = gofakeit.Int64()
		userID = gofakeit.Int64()
		users  = []int64{userID}

		req = &pb.CreateRequest{Users: users}

		wantResp     = &pb.CreateResponse{Id: id}
		wantChat     = &model.Chat{ID: id}
		wantChatUser = &model.ChatUser{
			ChatID: id,
			UserID: userID,
		}
		wantErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *pb.CreateResponse
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: wantResp,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateMock.Set(func(_ context.Context, _ []int64) (int64, error) {
					require.Equal(t, wantChat.ID, id)
					require.Equal(t, wantChatUser.ChatID, id)
					require.Equal(t, wantChatUser.UserID, userID)
					return id, nil
				})
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  customerrors.ConvertError(wantErr),
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.CreateMock.Set(func(_ context.Context, _ []int64) (int64, error) {
					return 0, wantErr
				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatServiceMock := tt.chatServiceMock(mc)
			api := chatAPI.NewMockImplementation(chatServiceMock)

			resp, grpcErr := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, grpcErr)
			require.Equal(t, tt.want, resp)
		})
	}
}
