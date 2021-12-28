package listing

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/natures-networker/common"
)

func (bot *Bot) delist(ctx *bcr.Context) (err error) {
	id, err := strconv.ParseInt(ctx.Args[0], 10, 64)
	if err != nil {
		return err
	}
	reason := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))

	listing, err := bot.DB.Listing(id)
	if err != nil {
		return ctx.SendfX("%v is not a valid listing ID!", id)
	}

	yes, _ := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
		Message: fmt.Sprintf("Are you sure you want to delist %v?", listing.Name),
		Embeds: []discord.Embed{{
			Title:       "Reason",
			Description: reason,
			Color:       bot.Colour,
		}},
		YesPrompt: "Delist",
		YesStyle:  discord.DangerButtonStyle(),
		Timeout:   5 * time.Minute,
	})
	if !yes {
		return ctx.SendX("Cancelled.")
	}

	msgs, err := bot.DB.ListingMessages(listing.ID)
	if err != nil {
		return err
	}

	var failed string

	for _, msg := range msgs {
		err = ctx.State.DeleteMessage(msg.ChannelID, msg.ID, "Server delisted")
		if err != nil {
			failed += fmt.Sprintf("%v in %v\n", msg.ID, msg.ChannelID.Mention())
		}
	}

	s := "**" + listing.Name + "** delisted!"
	if failed != "" {
		s += "\nFailed to delete the following messages:\n" + failed
	}

	err = ctx.SendX(s)
	if err != nil {
		return err
	}

	bot.checkRepRole(ctx, listing.Representatives...)

	e := discord.Embed{
		Title:       "Server delisted",
		Description: "**" + listing.Name + "**",
		Fields: []discord.EmbedField{{
			Name:  "Reason",
			Value: reason,
		}},
		Timestamp: discord.NowTimestamp(),
		Color:     bcr.ColourRed,
	}

	_, err = ctx.State.SendEmbeds(common.Conf.ListingLog, e)
	return
}
