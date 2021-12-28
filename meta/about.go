package meta

import (
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/natures-networker/common"
)

func (bot *Bot) about(ctx *bcr.Context) (err error) {
	g, err := ctx.State.Guild(common.Conf.GuildID)
	if err != nil {
		return err
	}

	e := discord.Embed{
		Title: "About " + bot.Router.Bot.Username,
		Color: bot.Colour,
		Thumbnail: &discord.EmbedThumbnail{
			URL: g.IconURL(),
		},
		Description: "[GitHub](https://github.com/starshine-sys/natures-networker)",
	}

	return ctx.SendX("", e)
}

func (bot *Bot) ping(ctx *bcr.Context) (err error) {
	t := time.Now()

	err = ctx.SendX("...")
	if err != nil {
		return err
	}

	latency := time.Since(t).Round(time.Millisecond)

	// this will return 0ms in the first minute after the bot is restarted
	// can't do much about that though
	heartbeat := ctx.Session().Gateway().EchoBeat().Sub(ctx.Session().Gateway().SentBeat()).Round(time.Millisecond)

	_, err = ctx.EditOriginal(api.EditInteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("Pong! Heartbeat: %s | Message: %s", heartbeat, latency)),
	})
	return err
}
