package realtime

import (
	"log"
	"realTimeChat/servergrpc"
	"strconv"
)

func GetFriends(userID string) <-chan []string {
	ch := make(chan []string)
	go func() {
		defer close(ch)
		// friends := []string{"1", "2", "3"}
		friends, err := servergrpc.GetUserFollowingFollowers(userID)
		if err != nil {
			// Handle error appropriately
			log.Printf("Error getting friends for user %s: %v", userID, err)
			return
		}
		var friendIDs []string
		for _, friend := range friends {
			for _, id := range friend.UserIds {
				friendIDs = append(friendIDs, strconv.FormatUint(uint64(id), 10))
			}
		}
		ch <- friendIDs
	}()
	return ch
}
