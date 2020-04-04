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
	faddr     = flag.String("addr", "localhost:80", "Port which will be used to expose the app's web interface")
	fversion  = flag.Bool("version", false, "Display build version")
	fcertfile = flag.String("tls-cert", "", "The TLS cert file")
	fkeyfile  = flag.String("tls-key", "", "The TLS key file")
	ftls      = flag.Bool("tls", false, "Wether or not to enable TLS")
	fwebdir   = flag.String("web", "./web/build", "static web files location")
)

type GuildInfoTile struct {
	Guild *discord.GuildInfo `json:"guild"`
	N     int                `json:"n"`
	NMax  int                `json:"nmax"`
}

func sendInvalid(w *ws.WebSocket, cid int) {
	w.Send("invalid", cid, nil)
}

func sendValid(w *ws.WebSocket, cid int, info *discord.User) {
	w.Send("valid", cid, info)
}

func sendError(w *ws.WebSocket, cid int, err error) {
	w.Send("error", cid, err.Error())
}

func sendErrorS(w *ws.WebSocket, cid int, err string, v ...interface{}) {
	sendError(w, cid, fmt.Errorf(err, v))
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
		if s, err := ws.NewWebSocket(w, r); err == nil {

			var dc *discord.Discord
			var nguilds = -1

			s.SetHandler("init", func(e *ws.Event) {
				token, _ := e.Data.(string)
				if token == "" {
					sendErrorS(s, e.CID, "invalid token payload")
					return
				}

				dc = discord.NewDiscord(token)
				nguilds = -1
			})

			s.SetHandler("check", func(e *ws.Event) {
				if dc == nil {
					sendErrorS(s, e.CID, "not initialized")
					return
				}

				info, err := dc.GetInfo()
				if err != nil {
					sendInvalid(s, e.CID)
					return
				}
				nguilds = info.Guilds
				sendValid(s, e.CID, info)
			})

			s.SetHandler("guildinfo", func(e *ws.Event) {
				if dc == nil {
					sendErrorS(s, e.CID, "not initialized")
					return
				}

				if nguilds < 0 {
					info, err := dc.GetInfo()
					if err != nil {
						sendError(s, e.CID, err)
						return
					}
					nguilds = info.Guilds
					time.Sleep(1000 * time.Millisecond)
				}

				if nguilds == 0 {
					return
				}

				guilds := make(chan *discord.GuildInfo, nguilds)

				err := dc.GetGuilds(guilds)
				if err != nil {
					sendError(s, e.CID, err)
					log.Println("[ERR]", err)
					return
				}

				counter := 0
				n := nguilds
				for {
					select {
					case g := <-guilds:
						counter++
						go func() {
							s.Send("guildinfo", e.CID, &GuildInfoTile{
								Guild: g,
								N:     counter,
								NMax:  n,
							})
						}()
					}
					if counter == n {
						break
					}
				}
			})

			s.SetHandler("userinfo", func(e *ws.Event) {
				if dc == nil {
					sendErrorS(s, e.CID, "not initialized")
					return
				}

				uid, _ := e.Data.(string)
				if uid == "" {
					sendErrorS(s, e.CID, "invalid user ID payload")
					return
				}

				user, err := dc.GetUser(uid)
				if err != nil {
					sendError(s, e.CID, err)
					return
				}

				s.Send("userinfo", e.CID, user)
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
		err = http.ListenAndServeTLS(*faddr, *fcertfile, *fkeyfile, nil)
	} else {
		err = http.ListenAndServe(*faddr, nil)
	}
	if err != nil {
		log.Fatal(err)
	}

}
