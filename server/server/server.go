package server

import (
	"fmt"
	"strings"
	"time"
)

type client interface {
	GetName() string
	SendMessage([]byte) error
	ReadMessage() (int, []byte, error)
	Close() error
}

type logger interface {
	Write([]byte) error
}

type Server struct {
	clients    map[client]bool
	broadcast  chan []byte
	register   chan client
	unregister chan client
	logger     logger
}

func NewServer(logger logger) *Server {
	return &Server{
		clients:    make(map[client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan client),
		unregister: make(chan client),
		logger:     logger,
	}
}
func (s *Server) Broadcast(msg []byte) {
	messageWithTimestamp := append([]byte(fmt.Sprintf("[%s] ", time.Now().Format(time.RFC3339))), msg...)
	err := s.logger.Write(messageWithTimestamp)
	if err != nil {
		fmt.Println(err)
	}
	s.broadcast <- msg
}

func (s *Server) Register(client client) {
	s.register <- client
}

func (s *Server) Unregister(client client) {
	s.unregister <- client
}

func (s *Server) Start() {
	for {
		select {
		case c := <-s.register:
			s.clients[c] = true
			fmt.Printf("Client [%s] connected\n", c.GetName())
			connectedUsers := strings.Builder{}
			for client := range s.clients {
				if connectedUsers.Len() > 0 {
					connectedUsers.WriteString(", ")
				}
				connectedUsers.WriteString(client.GetName())
			}
			c.SendMessage([]byte(fmt.Sprintf("Welcome %s! Users online [%s]", c.GetName(), connectedUsers.String())))
			go s.listen(c)
		case c := <-s.unregister:
			err := c.Close()
			if err != nil {
				fmt.Println(err)
				break
			}

			delete(s.clients, c)
			fmt.Printf("Client [%s] disconnected\n", c.GetName())
		case message := <-s.broadcast:
			for client := range s.clients {
				err := client.SendMessage(message)
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		}
	}
}

func (s *Server) listen(c client) {
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		if string(msg) == "bye" {
			fmt.Printf("Client [%s] requested to disconnect\n", c.GetName())
			s.Unregister(c)
			return
		}

		msg = append([]byte(c.GetName()+": "), msg...)
		s.Broadcast(msg)
	}

}
