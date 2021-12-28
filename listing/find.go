package listing

import (
	"context"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/mozillazg/go-unidecode"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/natures-networker/db"
)

func (bot *Bot) findByRep(ctx *bcr.Context) (err error) {
	u, err := ctx.ParseUser(ctx.RawArgs)
	if err != nil {
		return ctx.SendX("User not found.")
	}

	var listings []db.Listing
	err = pgxscan.Select(context.Background(), bot.DB, &listings, "select * from listings where $1 = any(representatives) order by id", u.ID)
	if err != nil {
		return err
	}

	if len(listings) == 0 {
		return ctx.SendfX("%v is not a representative.", u.Tag())
	}

	s := fmt.Sprintf("%v/%v is representative for the following listing(s):\n", u.Tag(), u.ID)
	for _, listing := range listings {
		s += fmt.Sprintf("`%v`. %v (%v)\n", listing.ID, listing.Name, unidecode.Unidecode(listing.Name))
	}

	return ctx.SendX(s)
}

func (bot *Bot) findByName(ctx *bcr.Context) (err error) {
	listings, err := bot.DB.AllListings()
	if err != nil {
		return err
	}

	var found []db.Listing
	for _, listing := range listings {
		if strings.Contains(strings.ToLower(unidecode.Unidecode(listing.Name)), strings.ToLower(unidecode.Unidecode(ctx.RawArgs))) {
			found = append(found, listing)
		}
	}

	if len(found) == 0 {
		return ctx.SendfX("No listing's name contains ``%v``", bcr.EscapeBackticks(ctx.RawArgs))
	}

	s := ""
	for _, listing := range found {
		s += fmt.Sprintf("`%v`. %v (%v)\n", listing.ID, listing.Name, unidecode.Unidecode(listing.Name))
	}

	return ctx.SendX(s)
}
