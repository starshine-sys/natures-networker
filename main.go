package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/ws"
	"github.com/starshine-sys/natures-networker/ban"
	"github.com/starshine-sys/natures-networker/bot"
	"github.com/starshine-sys/natures-networker/common"
	"github.com/starshine-sys/natures-networker/listing"
	"github.com/starshine-sys/natures-networker/meta"
)

func main() {
	ws.WSDebug = common.Log.Debug

	common.Log.Infof("Starting Nature's Networker version %v", common.Version)

	bot, err := bot.New()
	if err != nil {
		common.Log.Fatalf("Error creating bot: %v", err)
	}

	state, _ := bot.Router.StateFromGuildID(0)
	botUser, err := state.Me()
	if err != nil {
		common.Log.Fatal(err)
	}
	bot.Router.Bot = botUser

	bot.Add(meta.Init)
	bot.Add(listing.Init)
	bot.Add(ban.Init)

	bot.Router.AddHandler(ready(state))

	err = bot.Start(context.Background())
	if err != nil {
		common.Log.Fatalf("Error opening connection to Discord: %v", err)
	}

	// Defer this to make sure that things are always cleanly shutdown even in the event of a crash
	defer func() {
		bot.Router.ShardManager.Close()
		common.Log.Infof("Disconnected from Discord")
		bot.DB.Close()
		common.Log.Infof("Database connection closed")
	}()

	common.Log.Info("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")
	common.Log.Infof("User: %v (%v)", botUser.Tag(), botUser.ID)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	<-ctx.Done()

	common.Log.Infof("Interrupt signal received. Shutting down...")
}

func ready(s *state.State) func(*gateway.ReadyEvent) {
	return func(*gateway.ReadyEvent) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.Gateway().Send(ctx, &gateway.UpdatePresenceCommand{
			Activities: []discord.Activity{{
				Name: common.Conf.Prefix + "help",
			}},
			Status: discord.OnlineStatus,
		})
		if err != nil {
			common.Log.Errorf("Error setting status: %v", err)
		}
	}
}
