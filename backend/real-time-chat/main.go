package main

import (
	"log"
	"net/http"
	"realTimeChat/realtime"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// 浏览器 WS 握手会带 Origin，需要校验；开发阶段直接放行
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	app := gin.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
	}))
	manager, err := realtime.NewConnectionManager(realtime.GetFriends)
	if err != nil {
		log.Fatal("Failed to create connection manager:", err)
	}
	app.GET("/ws/:id", func(c *gin.Context) {
		id := c.Param("id")
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer conn.Close()

		log.Println("user connected:", id)
		conn.WriteMessage(websocket.TextMessage, []byte("Welcome user "+id+"!"))

		manager.AddConnection(id, conn)
		defer conn.Close()
		defer manager.RemoveConnection(id)
		// 读循环保持连接，handler 返回连接就会被关掉
		var message realtime.Message
		for {
			if err := conn.ReadJSON(&message); err != nil {
				log.Printf("user disconnected: %s because: %v", id, err)
				break
			}
			manager.SendMessage(&message)
		}
	})
	log.Fatal(app.Run(":8083"))
}
