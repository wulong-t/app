package realtime

func GetFriends(userID string) <-chan []string { 
	ch := make(chan []string)
	go func() { 
		defer close(ch)
		friends := []string{"1", "2", "3"}
		ch <- friends
	}()
	return ch
}