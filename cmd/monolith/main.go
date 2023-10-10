package main

import (
	"DiscordBotGO/app/guild"
	"DiscordBotGO/app/settings"
	"DiscordBotGO/app/strategicpoint"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	//"github.com/joho/godotenv"
)

func main() {
	/* 	err := godotenv.Load()
	   	if err != nil {
	   		fmt.Println("Error loading .env file" + err.Error())
	   		return
	   	} */

	// Set up the connection pool
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		panic(err)
	}
	defer dbpool.Close()

	// Ping the database to ensure a successful connection
	err = dbpool.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to PostgreSQL database successfully!")

	// Development token
	//token := os.Getenv("DEVELOPMENTTOKEN")
	// Deploy token (UWS)
	token := os.Getenv("UWSTOKEN")

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
	guildRepo := guild.NewRepository(dbpool)
	strategicpointRepo := strategicpoint.NewRepository(dbpool)
	settingsRepo := settings.NewRepository(dbpool)

	// Services
	guildService := guild.NewService(guildRepo)
	strategicpointService := strategicpoint.NewService(strategicpointRepo)
	settingsService := settings.NewService(settingsRepo)

	// Handlers
	guild.Register(bot, guildService)
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
