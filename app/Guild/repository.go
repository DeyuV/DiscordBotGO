package guild

import (
	"context"
	"database/sql"
	"log"
)

func NewRepository(db *sql.DB) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db *sql.DB
}

func (r *repositoryImpl) DeleteEmojis(ctx context.Context, guildId string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM guildemojis WHERE guildid = ?`, guildId)
	return err
}

func (r *repositoryImpl) AddEmojis(ctx context.Context, guildId, emojiId, emojiName string, animated bool) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO guildemojis (guildid, emojiid, emojiname, animated) VALUES (?,?,?,?)`, guildId, emojiId, emojiName, animated)

	return err
}

func (r *repositoryImpl) AddDefaultCommands(ctx context.Context, guildId string) error {
	var commandsID []int

	rows, err := r.db.QueryContext(ctx, `SELECT commandid FROM commands WHERE defaultcommand = true`)

	if err != nil {
		return err
	}

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		commandsID = append(commandsID, id)
	}

	for _, commandID := range commandsID {
		_, err = r.db.ExecContext(ctx, `INSERT INTO guildcommands(guildid, commandid) VALUES (?, ?)`, guildId, commandID)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *repositoryImpl) GetSlashCommands(ctx context.Context, guildId string) ([]Guildcommands, error) {
	var commands []Guildcommands

	rows, err := r.db.QueryContext(ctx, `SELECT * FROM guildcommands WHERE guildid = ?`, guildId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var command Guildcommands
		err := rows.Scan(&command.GuildID, &command.CommandID)
		if err != nil {
			log.Fatal(err)
		}
		err = r.db.QueryRowContext(ctx, `SELECT * FROM commands WHERE commandid = ?`, command.CommandID).Scan(&command.Command.CommandID, &command.Command.CommandName, &command.Command.CommandDescription, &command.Command.DefaultCommand)
		if err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}

	return commands, nil
}

func (r *repositoryImpl) Add(ctx context.Context, guildId string, guildName string) error {

	_, err := r.db.ExecContext(ctx, `INSERT INTO guilds (id, name) VALUES (?, ?)`, guildId, guildName)

	return err
}

func (r *repositoryImpl) GetById(ctx context.Context, guildId string) (*Guilds, error) {
	var guild Guilds

	err := r.db.QueryRowContext(ctx, `SELECT * FROM guilds WHERE id = ?`, guildId).Scan(&guild.GuildID, &guild.GuildName)
	if err != nil {
		return nil, err
	}
	return &guild, nil
}

func (r *repositoryImpl) DeleteGuild(ctx context.Context, guildId string) error {
	r.db.ExecContext(context.Background(), `PRAGMA foreign_keys = ON`)
	_, err := r.db.ExecContext(ctx, `DELETE FROM guilds WHERE id = ?`, guildId)

	return err
}
