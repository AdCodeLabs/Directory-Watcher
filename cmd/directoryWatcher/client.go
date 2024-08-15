package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	ipAddr := flag.String("server", "127.0.0.1:8089", "enter server ip address")
	flag.Parse()

	serverURL := fmt.Sprintf("ws://%s/ws", *ipAddr)
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatalf("Error connecting to WebSocket server: %v", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("connection closed")
		}
	}(conn)

	log.Println("Connected to WebSocket server:", serverURL)

	go listenForMessages(conn, &wg)
	wg.Wait()
	fmt.Println("Exiting...")
}

func listenForMessages(conn *websocket.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}
		fmt.Println("Received:", string(msg))
	}
}
