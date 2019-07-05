package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zekroTJA/tokenToolsR/internal/discord"
)

var statusMessages = map[int]string{
	200: "OK",
	400: "Bad request",
	427: "You are being rate limited",
}

var tokenChache = map[string]int{}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type TokenState struct {
	Valid         bool   `json:"valid"`
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Guilds        int    `json:"guilds"`
}

func SendResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response := &Response{
		Code:    code,
		Message: statusMessages[code],
		Data:    data,
	}
	bdata, _ := json.MarshalIndent(response, "", "  ")
	w.Write(bdata)
}

func InitApi(r *mux.Router, root string) {

	// GET check/:token
	r.HandleFunc(root+"/check/{token}", func(w http.ResponseWriter, r *http.Request) {
		if !CheckRatelimit(w, r) {
			return
		}

		params := mux.Vars(r)

		token, ok := params["token"]
		if !ok {
			SendResponse(w, 400, nil)
			return
		}

		if len(token) <= 40 {
			SendResponse(w, 200, &TokenState{Valid: false})
			return
		}

		discord, err := discord.NewDiscord(token)
		if err != nil {
			SendResponse(w, 500, err)
			return
		}
		info, err := discord.GetInfo()
		if err != nil {
			SendResponse(w, 200, &TokenState{Valid: false})
			return
		}
		tokenChache[token] = info.Guilds
		SendResponse(w, 200, &TokenState{
			Valid:         true,
			ID:            info.ID,
			Username:      info.Username,
			Discriminator: info.Discriminator,
			Avatar:        info.Avatar,
			Guilds:        info.Guilds,
		})

	}).Methods("GET")

	r.HandleFunc(root+"/guilds/{token}", func(w http.ResponseWriter, r *http.Request) {
		if !CheckRatelimit(w, r) {
			return
		}

		params := mux.Vars(r)

		token, ok := params["token"]
		if !ok {
			SendResponse(w, 400, nil)
			return
		}

		if len(token) <= 40 {
			SendResponse(w, 400, "token invalid")
			return
		}

		dc, err := discord.NewDiscord(token)
		if err != nil {
			SendResponse(w, 500, err)
			return
		}

		nguild, ok := tokenChache[token]
		if !ok {
			info, err := dc.GetInfo()
			if err != nil {
				SendResponse(w, 400, "token invalid")
				return
			}
			nguild = info.Guilds
			tokenChache[token] = info.Guilds
			<-time.After(300 * time.Millisecond)
		}

		guilds := make(chan *discord.GuildInfo, nguild)

		err = dc.GetGuilds(guilds)
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

		SendResponse(w, 200, collectedGuilds)

	}).Methods("GET")

}
