package ban

import (
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/natures-networker/bot"
)

type Bot struct {
	*bot.Bot
}

func Init(b *bot.Bot) (s string, cmds []*bcr.Command) {
	s = "Ban commands"

	bot := &Bot{b}

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:              "ban",
		Summary:           "Post a ban log for the given users",
		Usage:             "<users...>\n<reason>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: bot.RequireStaff,
		Command:           bot.ban,
	}))

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:              "evidence",
		Summary:           "Update the evidence for a ban",
		Usage:             "<id> <evidence>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: bot.RequireStaff,
		Command:           bot.evidence,
	}))

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:              "reason",
		Summary:           "Update the reason for a ban",
		Usage:             "<id> <reason>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: bot.RequireStaff,
		Command:           bot.reason,
	}))

	return s, cmds
}
