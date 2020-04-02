package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zekroTJA/tokenToolsR/internal/api"
	"github.com/zekroTJA/tokenToolsR/internal/discord"
	"github.com/zekroTJA/tokenToolsR/internal/spa"
	"github.com/zekroTJA/tokenToolsR/internal/static"
	"github.com/zekroTJA/tokenToolsR/internal/ws"
)

var (
	fport     = flag.String("port", "80", "Port which will be used to expose the app's web interface")
	fversion  = flag.Bool("version", false, "Display build version")
	fcertfile = flag.String("tls-cert", "", "The TLS cert file")
	fkeyfile  = flag.String("tls-key", "", "The TLS key file")
	ftls      = flag.Bool("tls", false, "Wether or not to enable TLS")
	fwebdir   = flag.String("web", "./web/build", "static web files location")
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

func sendError(w *ws.WebSocket, err error) {
	go func() {
		w.Out <- (&ws.Event{
			Name: "error",
			Data: err.Error(),
		}).Raw()
	}()
}

func main() {

	flag.Parse()

	if *fversion {
		fmt.Printf("tokenToolsR © 2020 zekro Development\n"+
			"Version:   %s\n"+
			"Commit:    %s\n"+
			"Date:      %s\n",
			static.AppVersion, static.AppCommit, static.AppDate)
		return
	}

	router := mux.NewRouter()

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
				token := e.Data.(string)
				if token != "" {
					if dc, err = discord.NewDiscord(token); err != nil {
						sendError(w, err)
						return
					}
					info, err := dc.GetInfo()
					if err != nil {
						sendError(w, err)
						return
					}
					nguild = info.Guilds
				}

				if dc == nil {
					return
				}

				guilds := make(chan *discord.GuildInfo, nguild)

				time.Sleep(750 * time.Millisecond)

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

	spaHandler := spa.NewSPA(*fwebdir, "index.html")
	router.PathPrefix("/").Handler(spaHandler)

	http.Handle("/", router)

	log.Println("[INFO] listening...")
	var err error
	if *ftls && *fcertfile != "" && *fkeyfile != "" {
		err = http.ListenAndServeTLS(":"+*fport, *fcertfile, *fkeyfile, nil)
	} else {
		err = http.ListenAndServe(":"+*fport, nil)
	}
	if err != nil {
		log.Fatal(err)
	}

}
