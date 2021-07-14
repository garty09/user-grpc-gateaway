package internal

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	test "user_test/gen/go/proto"
	"user_test/internal/model"
	"user_test/internal/pagination"
)

var _ test.UserServiceServer = (*RPCServer)(nil)

type RPCServer struct {
	test.UserServiceServer
	Repo *model.Repository
}

func (r *RPCServer) DeleteUser(ctx context.Context, input *test.DeleteUserRequest) (*test.DeleteUserResponse, error) {
	err := r.Repo.DeleteUser(ctx, input.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "internal server error")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &test.DeleteUserResponse{
		Deleted: true,
	}, nil
}

func (r *RPCServer) InsertUser(ctx context.Context, input *test.InsertUserRequest) (*test.InsertUserResponse, error) {
	err := r.Repo.InsertUser(ctx, input)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "internal server error")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &test.InsertUserResponse{
		Added: true,
	}, nil
}

func (r *RPCServer) ListUsers(ctx context.Context, input *test.ListUserRequest) (*test.ListUserResponse, error) {
	p := pagination.New(int(input.Page))
	users, err := r.Repo.ListUsers(ctx, p.Offset(), p.Limit())
	if err != nil {
		return nil, err
	}

	resp := &test.ListUserResponse{
		Page:  input.Page,
		Limit: int64(p.Limit()),
	}

	for _, user := range users {
		resp.Users = append(resp.Users, &test.User{
			Id:    user.ID,
			Fio:   user.FIO,
			Email: user.Email,
			Phone: user.PhoneNumber,
		})
	}

	return resp, nil
}
