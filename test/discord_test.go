package test_test

import (
	"testing"
	"time"

	"github.com/zekroTJA/tokenToolsR/pkg/discord"
)

func TestNewDiscord(t *testing.T) {
	dc, err := discord.NewDiscord(apiToken)
	if err != nil {
		t.Fatal(err)
	}
	if dc == nil {
		t.Fatal("discord object was nil")
	}
}

func TestGetInfo(t *testing.T) {
	dc, err := discord.NewDiscord(apiToken)
	if err != nil {
		t.Fatal(err)
	}

	user, err := dc.GetInfo()
	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user was nil")
	}

	assertNotEqual(t, user.Avatar, "", "user.Avatar was emmpty")
	assertNotEqual(t, user.Discriminator, "", "user.Discriminator was emmpty")
	assertNotEqual(t, user.ID, "", "user.ID was emmpty")
	assertNotEqual(t, user.Guilds, 0, "user.Guilds was 0")
	assertNotEqual(t, user.Username, "", "user.Username was emmpty")
}

func TestGetUser(t *testing.T) {
	dc, err := discord.NewDiscord(apiToken)
	if err != nil {
		t.Fatal(err)
	}

	// For testing purposes, it does not matter if testing for
	// another user or the apps user itself. But it can be sure
	// that at least the apps user itself is present and can be get.
	user, err := dc.GetUser("@me")
	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user was nil")
	}

	assertNotEqual(t, user.Avatar, "", "user.Avatar was emmpty")
	assertNotEqual(t, user.Discriminator, "", "user.Discriminator was emmpty")
	assertNotEqual(t, user.ID, "", "user.ID was emmpty")
	assertNotEqual(t, user.Username, "", "user.Username was emmpty")
}

func TestGetGuilds(t *testing.T) {
	time.Sleep(10 * time.Second)
	dc, err := discord.NewDiscord(apiToken)
	if err != nil {
		t.Fatal(err)
	}

	me, err := dc.GetInfo()
	if err != nil {
		t.Fatal(err)
	}

	guildsResp := make([]*discord.GuildInfo, 0)
	guildsChan := make(chan *discord.GuildInfo)

	err = dc.GetGuilds(guildsChan)
	if err != nil {
		t.Fatal(err)
	}

	timeout := time.AfterFunc(1*time.Minute, func() {
		t.Fatal("GetGuilds timed out")
	})

	for i := 0; i < me.Guilds; i++ {
		guildsResp = append(guildsResp, <-guildsChan)
	}

	timeout.Stop()

	if len(guildsResp) < 1 {
		t.Fatal("guilds len was 0")
	}

	guild := guildsResp[0]
	assertNotEqual(t, guild.ID, "", "guild.ID was empty")
	assertNotEqual(t, guild.Members, "", "guild.Members was 0")
	assertNotEqual(t, guild.Name, "", "guild.Name was empty")
	assertNotEqual(t, guild.Owner, "", "guild.Owner was empty")
}
