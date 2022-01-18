package main

import (
	"UserMicro/proto"
	"UserMicro/service"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func serverInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	var err error
	switch info.FullMethod {
	case "/proto.UserService/SetUsersRole", "/proto.UserService/GetUserByID":
		err = checkIfAdmin(ctx)
	}
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func checkIfAdmin(ctx context.Context) error {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "No metadata retreived")
	}

	auth, ok := meta["authorization"]
	if !ok {
		return status.Errorf(codes.Unauthenticated, "No authorization token found")
	}

	token := &proto.Token{Token: auth[0]}
	userService := service.NewUserService(nil, &service.TokenService{})
	_, err := userService.ValidateToken(ctx, token)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, err.Error())
	}

	return nil
}
