package main

import (
	"log"
	"fmt"
	"net/http"
)

func sendInvalid(ws *WebSocket) {
	go func() {
		ws.Out <- (&Event{
			Name: "tokenInvalid",
			Data: nil,
		}).Raw()
	}()
}

func sendValid(ws *WebSocket, info *User) {
	info.Avatar = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", info.ID, info.Avatar)
	go func() {
		ws.Out <- (&Event{
			Name: "tokenValid",
			Data: info,
		}).Raw()
	}()
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("./assets")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if ws, err := NewWebSocket(w, r); err == nil {
			ws.SetHandler("checkToken", func(e *Event) {
				discord, err := NewDiscord(e.Data.(string))
				if err != nil {
					sendInvalid(ws)
				}
				info, err := discord.GetInfo()
				if err != nil {
					sendInvalid(ws)
				}
				sendValid(ws, info)
			})
		} else {
			log.Println("[ERR] ", err)
		}
	})

	log.Println("[INFO] listening...")
	err := http.ListenAndServe(":1337", nil)
	if err != nil {
		log.Fatal(err)
	}

}