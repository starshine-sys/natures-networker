package ban

import (
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/natures-networker/common"
)

func (bot *Bot) evidence(ctx *bcr.Context) (err error) {
	id, err := strconv.ParseInt(ctx.Args[0], 10, 64)
	if err != nil {
		return err
	}

	evidence := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))

	log, err := bot.DB.UpdateEvidence(id, evidence)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ctx.SendX("That doesn't seem to be a ban log ID.")
		}
		return err
	}

	if log.MessageID == nil {
		msg, err := ctx.State.SendEmbeds(common.Conf.BanLog, log.Embed(ctx))
		if err != nil {
			return err
		}

		err = bot.DB.SetBanMessage(log.ID, msg.ID)
		if err != nil {
			return err
		}
	} else {
		_, err = ctx.State.EditEmbeds(common.Conf.BanLog, *log.MessageID, log.Embed(ctx))
		if err != nil {
			return err
		}
	}

	return ctx.SendX("Evidence updated!")
}

func (bot *Bot) reason(ctx *bcr.Context) (err error) {
	id, err := strconv.ParseInt(ctx.Args[0], 10, 64)
	if err != nil {
		return err
	}

	reason := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))

	log, err := bot.DB.UpdateReason(id, reason)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ctx.SendX("That doesn't seem to be a ban log ID.")
		}
		return err
	}

	if log.MessageID == nil {
		msg, err := ctx.State.SendEmbeds(common.Conf.BanLog, log.Embed(ctx))
		if err != nil {
			return err
		}

		err = bot.DB.SetBanMessage(log.ID, msg.ID)
		if err != nil {
			return err
		}
	} else {
		_, err = ctx.State.EditEmbeds(common.Conf.BanLog, *log.MessageID, log.Embed(ctx))
		if err != nil {
			return err
		}
	}

	return ctx.SendX("Reason updated!")
}
