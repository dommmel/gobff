package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/dommmel/gobff/user"
)

type server struct {
	db *sql.DB
	pb.UnimplementedUserServiceServer
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (s *server) GetUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user := &User{}
	err := s.db.QueryRow("SELECT id, name FROM users WHERE id = ?", in.Id).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with ID %d", in.Id)
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		log.Printf("Error fetching user with ID %d: %v", in.Id, err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}
	return &pb.UserResponse{Id: int32(user.ID), Name: user.Name}, nil
}

func main() {
	db, err := sql.Open("sqlite3", "./user.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{db: db})

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down GRPC server...")
		s.GracefulStop()
	}()

	log.Println("Starting GRPC server on :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
