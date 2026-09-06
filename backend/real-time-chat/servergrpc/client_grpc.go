package servergrpc

import (
	"context"
	"log"
	"realTimeChat/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendMessage(senderID string, receiverID string, message string) error {
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := protos.NewRealtimeChatServiceClient(conn)
	_, err = client.SendMessage(context.Background(), &protos.MessageRequest{
		Sender:  senderID,
		Receiver: receiverID,
		Content: message,
	})
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}
	return nil
}

func GetUserFollowingFollowers(userID string) ([]*protos.UserIDs, error) {
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := protos.NewRealtimeChatServiceClient(conn)
	resp, err := client.GetUserFollowingFollowers(context.Background(), &protos.GetUserFollowingFollowersRequest{
		UserId: userID,
	})
	if err != nil {
		log.Fatalf("could not get user following/followers: %v", err)
	}
	return resp.Followers, nil
}