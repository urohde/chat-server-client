package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/url"
	"os"
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
	if len(os.Args) < 3 {
		fmt.Println("Please provide a URL and a username")
		os.Exit(1)
	}

	input := make(chan string, 1)
	go getInput(input)

	query := url.Values{}
	query.Add("username", os.Args[2])

	URL := url.URL{Scheme: "ws", Host: os.Args[1], Path: "/", RawQuery: query.Encode()}
	conn, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	if err != nil {
		fmt.Println("Error connecting to server: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if err == io.EOF || websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					fmt.Println("Connection closed by server")
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
		}
	}
}
