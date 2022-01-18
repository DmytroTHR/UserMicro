//go:generate protoc -I=./proto --go_out=./proto ./proto/user.proto --go-grpc_out=./proto ./proto/user.proto

package main

import (
	"UserMicro/configs"
	"UserMicro/proto"
	"UserMicro/service"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const MicroName = "user_service"

func main() {
	connectionDB := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		configs.PG_HOST,
		configs.PG_PORT,
		configs.POSTGRES_USER,
		configs.POSTGRES_PASSWORD,
		configs.POSTGRES_DB)

	db, err := sql.Open("postgres", connectionDB)
	if err != nil {
		log.Panicf("%s: failed to open db connection - %v", MicroName, err)
	}
	defer db.Close()

	service := service.NewUserService(db, &service.TokenService{})

	listener, err := net.Listen("tcp", net.JoinHostPort("", configs.GRPC_PORT))
	if err != nil {
		log.Panicf("%s: failed to listen on port - %v", MicroName, err)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(serverInterceptor))
	defer server.GracefulStop()
	proto.RegisterUserServiceServer(server, service)
	reflection.Register(server)

	if err := server.Serve(listener); err != nil {
		log.Panicf("%s: failed to start grpc - %v", MicroName, err)
	}
}
