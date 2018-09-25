package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	APIROOT = "https://discordapp.com/api"
)

type Discord struct {
	token  string
	client *http.Client
	header http.Header
}

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Guilds        int    `json:"guilds"`
}

type Guild struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Owner   string `json:"owner_id"`
	Members int    `json:"member_count"`
}

type GuildInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Owner   *User  `json:"owner"`
	Members int    `json:"members"`
}

func NewDiscord(token string) (*Discord, error) {
	client := &http.Client{}
	header := http.Header{}
	header.Add("User-Agent", "DiscordBot (https://github.com/zekrotja, 0.1.0)")
	header.Add("Authorization", "Bot "+token)
	discord := &Discord{
		token:  "Bot " + token,
		client: client,
		header: header,
	}
	return discord, nil
}

func (d *Discord) request(method, endpoint string, data []byte, output interface{}) error {
	req, err := http.NewRequest(method, APIROOT+"/"+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header = d.header

	res, err := d.client.Do(req)
	if err != nil {
		return err
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	apiError := &ApiError{Code: -1}
	err = json.Unmarshal(resData, apiError)
	if err == nil && apiError.Code != -1 && apiError.Message != "" {
		return errors.New("API ERROR: " + apiError.Message)
	}

	err = json.Unmarshal(resData, output)
	return err
}

func (d *Discord) GetInfo() (*User, error) {
	user := new(User)
	err := d.request("GET", "users/@me", nil, user)
	if err != nil {
		return nil, err
	}
	guilds := make([]*struct {
		ID string `json:"id"`
	}, 0)
	d.request("GET", "users/@me/guilds", nil, &guilds)
	user.Guilds = len(guilds)
	return user, nil
}

func (d *Discord) GetGuilds(guilds chan *GuildInfo) error {
	guildsResp := make([]*Guild, 0)
	err := d.request("GET", "users/@me/guilds", nil, &guildsResp)
	if err != nil {
		return err
	}

	for _, gld := range guildsResp {
		go func(g *Guild) {
			guild := new(Guild)
			err = d.request("GET", "guilds/"+g.ID, nil, guild)
			if err == nil {
				ownerid := guild.Owner
				owner := new(User)
				d.request("GET", "users/"+ownerid, nil, owner)

				guildMembers := make([]*struct {
					ID string `json:"id"`
				}, 0)
				d.request("GET", "guilds/"+g.ID+"/members?limit=1000", nil, &guildMembers)

				guildInfo := &GuildInfo{
					ID:      guild.ID,
					Name:    guild.Name,
					Owner:   owner,
					Members: len(guildMembers),
				}
				guilds <- guildInfo
			}
		}(gld)
	}

	return nil
}
