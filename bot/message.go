package bot

import (
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/starshine-sys/natures-networker/common"
)

func (bot *Bot) messageCreate(m *gateway.MessageCreateEvent) {
	// if the author is a bot, return
	if m.Author.Bot {
		return
	}

	// if the message does not start with any of the bot's prefixes (including mentions), return
	if !bot.Router.MatchPrefix(m.Message) {
		return
	}

	// get the context
	ctx, err := bot.Router.NewContext(m)
	if err != nil {
		ctx.Router.Logger.Error("getting context: %v", err)
		return
	}

	err = bot.Router.Execute(ctx)
	if err != nil {
		switch err.(type) {
		case httputil.HTTPError, *httputil.HTTPError:
			common.Log.Warnf("HTTP error in command: %v", err)
		default:
			common.Log.Errorf("Error in command %v: %v", ctx.Command, err)
			err = ctx.SendfX("Internal error occurred.")
			if err != nil {
				common.Log.Warnf("HTTP error in command: %v", err)
			}
		}
	}
}
