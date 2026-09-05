package realtime

import (
	"errors"
	"log"
	"realTimeChat/servergrpc"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
}

type ConnectionManager struct {
	connections    map[string]*websocket.Conn
	onlineFriends  map[string][]string
	getUserFriends func(string) <-chan []string
	lock           sync.Mutex
}

func NewConnectionManager(getUserFriends func(string) <-chan []string) (*ConnectionManager, error) {
	if getUserFriends == nil {
		return nil, errors.New("getUserFriends is required")
	}
	return &ConnectionManager{
		connections:    make(map[string]*websocket.Conn),
		onlineFriends:  make(map[string][]string),
		getUserFriends: getUserFriends,
	}, nil
}

func (cm *ConnectionManager) AddConnection(userID string, conn *websocket.Conn) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.connections[userID] = conn
	cm.onlineFriends[userID] = []string{}

	for onlineID := range cm.onlineFriends {
		yes, err := cm.isFriend(onlineID, userID)
		if err != nil {
			continue
		}
		if onlineID != userID && yes {
			cm.onlineFriends[onlineID] = append(cm.onlineFriends[onlineID], userID)
			err := cm.connections[onlineID].WriteJSON(gin.H{
				"type":   "friend_online",
				"userID": userID,
			})
			if err != nil {
				log.Printf("Error notifying %s about %s", onlineID, userID)
				return
			}
		}
	}

	go func() {
		for friends := range cm.getUserFriends(userID) {
			if friends == nil {
				log.Printf("Failed to get friends for user %s", userID)
				return
			}
			cm.lock.Lock()
			defer cm.lock.Unlock()
			for _, friendID := range friends {
				if cm.connections[friendID] == nil {
					continue
				}
				cm.onlineFriends[userID] = append(cm.onlineFriends[userID], friendID)
				err := cm.connections[userID].WriteJSON(gin.H{
					"type":   "friend_online",
					"userID": friendID,
				})
				if err != nil {
					log.Printf("Error notifying %s about %s", userID, friendID)
					return
				}
			}
		}
	}()
}
func (cm *ConnectionManager) RemoveConnection(userID string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	delete(cm.connections, userID)
	delete(cm.onlineFriends, userID)
	for onlineID := range cm.onlineFriends {
		yes, err := cm.isFriend(onlineID, userID)
		if err != nil {
			log.Printf("Error checking friendship between %s and %s: %v", onlineID, userID, err)
			return
		}
		if onlineID != userID && yes {
			cm.onlineFriends[onlineID] = removeFriend(cm.onlineFriends[onlineID], userID)
			err := cm.connections[onlineID].WriteJSON(gin.H{
				"type":   "friend_offline",
				"userID": userID,
			})
			if err != nil {
				log.Printf("Error notifying %s about %s", onlineID, userID)
				return
			}
		}
	}
}

func (cm *ConnectionManager) SendMessage(message *Message) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	conn, exists := cm.connections[message.Receiver]
	if !exists {
		log.Printf("Error: receiver %s is not connected", message.Receiver)
		return
	}
	err := conn.WriteJSON(gin.H{
		"type":    "message",
		"sender":  message.Sender,
		"content": message.Content,
	})
	if err != nil {
		log.Printf("Error sending message to %s: %v", message.Receiver, err)
		return
	}
	err = servergrpc.SendMessage(message.Sender, message.Receiver, message.Content)
	if err != nil {
		log.Fatalf("Error sending message via gRPC: %v", err)
	}
}

func (cm *ConnectionManager) isFriend(userID, friendID string) (bool, error) {
	friendsCh := cm.getUserFriends(userID)
	friends, ok := <-friendsCh
	if !ok {
		return false, errors.New("failed to receive friends list")
	}
	for _, friend := range friends {
		if friend == friendID {
			return true, nil
		}
	}
	return false, nil
}

func removeFriend(friends []string, friendID string) []string {
	for i, friend := range friends {
		if friend == friendID {
			return append(friends[:i], friends[i+1:]...)
		}
	}
	return friends
}
