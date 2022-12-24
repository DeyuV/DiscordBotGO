package Guild

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func NewRepository(db *pgxpool.Pool) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db *pgxpool.Pool
}

func (r *repositoryImpl) DeleteEmojis(ctx context.Context, guildId string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM guildemojis WHERE guildid = $1`, guildId)
	return err
}

func (r *repositoryImpl) AddEmojis(ctx context.Context, guildId, emojiId, emojiName string, animated bool) error {
	_, err := r.db.Exec(ctx, `INSERT INTO guildemojis (guildid, emojiid, emojiname, animated) VALUES ($1,$2,$3,$4)`, guildId, emojiId, emojiName, animated)

	return err
}

func (r *repositoryImpl) AddDefaultCommands(ctx context.Context, guildId string) error {
	var commandsID []int

	rows, err := r.db.Query(ctx, `SELECT commandid FROM commands WHERE defaultcommand = true`)

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
		_, err = r.db.Exec(ctx, `INSERT INTO guildcommands(guildid, commandid) VALUES ($1, $2)`, guildId, commandID)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *repositoryImpl) GetSlashCommands(ctx context.Context, guildId string) ([]Guildcommands, error) {
	var commands []Guildcommands

	rows, err := r.db.Query(ctx, `SELECT * FROM guildcommands WHERE guildid = $1`, guildId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var command Guildcommands
		err := rows.Scan(&command.GuildID, &command.CommandID)
		if err != nil {
			log.Fatal(err)
		}
		err = r.db.QueryRow(ctx, `SELECT * FROM commands WHERE commandid = $1`, command.CommandID).Scan(&command.Command.CommandID, &command.Command.CommandName, &command.Command.CommandDescription, &command.Command.DefaultCommand)
		if err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}

	return commands, nil
}

func (r *repositoryImpl) Add(ctx context.Context, guildId string, guildName string) error {

	_, err := r.db.Exec(ctx, `INSERT INTO guilds (id, name) VALUES ($1, $2)`, guildId, guildName)

	return err
}

func (r *repositoryImpl) GetById(ctx context.Context, guildId string) (*Guilds, error) {
	var guild Guilds

	err := r.db.QueryRow(ctx, `SELECT * FROM guilds WHERE id = $1`, guildId).Scan(&guild)
	if err != nil {
		return nil, err
	}
	return &guild, nil
}

// DeleteGuild deletes a guild from database
func (r *repositoryImpl) DeleteGuild(ctx context.Context, guildId string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM guilds WHERE id = $1`, guildId)

	return err
}
