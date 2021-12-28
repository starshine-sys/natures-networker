package common

import "github.com/diamondburned/arikawa/v3/discord"

type Config struct {
	Token       string `toml:"token"`
	DatabaseURL string `toml:"database_url"`

	Prefix  string          `toml:"prefix"`
	Owner   discord.UserID  `toml:"owner"`
	GuildID discord.GuildID `toml:"guild_id"`

	StaffRole  discord.RoleID `toml:"staff_role"`
	HelperRole discord.RoleID `toml:"helper_role"`

	ListingLog discord.ChannelID `toml:"listing_log"`
	BanLog     discord.ChannelID `toml:"ban_log"`
}

var Conf Config
