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

type Logger interface {
	Write([]byte) (error)
	Close() error
}

func main() {
	var l Logger
	var err error
	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		l = logger.NewNoopLogger()
	} else {
		l, err = logger.NewFileLogger(logFile)
		if err != nil {
			fmt.Println("Error creating logger: ", err)
			return
		}
		defer l.Close()
	}

	chatServer := server.NewServer(l)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
