package db

import (
	"context"
	"fmt"
	"strings"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/natures-networker/common"
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
		reps = append(reps, discord.UserID(rep).Mention())
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

func (db *DB) AllListings() (l []Listing, err error) {
	err = pgxscan.Select(context.Background(), db, &l, "select * from listings order by id")
	return l, errors.Cause(err)
}

func (db *DB) ListingMessages(id int64) (m []ListingMessage, err error) {
	err = pgxscan.Select(context.Background(), db, &m, "select * from listing_messages where listing_id = $1 order by id", id)
	return m, errors.Cause(err)
}

func (db *DB) AddListing(name, url string, embed bool) (l Listing, err error) {
	err = pgxscan.Get(context.Background(), db, &l, "insert into listings (name, url, embed_url) values ($1, $2, $3) returning *", name, url, embed)
	return l, errors.Cause(err)
}

func (db *DB) IsRepresentative(id uint64) (isRep bool) {
	err := db.QueryRow(context.Background(), "select exists(select * from listings where $1 = any(representatives))", id).Scan(&isRep)
	if err != nil {
		common.Log.Errorf("Error checking if user is representative: %v", err)
	}
	return isRep
}

func (db *DB) AddListingMessage(listing int64, id discord.MessageID, ch discord.ChannelID) error {
	_, err := db.Exec(context.Background(), "insert into listing_messages (id, channel_id, listing_id) values ($1, $2, $3)", id, ch, listing)
	return err
}
