package service

import (
	"UserMicro/proto"
	"UserMicro/repository"
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo         *repository.UserRepo
	tokenService *TokenService
	*proto.UnimplementedUserServiceServer
}

func NewUserService(repo *sql.DB, tokenServ *TokenService) *UserService {
	return &UserService{
		Repo:         repository.NewUserRepo(repo),
		tokenService: tokenServ,
	}
}

func (serv *UserService) CreateUser(ctx context.Context, user *proto.User) (*proto.Response, error) {
	userCreated, err := serv.Repo.Create(ctx, user)
	return &proto.Response{
		Success: err == nil,
		User:    userCreated,
	}, err
}

func (serv *UserService) GetUserByID(ctx context.Context, user *proto.User) (*proto.Response, error) {
	userFound, err := serv.Repo.ReadByID(ctx, user.Id)
	return &proto.Response{
		Success: err == nil,
		User:    userFound,
	}, err
}

func (serv *UserService) GetUsersRole(ctx context.Context, user *proto.User) (*proto.Response, error) {
	roleFound, err := serv.Repo.GetRoleByUser(ctx, user)
	return &proto.Response{
		Success: err == nil,
		Role:    roleFound,
	}, err
}

func (serv *UserService) SetUsersRole(ctx context.Context, user *proto.User) (*proto.Response, error) {
	err := serv.Repo.UpdateRole(ctx, user, user.Role)
	return &proto.Response{
		Success: err == nil,
	}, err
}

func (serv *UserService) AuthUser(ctx context.Context, user *proto.User) (*proto.Response, error) {
	userDB, err := serv.Repo.ReadByEmail(ctx, user.Email)
	if err != nil {
		return &proto.Response{
			Success: err == nil,
		}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password)); err != nil {
		return &proto.Response{Success: err == nil}, err
	}

	token, err := serv.tokenService.Encode(userDB)
	return &proto.Response{
		Success: err == nil,
		Token:   &proto.Token{Token: token},
	}, err
}

func (serv *UserService) ValidateToken(ctx context.Context, token *proto.Token) (*proto.Response, error) {
	claims, err := serv.tokenService.Decode(token.Token)

	if err != nil {
		return &proto.Response{Success: err == nil}, err
	}

	if claims.User.Id == 0 {
		err = errors.New("invalid user")
		return &proto.Response{Success: err == nil}, err
	}

	token.Valid = true

	return &proto.Response{
		Success: true,
		User:    claims.User,
		Role:    claims.User.Role,
		Token:   token}, nil

}
