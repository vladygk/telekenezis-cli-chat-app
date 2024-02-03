package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/vladygk/chat-app/model"
)

type Connection struct {
	ID     string
	Socket *websocket.Conn
}

type Service struct {
	connections map[string]*Connection
}

func Initialize() *Service {
	service := &Service{

		connections: make(map[string]*Connection),
	}
	return service
}

func (s *Service) Run() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.HandleRequest(w, r)
	})
}

func (s *Service) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var doneCh = make(chan struct{})
	defer close(doneCh)
	go s.OpenSocketConnection(w, r, doneCh)

	<-doneCh
}

func (s *Service) OpenSocketConnection(w http.ResponseWriter, r *http.Request, doneCh chan struct{}) {
	respHeader := make(http.Header)
	c, err := upgrader.Upgrade(w, r, respHeader)
	if err != nil {
		return
	}
	defer c.Close()

	s.registerClient(c)
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}
		var message model.Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			_ = fmt.Errorf("Error: %s", err)
		}

		s.handleMessage(message)
	}
	doneCh <- struct{}{}
}

func (s *Service) registerClient(connection *websocket.Conn) {

	connectionID := uuid.New().String()
	fmt.Printf("Your id is: %s\n", connectionID)
	s.connections[connectionID] = &Connection{
		ID:     connectionID,
		Socket: connection,
	}
	connection.WriteJSON(connectionID)
	connection.SetCloseHandler(func(code int, text string) error {
		s.unregisterClient(connectionID, connectionID)
		return nil
	})
}

func (s *Service) unregisterClient(userID, connectionID string) {
	delete(s.connections, userID)
}

func (s *Service) handleMessage(message model.Message) {
	connection, ok := s.connections[message.To]
	if !ok {
		fmt.Println("Invalid destination id.")
		return
	}
	err := connection.Socket.WriteJSON(message)
	if err != nil {
		fmt.Println("Error sending msg to client.")
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Service) StartListening(port int) {
	fmt.Printf("Server started at port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
