package bot

import (
	"sort"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/bcr/bot"
	"github.com/starshine-sys/natures-networker/common"
	"github.com/starshine-sys/natures-networker/db"
)

const Colour = bcr.ColourBlurple

type Bot struct {
	*bot.Bot
	Colour discord.Color

	DB *db.DB

	RequireStaff, RequireHelper bcr.CustomPerms
}

func New() (*Bot, error) {
	db, err := db.New(common.Conf.DatabaseURL)
	if err != nil {
		return nil, err
	}

	bcrbot, err := bot.New(common.Conf.Token)
	if err != nil {
		return nil, err
	}
	bcrbot.Prefix(common.Conf.Prefix)
	bcrbot.Owner(common.Conf.Owner)
	bcrbot.Router.EmbedColor = Colour

	bot := &Bot{Bot: bcrbot, DB: db, Colour: Colour}

	bot.Router.AddHandler(bot.messageCreate)
	bot.Router.AddHandler(bot.Router.InteractionCreate)

	bot.RequireStaff = bot.Router.RequireRole(
		"Staff", common.Conf.StaffRole,
	)
	bot.RequireHelper = bot.Router.RequireRole(
		"Helper", common.Conf.HelperRole, common.Conf.StaffRole,
	)

	return bot, nil
}

// Add adds a module to the bot
func (bot *Bot) Add(f func(*Bot) (string, []*bcr.Command)) {
	m, c := f(bot)

	// sort the list of commands
	sort.Sort(bcr.Commands(c))

	// add the module
	bot.Bot.Modules = append(bot.Bot.Modules, &botModule{
		name:     m,
		commands: c,
	})
}

type botModule struct {
	name     string
	commands bcr.Commands
}

// String returns the module's name
func (b botModule) String() string {
	return b.name
}

// Commands returns a list of commands
func (b *botModule) Commands() []*bcr.Command {
	return b.commands
}
