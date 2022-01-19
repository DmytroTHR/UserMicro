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
	case "/user.UserService/SetUsersRole":
		err = checkIfAdmin(ctx)
	case "/user.UserService/GetUserByID":
		err = checkIfYourself(ctx, req)
		if err != nil {
			err = checkIfAdmin(ctx)
		}
	}
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func checkIfYourself(ctx context.Context, request interface{}) error {
	userRequest, ok := request.(*proto.User)
	if !ok {
		return status.Errorf(codes.PermissionDenied, "Wrong request parameter")
	}

	validationResult, err := getTokenValidationResult(ctx)
	if err != nil {
		return err
	}

	if validationResult.User == nil || validationResult.User.Id != userRequest.Id {
		return status.Errorf(codes.PermissionDenied, "No permission for this operation")
	}

	return nil
}

func checkIfAdmin(ctx context.Context) error {

	validationResult, err := getTokenValidationResult(ctx)
	if err != nil {
		return err
	}
	role := validationResult.Role
	if role == nil || !role.IsAdmin {
		return status.Errorf(codes.PermissionDenied, "No permission for this operation")
	}

	return nil
}

func getTokenValidationResult(ctx context.Context) (*proto.Response, error)  {
	result := &proto.Response{}
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return result, status.Errorf(codes.InvalidArgument, "No metadata retreived")
	}

	auth, ok := meta["authorization"]
	if !ok {
		return result, status.Errorf(codes.Unauthenticated, "No authorization token found")
	}

	token := &proto.Token{Token: auth[0]}
	userService := service.NewUserService(nil, &service.TokenService{})
	result, err := userService.ValidateToken(ctx, token)
	if err != nil {
		return result, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return result, nil
}