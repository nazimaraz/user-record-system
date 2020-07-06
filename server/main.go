package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	pb "UserRecordSystem/server/proto"

	"github.com/fatih/structs"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
)

type server struct{}

type user struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	IPAddress string `json:"ip_address"`
	UserName  string `json:"user_name"`
	Agent     string `json:"agent"`
	Country   string `json:"country"`
}

var (
	ctx    = context.Background()
	client = newClient()
	port   = flag.Int("port", 9000, "The server port")
)

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}

func getLastKey() (string, error) {
	lastKey, err := client.Get("last_key").Result()
	if err != nil {
		err = client.Set("last_key", 1, 0).Err()
		lastKey = "1"
		if err != nil {
			return "", nil
		}
	} else {
		var lastKeyInt, err = strconv.Atoi(lastKey)
		lastKeyInt++
		if err != nil {
			return "", err
		}
		lastKey = strconv.Itoa(lastKeyInt)
		err = client.Set("last_key", lastKey, 0).Err()
		if err != nil {
			return "", err
		}
	}
	return lastKey, nil
}

func saveUser(client *redis.Client, r *pb.User) error {
	usr := user{
		ID:        r.Id,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Gender:    r.Gender,
		IPAddress: r.IpAddress,
		UserName:  r.UserName,
		Agent:     r.Agent,
		Country:   r.Country,
	}

	usrM := structs.Map(usr)
	lastKey, err := getLastKey()
	if err != nil {
		return err
	}

	err = client.HMSet("user:"+lastKey, usrM).Err()
	if err != nil {
		return err
	}
	log.Printf("user:%s saved.", lastKey)
	return nil
}

func (s *server) Save(ctx context.Context, r *pb.User) (*pb.SaveResponse, error) {
	err := saveUser(client, r)
	if err != nil {
		fmt.Println(err)
	}
	return &pb.SaveResponse{Body: "User saved!"}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
