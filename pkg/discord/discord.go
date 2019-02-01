package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	apiRoot   = "https://discordapp.com/api"
	defAvatar = "https://discordapp.com/assets/0e291f67c9274a1abdddeb3fd919cbaa.png"
)

// Discord contains the Discord APP token,
// the HTTP Client and the HTTP request
// headers.
type Discord struct {
	token  string
	client *http.Client
	header http.Header
}

// APIError contains the error code
// and the error message of an error
// returned from the Discord API.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// User contains a Discord API user object:
// https://discordapp.com/developers/docs/resources/user
type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Guilds        int    `json:"guilds"`
}

// Guild contains a Discord API guild object:
// https://discordapp.com/developers/docs/resources/guild
type Guild struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Owner   string `json:"owner_id"`
	Members int    `json:"member_count"`
}

// GuildInfo contains the ID, name and the owner ID of
// the guild as well as the number of members on the guild.
type GuildInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	Members int    `json:"members"`
}

// NewDiscord creates a new instance of Discord.
// The passed token will beset as Authentication
// header and the HTTP client will be created.
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
	req, err := http.NewRequest(method, apiRoot+"/"+endpoint, bytes.NewBuffer(data))
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
	apiError := &APIError{Code: -1}
	err = json.Unmarshal(resData, apiError)
	if err == nil && apiError.Code != -1 && apiError.Message != "" {
		return errors.New("API ERROR: " + apiError.Message)
	}

	fmt.Println(">>>>>>>>>>>>\n", method, "\n", endpoint, "\n", string(resData), "\n<<<<<<<<<<<<")

	err = json.Unmarshal(resData, output)
	return err
}

// GetInfo returns the user object of the current
// user impersonated by the token (users/@me).
func (d *Discord) GetInfo() (*User, error) {
	user := new(User)
	err := d.request("GET", "users/@me", nil, user)
	if err != nil {
		return nil, err
	}
	if user.Avatar == "" {
		user.Avatar = defAvatar
	} else {
		user.Avatar = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", user.ID, user.Avatar)
	}
	guilds := make([]*struct {
		ID string `json:"id"`
	}, 0)
	d.request("GET", "users/@me/guilds", nil, &guilds)
	user.Guilds = len(guilds)
	return user, nil
}

// GetUser returns the user object of a user specified
// by their ID (users/:ID).
func (d *Discord) GetUser(uid string) (*User, error) {
	user := new(User)
	err := d.request("GET", "users/"+uid, nil, user)
	if user.Avatar == "" {
		user.Avatar = defAvatar
	} else {
		user.Avatar = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", user.ID, user.Avatar)
	}
	return user, err
}

// GetGuilds creates a request for every guild asyncronously
// which results will be sent into the passed GuildInfo channel.
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

				guildMembers := make([]*struct {
					ID string `json:"id"`
				}, 0)
				d.request("GET", "guilds/"+g.ID+"/members?limit=1000", nil, &guildMembers)

				guildInfo := &GuildInfo{
					ID:      guild.ID,
					Name:    guild.Name,
					Owner:   ownerid,
					Members: len(guildMembers),
				}
				guilds <- guildInfo
			}
		}(gld)
	}

	return nil
}
