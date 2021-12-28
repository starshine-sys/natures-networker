package db

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

type Ban struct {
	ID        int64
	MessageID *discord.MessageID
	UserIDs   []uint64 `db:"user_ids"`
	Reason    string
	Evidence  *string
	CreatedAt time.Time
}

func (b Ban) Embed(ctx *bcr.Context) discord.Embed {
	plural := ""
	if len(b.UserIDs) != 1 {
		plural = "s"
	}

	e := discord.Embed{
		Title:       fmt.Sprintf("User%v banned", plural),
		Color:       bcr.ColourRed,
		Timestamp:   discord.NewTimestamp(b.CreatedAt),
		Description: fmt.Sprintf("**User ID%v**\n", plural),
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v", b.ID),
		},
	}

	usernames := make([]string, 0, len(b.UserIDs))
	for _, id := range b.UserIDs {
		e.Description += strconv.FormatUint(id, 10) + "\n"

		user, err := ctx.State.User(discord.UserID(id))
		if err != nil {
			usernames = append(usernames, fmt.Sprintf("*deleted user %v*", id))
		} else {
			usernames = append(usernames, user.Tag())
		}
	}

	if len(usernames) > 0 {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Last known username" + plural,
			Value: strings.Join(usernames, "\n"),
		})
	}

	e.Fields = append(e.Fields, discord.EmbedField{
		Name:  "Reason",
		Value: b.Reason,
	})

	if b.Evidence != nil {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Evidence",
			Value: *b.Evidence,
		})
	}

	return e
}

func (db *DB) CreateBan(users []*discord.User, reason string) (b Ban, err error) {
	ids := make([]uint64, 0, len(users))
	for _, u := range users {
		ids = append(ids, uint64(u.ID))
	}

	err = pgxscan.Get(context.Background(), db, &b, "insert into bans (user_ids, reason) values ($1, $2) returning *", ids, reason)
	return b, errors.Cause(err)
}

func (db *DB) SetBanMessage(id int64, msg discord.MessageID) error {
	_, err := db.Exec(context.Background(), "update bans set message_id = $1 where id = $2", msg, id)
	return err
}

func (db *DB) UpdateEvidence(id int64, evidence string) (b Ban, err error) {
	err = pgxscan.Get(context.Background(), db, &b, "update bans set evidence = $1 where id = $2 returning *", evidence, id)
	return b, errors.Cause(err)
}

func (db *DB) UpdateReason(id int64, reason string) (b Ban, err error) {
	err = pgxscan.Get(context.Background(), db, &b, "update bans set reason = $1 where id = $2 returning *", reason, id)
	return b, errors.Cause(err)
}
