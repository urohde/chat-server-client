package main

import (
	"fmt"
	"net/http"
	"os"
	"server/client"
	"server/logger"
	"server/server"

	"github.com/gorilla/websocket"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing log file path")
		os.Exit(1)
	}

	// logger, err := logger.NewNoopLogger()
	logger, err := logger.NewFileLogger(os.Args[1])
	if err != nil {
		fmt.Println("Error creating logger: ", err)
		return
	}
	defer logger.Close()

	chatServer := server.NewServer(logger)
	go chatServer.Start()

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		client := client.NewClient(ws, username)
		chatServer.Register(client)
	})

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
