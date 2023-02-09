package main

import (
	"DiscordBotGO/app/guild"
	"DiscordBotGO/app/serverstatus"
	"DiscordBotGO/app/settings"
	"DiscordBotGO/app/strategicpoint"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	loc, err := time.LoadLocation("America/Antigua")
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Local = loc

	// Loading environment variables
	err = godotenv.Load()
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error loading .env file")
	}

	// Opening connection to database
	conn, err := sql.Open("sqlite3", "./UWSbot.sqlite3")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	_, err = conn.ExecContext(context.Background(), `PRAGMA foreign_keys = ON`)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Development token
	token := os.Getenv("DEVELOPMENTTOKEN")
	// Deploy token (UWS)
	//token := os.Getenv("UWSTOKEN")

	// Create a new Discord session using the provided token
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("invalid token")
	}

	// Setting intents for bot
	bot.Identify.Intents =
		discordgo.IntentsGuildMessages |
			discordgo.IntentsDirectMessages |
			discordgo.IntentsMessageContent |
			discordgo.IntentsGuildMembers |
			discordgo.IntentsGuildMessageReactions |
			discordgo.IntentGuilds |
			discordgo.IntentGuildEmojis

	// Repositories
	guildRepo := guild.NewRepository(conn)
	serverstatusRepo := serverstatus.NewRepository(conn)
	strategicpointRepo := strategicpoint.NewRepository(conn)
	settingsRepo := settings.NewRepository(conn)

	// Services
	guildService := guild.NewService(guildRepo)
	serverstatusService := serverstatus.NewService(serverstatusRepo)
	strategicpointService := strategicpoint.NewService(strategicpointRepo)
	settingsService := settings.NewService(settingsRepo)

	// Handlers
	guild.Register(bot, guildService)
	serverstatus.Register(bot, serverstatusService)
	strategicpoint.Register(bot, strategicpointService)
	settings.Register(bot, settingsService)

	// Open connection to discord and start listening
	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Close Discord session
	err = bot.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
