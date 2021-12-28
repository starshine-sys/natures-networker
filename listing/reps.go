package listing

import (
	"context"
	"strconv"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/natures-networker/common"
)

func (bot *Bot) reps(ctx *bcr.Context) (err error) {
	id, err := strconv.ParseInt(ctx.Args[0], 10, 64)
	if err != nil {
		return err
	}

	listing, err := bot.DB.Listing(id)
	if err != nil {
		return ctx.SendfX("%v is not a valid listing ID!", id)
	}

	users, n := ctx.GreedyUserParser(ctx.Args[1:])
	if n != -1 {
		return ctx.SendfX("Not all users were parsed: user %v was not found.", ctx.Args[n+1])
	}

	reps := make([]uint64, 0, len(users))
	for _, u := range users {
		reps = append(reps, uint64(u.ID))
	}

	_, err = bot.DB.Exec(context.Background(), "update listings set representatives = $1 where id = $2", reps, id)
	if err != nil {
		return err
	}

	var noLongerReps []uint64
	for _, current := range listing.Representatives {
		isRep := false
		for _, new := range reps {
			if current == new {
				isRep = true
				break
			}
		}

		if !isRep {
			noLongerReps = append(noLongerReps, current)
		}
	}

	bot.checkRepRole(ctx, noLongerReps...)

	for _, u := range users {
		err = ctx.State.AddRole(common.Conf.GuildID, u.ID, common.Conf.RepRole, api.AddRoleData{
			AuditLogReason: "Add representative role",
		})
		if err != nil {
			return ctx.SendfX("Error adding rep role to %v: %v\nPlease add the role manually.", u.Tag(), err)
		}
	}

	return ctx.SendX("Representatives updated!")
}

func (bot *Bot) checkRepRole(ctx *bcr.Context, users ...uint64) {
	for _, u := range users {
		if bot.DB.IsRepresentative(u) {
			err := ctx.State.AddRole(common.Conf.GuildID, discord.UserID(u), common.Conf.RepRole, api.AddRoleData{
				AuditLogReason: "Add representative role",
			})
			if err != nil {
				common.Log.Errorf("Error adding rep role to %v: %v", u, err)
			}
		} else {
			err := ctx.State.RemoveRole(common.Conf.GuildID, discord.UserID(u), common.Conf.RepRole, "Remove representative role")
			if err != nil {
				common.Log.Errorf("Error removing rep role from %v: %v", u, err)
			}
		}
	}
}
