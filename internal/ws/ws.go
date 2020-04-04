package ws

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocket struct {
	conn   *websocket.Conn
	events map[string]EventHandler

	out chan []byte
	in  chan []byte
}

func NewWebSocket(w http.ResponseWriter, r *http.Request) (*WebSocket, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR | SOCKET CONNECT] %v", err)
		return nil, err
	}
	ws := &WebSocket{
		conn:   conn,
		out:    make(chan []byte),
		in:     make(chan []byte),
		events: make(map[string]EventHandler),
	}
	go ws.reader()
	go ws.writer()
	return ws, nil
}

func (ws *WebSocket) SetHandler(event string, action EventHandler) *WebSocket {
	ws.events[strings.ToLower(event)] = action
	return ws
}

func (ws *WebSocket) Send(event string, cid int, data interface{}) {
	go func() {
		ws.out <- (&Event{
			Name: strings.ToLower(event),
			CID:  cid,
			Data: data,
		}).Raw()
	}()
}

func (ws *WebSocket) reader() {
	defer func() {
		ws.conn.Close()
	}()
	for {
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ERROR] %v", err)
			}
			break
		}
		event, err := NewEventFromRaw(message)
		if err != nil {
			log.Printf("[ERROR | MSG] %v", err)
		} else {
			// log.Printf("[MSG] %v", event)
		}
		if action, ok := ws.events[strings.ToLower(event.Name)]; ok {
			action(event)
		}
	}
}

func (ws *WebSocket) writer() {
	for {
		select {
		case message, ok := <-ws.out:
			if !ok {
				ws.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := ws.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			w.Close()
		}
	}
}
