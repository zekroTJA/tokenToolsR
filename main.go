package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	fport = flag.String("port", "80", "Port which will be used to expose the app's web interface")
	fversion = flag.Bool("version", false, "Display build version")

	appVersion string = "testing build"
	appCommit  string = "testing build"
	appDate	   string = "testing build"
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

	flag.Parse()

	if *fversion {
		fmt.Printf("tokenToolsR © 2018 zekro Development\n"+
				   "Version:   %s\n"+
				   "Commit:    %s\n"+
				   "Date:      %s\n",
				   appVersion, appCommit, appDate)
		return
	}

	http.Handle("/", http.FileServer(http.Dir("./assets")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if ws, err := NewWebSocket(w, r); err == nil {

			var discord *Discord
			var nguild int
			var err error

			ws.SetHandler("checkToken", func(e *Event) {
				discord, err = NewDiscord(e.Data.(string))
				if err != nil {
					sendInvalid(ws)
					return
				}
				info, err := discord.GetInfo()
				if err != nil {
					sendInvalid(ws)
					return
				}
				nguild = info.Guilds
				sendValid(ws, info)
			})

			ws.SetHandler("getGuildInfo", func(e *Event) {
				if discord == nil {
					log.Println(discord)
					return
				}

				guilds := make(chan *GuildInfo, nguild)

				err := discord.GetGuilds(guilds)
				if err != nil {
					log.Println("[ERR]", err)
					return
				}

				collectedGuilds := make([]*GuildInfo, nguild)
				counter := 0
				for {
					select {
					case g := <-guilds:
						collectedGuilds[counter] = g
						counter++
					}
					if counter == nguild {
						break
					}
				}

				go func() {
					ws.Out <- (&Event{
						Name: "guildInfo",
						Data: collectedGuilds,
					}).Raw()
				}()
			})

		} else {
			log.Println("[ERR] ", err)
		}
	})

	log.Println("[INFO] listening...")
	err := http.ListenAndServe(":"+*fport, nil)
	if err != nil {
		log.Fatal(err)
	}

}
