package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/zekroTJA/tokenToolsR/internal/api"
	"github.com/zekroTJA/tokenToolsR/internal/websocket"
	"github.com/zekroTJA/tokenToolsR/pkg/discord"
)

var (
	fport     = flag.String("p", "80", "Port which will be used to expose the app's web interface")
	fversion  = flag.Bool("v", false, "Display build version")
	fcertfile = flag.String("c", "", "The cert config file with fiel link to cert file and key file (see example cert config)")

	appVersion = "testing build"
	appCommit  = "testing build"
	appDate    = "testing build"
)

type Cert struct {
	CertFile string `json:"cert"`
	KeyFile  string `json:"key"`
}

func sendInvalid(ws *websocket.WebSocket) {
	go func() {
		ws.Out <- (&websocket.Event{
			Name: "tokenInvalid",
			Data: nil,
		}).Raw()
	}()
}

func sendValid(ws *websocket.WebSocket, info *discord.User) {
	go func() {
		ws.Out <- (&websocket.Event{
			Name: "tokenValid",
			Data: info,
		}).Raw()
	}()
}

func main() {

	flag.Parse()

	cert := new(Cert)
	if *fcertfile != "" {
		filehandle, err := os.Open(*fcertfile)
		if err != nil {
			log.Fatal("[FATAL] ", err)
		}
		decoder := json.NewDecoder(filehandle)
		err = decoder.Decode(cert)
		if err != nil {
			log.Fatal("[FATAL] ", err)
		}
	}

	if *fversion {
		fmt.Printf("tokenToolsR © 2018 zekro Development\n"+
			"Version:   %s\n"+
			"Commit:    %s\n"+
			"Date:      %s\n",
			appVersion, appCommit, appDate)
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.New("index.html")
		t, _ = t.ParseFiles("./assets/views/index.html")
		t.Execute(w, struct {
			VERSION string
			COMMIT  string
			DATE    string
		}{
			appVersion, appCommit, appDate,
		})
	})

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if ws, err := websocket.NewWebSocket(w, r); err == nil {

			var dc *discord.Discord
			var nguild int
			var err error

			ws.SetHandler("checkToken", func(e *websocket.Event) {
				dc, err = discord.NewDiscord(e.Data.(string))
				if err != nil {
					sendInvalid(ws)
					return
				}
				info, err := dc.GetInfo()
				if err != nil {
					sendInvalid(ws)
					return
				}
				nguild = info.Guilds
				sendValid(ws, info)
			})

			ws.SetHandler("getGuildInfo", func(e *websocket.Event) {
				if dc == nil {
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
					ws.Out <- (&websocket.Event{
						Name: "guildInfo",
						Data: collectedGuilds,
					}).Raw()
				}()
			})

			ws.SetHandler("getUserInfo", func(e *websocket.Event) {
				if dc == nil {
					return
				}

				uid := e.Data.(string)
				user, err := dc.GetUser(uid)
				if err == nil {
					go func() {
						ws.Out <- (&websocket.Event{
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
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./assets/static"))))

	log.Println("[INFO] listening...")
	var err error
	if *fcertfile != "" {
		err = http.ListenAndServeTLS(":"+*fport, cert.CertFile, cert.KeyFile, nil)
	} else {
		err = http.ListenAndServe(":"+*fport, nil)
	}
	if err != nil {
		log.Fatal(err)
	}

}
