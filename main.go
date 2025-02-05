package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var broadcat = make(chan string)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg string
		err := ws.ReadJSON(&msg)

		if err != nil {
			delete(clients, ws)
			break
		}
		broadcat <- msg
	}
}

func handleMessage() {
	for {
		msg := <-broadcat
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	go handleMessage()

	fmt.Println("Server run on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
