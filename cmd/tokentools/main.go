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

	InitApi(router, "/api")

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
