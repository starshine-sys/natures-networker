package meta

import (
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/natures-networker/bot"
)

type Bot struct {
	*bot.Bot
}

func Init(b *bot.Bot) (s string, cmds []*bcr.Command) {
	bot := &Bot{b}

	s = "Information commands"

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:    "about",
		Summary: "Show some information about " + bot.Router.Bot.Username,
		Command: bot.about,
	}))

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:    "help",
		Summary: "Show a list of commands",
		Command: bot.CommandList,
	}))

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:    "ping",
		Summary: "Show " + bot.Router.Bot.Username + "'s latency",
		Command: bot.ping,
	}))

	return s, cmds
}
