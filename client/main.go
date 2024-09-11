package main

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var in = bufio.NewReader(os.Stdin)

func getInput(input chan string) {
	result, err := in.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input: ", err)
		os.Exit(1)
	}

	result = result[:len(result)-1]

	// Clear the input line
	fmt.Print("\033[A\033[K")

	input <- result
}

func main() {
	username := os.Getenv("USERNAME")
	if username == "" {
		username = "anonymous"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost:8080"
	}

	input := make(chan string, 1)
	go getInput(input)

	query := url.Values{}
	query.Add("username", username)

	URL := url.URL{Scheme: "ws", Host: host, Path: "/", RawQuery: query.Encode()}
	conn, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	if err != nil {
		fmt.Println("Error connecting to server: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	done := make(chan struct{})

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if err == io.EOF || websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					fmt.Println("Connection closed")
				} else {
					fmt.Println("Error reading message: ", err)
				}
				return
			}

			fmt.Println(string(message))
		}
	}()

	for {
		select {
		case <-done:
			return
		case msg := <-input:
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("Error sending message: ", err)
				return
			}
			go getInput(input)
		case <-interrupt:
			fmt.Println("Interrupt signal received")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("Error sending close message: ", err)
			}
			select {
			case <-done:
			case <-time.After(2 * time.Second):
			}
			return
		}
	}
}
