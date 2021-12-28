package listing

import (
	"regexp"

	"github.com/starshine-sys/bcr"
)

var inviteRegexp = regexp.MustCompile(`(https?:\/\/)?(discord\.gg|discord(app)?\.com\/invite)\/(\w+)`)

func (bot *Bot) addServer(ctx *bcr.Context) (err error) {
	matches := inviteRegexp.FindStringSubmatch(ctx.RawArgs)
	if len(matches) != 5 {
		return ctx.SendfX("``%v`` is not a valid server invite link :(\n(Format: `https://discord.gg/<code>`, `https://discordapp.com/invite/<code>`, `https://discord.com/invite/<code>`)", bcr.EscapeBackticks(ctx.RawArgs))
	}
	code := matches[4]

	inv, err := ctx.State.Invite(code)
	if err != nil {
		return ctx.SendfX("That doesn't seem to be a valid invite link.")
	}

	l, err := bot.DB.AddListing(inv.Guild.Name, "https://discord.gg/"+code, true)
	if err != nil {
		return err
	}

	return ctx.SendfX("**%v** (invite: <https://discord.gg/%v>) added with ID **%v**!", inv.Guild.Name, inv.Code, l.ID)
}

func (bot *Bot) addRaw(ctx *bcr.Context) (err error) {
	name := ctx.Args[0]
	url := ctx.Args[1]

	l, err := bot.DB.AddListing(name, url, false)
	if err != nil {
		return err
	}

	return ctx.SendfX("**%v** (link: <%v>) added with ID **%v**!", name, url, l.ID)
}
