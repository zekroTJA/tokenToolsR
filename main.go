package main

import (
	"html/template"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.New("index.html")
		t, _ = t.ParseFiles("./views/index.html")
		t.Execute(w, struct{
			VERSION string
			COMMIT  string
			DATE	string
		}{
			appVersion, appCommit, appDate,
		})
	})

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

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

			ws.SetHandler("getUserInfo", func(e *Event) {
				if discord == nil {
					return
				}

				uid := e.Data.(string)
				user, err := discord.GetUser(uid)
				if err == nil {
					go func() {
						ws.Out <- (&Event{
							Name: "userInfo",
							Data: user,
						}).Raw()
					}()
				} else {
					fmt.Println("[ERR]", err)
				}
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
