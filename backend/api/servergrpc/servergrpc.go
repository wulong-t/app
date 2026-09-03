package servergrpc

import (
	"Server/database"
	"Server/models"
	pb "Server/protos"
	"context"
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
