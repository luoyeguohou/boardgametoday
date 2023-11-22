package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// 初始化各个service
	ServiceMGR.InitAllService()
	// 创建HTTP服务器
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级HTTP连接为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	ply := Player{
		NickName: "落叶过后",
		HeadPic:  "111",
		Conn:     conn,
	}

	defer func() {
		conn.Close()
		ply.OnExit()
		ServiceMGR.OnPlyExit(&ply)
	}()
	ply.OnInit()

	// 处理WebSocket连接
	for {
		// 读取消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		ply.ProcessMsg(p)
		log.Println(messageType)
		log.Println("Received message:", string(p))

		// // 发送消息
		// err = conn.WriteMessage(messageType, []byte("Hello, world!"))
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
	}
}

// 明日偷师任务  SIMap Modules
