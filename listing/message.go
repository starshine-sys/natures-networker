package listing

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/natures-networker/common"
)

func (bot *Bot) post(ctx *bcr.Context) (err error) {
	id, err := strconv.ParseInt(ctx.Args[0], 10, 64)
	if err != nil {
		return err
	}

	listing, err := bot.DB.Listing(id)
	if err != nil {
		return ctx.SendfX("%v is not a valid listing ID!", id)
	}

	if listing.Name == "" {
		return ctx.SendX("The listing must have a name!")
	}
	if listing.Description == "" {
		return ctx.SendX("The listing must have a description!")
	}
	if listing.URL == "" {
		return ctx.SendX("The listing must have a URL!")
	}
	if len(listing.Representatives) == 0 {
		return ctx.SendfX("The listing must have at least one representative!")
	}

	channels, n := ctx.GreedyChannelParser(ctx.Args[1:])
	if n != -1 {
		return ctx.SendfX("Not all channels were parsed: channel %v was not found.", ctx.Args[n+1])
	}
	for _, ch := range channels {
		if ch.GuildID != ctx.Message.GuildID || (ch.Type != discord.GuildText && ch.Type != discord.GuildNews) {
			return ctx.SendfX("%v is not in this guild, or it's not a text channel!", ch.Mention())
		}
	}

	allowedMentions := make([]discord.UserID, 0, len(listing.Representatives))
	for _, u := range listing.Representatives {
		allowedMentions = append(allowedMentions, discord.UserID(u))
	}

	for _, ch := range channels {
		msg, err := ctx.State.SendMessageComplex(ch.ID, api.SendMessageData{
			Content: listing.Message(),
			AllowedMentions: &api.AllowedMentions{
				Parse: []api.AllowedMentionType{},
				Users: allowedMentions,
			},
		})
		if err != nil {
			return ctx.SendfX("Error sending listing message in %v: %v", ch.Mention(), err)
		}

		err = bot.DB.AddListingMessage(listing.ID, msg.ID, ch.ID)
		if err != nil {
			return ctx.SendfX("Error saving listing message in %v: %v", ch.Mention(), err)
		}
	}

	err = ctx.SendX("Listing posted!")
	if err != nil {
		return err
	}

	e := discord.Embed{
		Title:       "New server listed!",
		Description: "**" + listing.Name + "**\n\n**Categories**\n",
		Timestamp:   discord.NowTimestamp(),
		Color:       bcr.ColourGreen,
	}

	for _, ch := range channels {
		e.Description += ch.Mention() + "\n"
	}

	_, err = ctx.State.SendEmbeds(common.Conf.ListingLog, e)
	return err
}

func (bot *Bot) preview(ctx *bcr.Context) (err error) {
	id, err := strconv.ParseInt(ctx.Args[0], 10, 64)
	if err != nil {
		return err
	}

	listing, err := bot.DB.Listing(id)
	if err != nil {
		return ctx.SendfX("%v is not a valid listing ID!", id)
	}

	if listing.Name == "" {
		return ctx.SendX("The listing must have a name!")
	}
	if listing.Description == "" {
		return ctx.SendX("The listing must have a description!")
	}
	if listing.URL == "" {
		return ctx.SendX("The listing must have a URL!")
	}
	if len(listing.Representatives) == 0 {
		return ctx.SendfX("The listing must have at least one representative!")
	}

	allowedMentions := make([]discord.UserID, 0, len(listing.Representatives))
	for _, u := range listing.Representatives {
		allowedMentions = append(allowedMentions, discord.UserID(u))
	}

	_, err = ctx.State.SendMessageComplex(ctx.Message.ChannelID, api.SendMessageData{
		Content: listing.Message(),
		AllowedMentions: &api.AllowedMentions{
			Parse: []api.AllowedMentionType{},
			Users: allowedMentions,
		},
	})
	return
}

func (bot *Bot) update(ctx *bcr.Context) (err error) {
	id, err := strconv.ParseInt(ctx.Args[0], 10, 64)
	if err != nil {
		return err
	}

	listing, err := bot.DB.Listing(id)
	if err != nil {
		return ctx.SendfX("%v is not a valid listing ID!", id)
	}

	if listing.Name == "" {
		return ctx.SendX("The listing must have a name!")
	}
	if listing.Description == "" {
		return ctx.SendX("The listing must have a description!")
	}
	if listing.URL == "" {
		return ctx.SendX("The listing must have a URL!")
	}
	if len(listing.Representatives) == 0 {
		return ctx.SendfX("The listing must have at least one representative!")
	}

	allowedMentions := make([]discord.UserID, 0, len(listing.Representatives))
	for _, u := range listing.Representatives {
		allowedMentions = append(allowedMentions, discord.UserID(u))
	}

	msgs, err := bot.DB.ListingMessages(listing.ID)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		_, err = ctx.State.EditMessageComplex(msg.ChannelID, msg.ID, api.EditMessageData{
			Content: option.NewNullableString(listing.Message()),
			Embeds:  &[]discord.Embed{},
			AllowedMentions: &api.AllowedMentions{
				Parse: []api.AllowedMentionType{},
				Users: allowedMentions,
			},
		})
		if err != nil {
			return ctx.SendfX("Error updating message %v in %v: %v", msg.ID, msg.ChannelID.Mention(), err)
		}
	}

	return ctx.SendfX("Updated %v messages!", len(msgs))
}

func (bot *Bot) forceUpdate(ctx *bcr.Context) (err error) {
	var count int64
	_ = bot.DB.QueryRow(context.Background(), "select count(*) from listing_messages").Scan(&count)

	yes, _ := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
		Message:   fmt.Sprintf("Warning! This will edit %v messages. Do you want to continue?\n(Estimated time: %s)", count, time.Duration(count)*time.Second),
		YesPrompt: "Confirm",
	})
	if !yes {
		return ctx.SendX("Cancelled.")
	}

	listings, err := bot.DB.AllListings()
	if err != nil {
		return err
	}

	for _, listing := range listings {
		allowedMentions := make([]discord.UserID, 0, len(listing.Representatives))
		for _, u := range listing.Representatives {
			allowedMentions = append(allowedMentions, discord.UserID(u))
		}

		msgs, err := bot.DB.ListingMessages(listing.ID)
		if err != nil {
			return err
		}

		for _, msg := range msgs {
			_, err = ctx.State.EditMessageComplex(msg.ChannelID, msg.ID, api.EditMessageData{
				Content: option.NewNullableString(listing.Message()),
				Embeds:  &[]discord.Embed{},
				AllowedMentions: &api.AllowedMentions{
					Parse: []api.AllowedMentionType{},
					Users: allowedMentions,
				},
			})
			if err != nil {
				return ctx.SendfX("Error updating message %v in %v: %v", msg.ID, msg.ChannelID.Mention(), err)
			}
		}
	}

	return ctx.SendfX("Updated %v messages!", count)
}

func (bot *Bot) messageDelete(ev *gateway.MessageDeleteEvent) {
	_, err := bot.DB.Exec(context.Background(), "delete from listing_messages where id = $1", ev.ID)
	if err != nil {
		common.Log.Errorf("Error deleting listing message %v from database: %v", ev.ID, err)
	}
}

func (bot *Bot) bulkMessageDelete(ev *gateway.MessageDeleteBulkEvent) {
	for _, id := range ev.IDs {
		_, err := bot.DB.Exec(context.Background(), "delete from listing_messages where id = $1", id)
		if err != nil {
			common.Log.Errorf("Error deleting listing message %v from database: %v", id, err)
		}
	}
}

func (bot *Bot) channelDelete(ev *gateway.ChannelDeleteEvent) {
	_, err := bot.DB.Exec(context.Background(), "delete from listing_messages where channel_id = $1", ev.ID)
	if err != nil {
		common.Log.Errorf("Error deleting listing channel %v from database: %v", ev.ID, err)
	}
}
