package listing

import (
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/natures-networker/bot"
)

type Bot struct {
	*bot.Bot
}

func Init(b *bot.Bot) (s string, cmds []*bcr.Command) {
	s = "Listing commands"

	return s, cmds
}
