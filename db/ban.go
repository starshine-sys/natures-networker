package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

type Ban struct {
	ID        int64
	MessageID *discord.MessageID
	UserIDs   []uint64
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
