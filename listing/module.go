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

	bot := &Bot{b}

	// delete handlers
	bot.Router.AddHandler(bot.messageDelete)
	bot.Router.AddHandler(bot.bulkMessageDelete)
	bot.Router.AddHandler(bot.channelDelete)

	add := bot.Router.AddCommand(&bcr.Command{
		Name:              "add",
		Summary:           "Add a server",
		Usage:             "<invite link>",
		Args:              bcr.MinArgs(1),
		CustomPermissions: bot.RequireHelper,
		Command:           bot.addServer,
	})

	add.AddSubcommand(&bcr.Command{
		Name:              "raw",
		Summary:           "Add a listing without parsing an invite link",
		Usage:             "<name> <link>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: bot.RequireHelper,
		Command:           bot.addRaw,
	})

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:              "description",
		Summary:           "Set a listing's description",
		Usage:             "<id> <description>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: bot.RequireHelper,
		Command:           bot.description,
	}))

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:              "rename",
		Summary:           "Set a listing's name",
		Usage:             "<id> <name>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: bot.RequireHelper,
		Command:           bot.rename,
	}))

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:              "reps",
		Aliases:           []string{"representatives"},
		Summary:           "Set a listing's representatives",
		Usage:             "<id> <reps...>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: bot.RequireHelper,
		Command:           bot.reps,
	}))

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:              "link",
		Summary:           "Set a listing's link",
		Usage:             "<id> <link> <embed>",
		Args:              bcr.MinArgs(3),
		CustomPermissions: bot.RequireHelper,
		Command:           bot.link,
	}))

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:              "post",
		Summary:           "Post a listing",
		Usage:             "<id> <channels...>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: bot.RequireHelper,
		Command:           bot.post,
	}))

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:              "update",
		Summary:           "Update a listing's messages",
		Usage:             "<id>",
		Args:              bcr.MinArgs(1),
		CustomPermissions: bot.RequireHelper,
		Command:           bot.update,
	}))

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:      "force-update",
		Summary:   "Update **all** listings' messages",
		OwnerOnly: true,
		Hidden:    true,
		Command:   bot.forceUpdate,
	}))

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:              "preview",
		Summary:           "Preview a listing",
		Usage:             "<id>",
		Args:              bcr.MinArgs(1),
		CustomPermissions: bot.RequireHelper,
		Command:           bot.preview,
	}))

	cmds = append(cmds, bot.Router.AddCommand(&bcr.Command{
		Name:              "delist",
		Summary:           "Delist a listing",
		Usage:             "<id> <reason>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: bot.RequireHelper,
		Command:           bot.delist,
	}))

	find := bot.Router.AddCommand(&bcr.Command{
		Name:              "find",
		Summary:           "Find the servers a user is representative for",
		Usage:             "<user>",
		Args:              bcr.MinArgs(1),
		CustomPermissions: bot.RequireHelper,
		Command:           bot.findByRep,
	})

	find.AddSubcommand(&bcr.Command{
		Name:              "name",
		Summary:           "Find a server by name",
		Usage:             "<query>",
		Args:              bcr.MinArgs(1),
		CustomPermissions: bot.RequireHelper,
		Command:           bot.findByName,
	})

	return s, append(cmds, add)
}
