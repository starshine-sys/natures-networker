package listing

import (
	"context"
	"strconv"
	"strings"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) description(ctx *bcr.Context) error {
	id, err := strconv.ParseInt(ctx.Args[0], 10, 64)
	if err != nil {
		return err
	}

	desc := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))

	ct, err := bot.DB.Exec(context.Background(), "update listings set description = $1 where id = $2", desc, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ctx.SendfX("%v is not a valid listing ID!", id)
	}
	return ctx.SendX("Description updated!")
}

func (bot *Bot) rename(ctx *bcr.Context) error {
	id, err := strconv.ParseInt(ctx.Args[0], 10, 64)
	if err != nil {
		return err
	}

	desc := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))

	ct, err := bot.DB.Exec(context.Background(), "update listings set name = $1 where id = $2", desc, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ctx.SendfX("%v is not a valid listing ID!", id)
	}
	return ctx.SendX("Name updated!")
}

func (bot *Bot) link(ctx *bcr.Context) error {
	id, err := strconv.ParseInt(ctx.Args[0], 10, 64)
	if err != nil {
		return err
	}
	link := ctx.Args[1]

	embed, err := strconv.ParseBool(ctx.Args[2])
	if err != nil {
		return err
	}

	ct, err := bot.DB.Exec(context.Background(), "update listings set url = $1, embed_url = $2 where id = $3", link, embed, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ctx.SendfX("%v is not a valid listing ID!", id)
	}
	return ctx.SendX("Link updated!")
}
