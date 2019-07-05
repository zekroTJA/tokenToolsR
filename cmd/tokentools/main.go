package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zekroTJA/tokenToolsR/internal/api"
	"github.com/zekroTJA/tokenToolsR/internal/discord"
	"github.com/zekroTJA/tokenToolsR/internal/static"
	"github.com/zekroTJA/tokenToolsR/internal/ws"
)

var (
	fport     = flag.String("port", "80", "Port which will be used to expose the app's web interface")
	fversion  = flag.Bool("version", false, "Display build version")
	fcertfile = flag.String("tls-cert", "", "The TLS cert file")
	fkeyfile  = flag.String("tls-key", "", "The TLS key file")
)

func sendInvalid(w *ws.WebSocket) {
	go func() {
		w.Out <- (&ws.Event{
			Name: "tokenInvalid",
			Data: nil,
		}).Raw()
	}()
}

func sendValid(w *ws.WebSocket, info *discord.User) {
	go func() {
		w.Out <- (&ws.Event{
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
			static.AppVersion, static.AppCommit, static.AppDate)
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.New("index.html")
		t, _ = t.ParseFiles("./web/views/index.html")
		t.Execute(w, struct {
			VERSION string
			COMMIT  string
			DATE    string
		}{
			static.AppVersion, static.AppCommit, static.AppDate,
		})
	})

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if w, err := ws.NewWebSocket(w, r); err == nil {

			var dc *discord.Discord
			var nguild int
			var err error

			w.SetHandler("checkToken", func(e *ws.Event) {
				dc, err = discord.NewDiscord(e.Data.(string))
				if err != nil {
					sendInvalid(w)
					return
				}
				info, err := dc.GetInfo()
				if err != nil {
					sendInvalid(w)
					return
				}
				nguild = info.Guilds
				sendValid(w, info)
			})

			w.SetHandler("getGuildInfo", func(e *ws.Event) {
				if dc == nil {
					log.Println(dc)
					return
				}

				guilds := make(chan *discord.GuildInfo, nguild)

				err := dc.GetGuilds(guilds)
				if err != nil {
					log.Println("[ERR]", err)
					return
				}

				collectedGuilds := make([]*discord.GuildInfo, nguild)
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
					w.Out <- (&ws.Event{
						Name: "guildInfo",
						Data: collectedGuilds,
					}).Raw()
				}()
			})

			w.SetHandler("getUserInfo", func(e *ws.Event) {
				if dc == nil {
					return
				}

				uid := e.Data.(string)
				user, err := dc.GetUser(uid)
				if err == nil {
					go func() {
						w.Out <- (&ws.Event{
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

	api.InitApi(router, "/api")

	http.Handle("/", router)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./web/assets"))))

	log.Println("[INFO] listening...")
	var err error
	if *fcertfile != "" && *fkeyfile != "" {
		err = http.ListenAndServeTLS(":"+*fport, *fcertfile, *fkeyfile, nil)
	} else {
		err = http.ListenAndServe(":"+*fport, nil)
	}
	if err != nil {
		log.Fatal(err)
	}

}
