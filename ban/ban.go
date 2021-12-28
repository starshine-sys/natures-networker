package ban

import (
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/natures-networker/common"
)

func (bot *Bot) ban(ctx *bcr.Context) (err error) {
	args := strings.Split(ctx.RawArgs, "\n")
	if len(args) < 2 {
		return ctx.SendX("You must give both user(s) to ban and a reason!")
	}

	users, _ := ctx.GreedyUserParser(strings.Fields(args[0]))
	if len(users) == 0 {
		return ctx.SendX("You must give at least one user to ban!")
	}

	reason := strings.TrimSpace(strings.Join(args[1:], "\n"))
	if reason == "" {
		return ctx.SendX("Reason cannot be empty!")
	}
	if len(reason) > 1000 {
		return ctx.SendX("Reason can only be 1000 characters in length.")
	}

	usernames := ""
	for _, u := range users {
		usernames += u.Tag() + " " + u.ID.String() + "\n"
	}

	yes, _ := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
		Message: "Are you sure you want to ban the following users?",
		Embeds: []discord.Embed{{
			Description: usernames,
			Color:       bcr.ColourRed,
			Fields: []discord.EmbedField{{
				Name:  "Reason",
				Value: reason,
			}},
		}},
		YesPrompt: "Ban",
		YesStyle:  discord.DangerButtonStyle(),
		Timeout:   5 * time.Minute,
	})
	if !yes {
		return ctx.SendX("Cancelled.")
	}

	log, err := bot.DB.CreateBan(users, reason)
	if err != nil {
		return err
	}

	msg, err := ctx.State.SendEmbeds(common.Conf.BanLog, log.Embed(ctx))
	if err != nil {
		return err
	}

	err = bot.DB.SetBanMessage(log.ID, msg.ID)
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(users))
	for _, u := range users {
		ids = append(ids, u.ID.String())
	}

	return ctx.SendfX("Ban log (ID: %v) sent!\nIf they haven't been banned yet, make sure to ban the following IDs: %v", log.ID, strings.Join(ids, ", "))
}
