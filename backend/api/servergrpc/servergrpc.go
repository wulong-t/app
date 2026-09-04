package servergrpc

import (
	"Server/database"
	"Server/models"
	pb "Server/protos"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Server struct {
	pb.UnimplementedRealtimeChatServiceServer
}

func (s *Server) GetUserFollowingFollowers(ctx context.Context, req *pb.GetUserFollowingFollowersRequest) (*pb.UserFollowingFollowersResponse, error) {
	userId := req.GetUserId()
	if userId == "" {
		return nil, fmt.Errorf("userId is required")
	}
	userIdUint, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid userId format")
	}
	// Implementation for fetching following and followers
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var followers []models.FollowModel
	followers, err = gorm.G[models.FollowModel](database.DB).Where("followed_id = ?", userIdUint).Find(c)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch followers: %w", err)
	}
	var friends []uint
	for _, follower := range followers {
		friends = append(friends, follower.FollowerID)
	}
	// Process the fetched followers and return the response
	ids := make([]uint32, 0, len(friends))
	for _, id := range friends {
		ids = append(ids, uint32(id))
	}
	return &pb.UserFollowingFollowersResponse{
		Followers: []*pb.UserIDs{
			{UserIds: ids},
		},
	}, nil
}

func (s *Server) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	if req.GetSender() == "" || req.GetReceiver() == "" {
		return nil, fmt.Errorf("senderId or receiverId is required")
	}
	senderID, err := strconv.ParseUint(req.GetSender(), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid senderId format")
	}
	receiverID, err := strconv.ParseUint(req.GetReceiver(), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid receiverId format")
	}

	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// Check if sender and receiver exist
	if _, err = gorm.G[models.UserModel](database.DB).Where("id = ?", senderID).First(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("sender not found")
		}
		return nil, fmt.Errorf("failed to fetch sender: %w", err)
	}
	if _, err = gorm.G[models.UserModel](database.DB).Where("id = ?", receiverID).First(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("receiver not found")
		}
		return nil, fmt.Errorf("failed to fetch receiver: %w", err)
	}

	msg := models.MessageModel{
		SenderID:   uint(senderID),
		ReceiverID: uint(receiverID),
		Content:    req.GetContent(),
	}
	if err := gorm.G[models.MessageModel](database.DB).Create(c, &msg); err != nil {
		return nil, fmt.Errorf("failed to store message: %w", err)
	}
	if _, err := gorm.G[models.UnReadNumModel](database.DB).Where("sender_id = ? AND receiver_id = ?", senderID, receiverID).First(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			unread := models.UnReadNumModel{
				SenderID:   uint(senderID),
				ReceiverID: uint(receiverID),
				Num:        1,
			}
			if err := gorm.G[models.UnReadNumModel](database.DB).Create(c, &unread); err != nil {
				return nil, fmt.Errorf("failed to create unread count: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to fetch unread count: %w", err)
		}
	}
	return &pb.MessageResponse{
		MessageId: uint64(msg.ID),
	}, nil
}
