package db

import (
	"context"
	"fmt"
	"strings"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
)

type Listing struct {
	ID              int64
	Name            string
	Description     string
	URL             string
	EmbedURL        bool
	Representatives []uint64
}

// Message returns the listing formatted for a message.
func (l Listing) Message() string {
	plural := ""
	if len(l.Representatives) != 1 {
		plural = "s"
	}

	s := fmt.Sprintf("**%v**\n\n%v\n\nContact%v: ", l.Name, l.Description, plural)
	var reps []string
	for _, rep := range l.Representatives {
		reps = append(reps, discord.RoleID(rep).Mention())
	}

	s += strings.Join(reps, " | ")

	if !l.EmbedURL {
		s += "\n\n<" + l.URL + ">"
	} else {
		s += "\n\n" + l.URL
	}

	return s
}

type ListingMessage struct {
	ID        discord.MessageID
	ChannelID discord.ChannelID
	ListingID int64
}

func (db *DB) Listing(id int64) (l Listing, err error) {
	err = pgxscan.Get(context.Background(), db, &l, "select * from listings where id = $1", id)
	return l, errors.Cause(err)
}
