package guild

// Guilds table name 'guilds' store every guild, and they're features status
type Guilds struct {
	GuildID   string
	GuildName string
}

// Commands table name 'slashcommands' store every command for each guild
type Commands struct {
	CommandID          int
	CommandName        string
	CommandDescription string
	DefaultCommand     bool
}

// Guildcommands tabel name 'guildcommands' store all commands that a guild can have
type Guildcommands struct {
	GuildID   string
	Guild     Guilds
	CommandID int
	Command   Commands
}

// Guildemojis tabel name 'guildemojis' store every custom emoji of a guild
type Guildemojis struct {
	GuildID   string
	EmojiId   string
	EmojiName string
	Animated  bool
}
