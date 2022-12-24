package main

import (
	"GOdiscordBOT/app/Commands"
	"GOdiscordBOT/app/Guild"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var token string

func main() {
	loc, err := time.LoadLocation("America/Antigua")
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Local = loc
	// Loading environment variables
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		os.Getenv("HOST"), os.Getenv("USER"), os.Getenv("NAME"), os.Getenv("PASSWORD"), os.Getenv("DBPORT"))

	// Opening connection to database
	conn, err := pgxpool.Connect(context.Background(), dbURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Development token
	token = os.Getenv("DEVELOPMENTTOKEN")
	// Deploy token (UWS)
	// token := os.Getenv("UWSTOKEN")

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
	guildRepo := Guild.NewRepository(conn)
	commandsRepo := Commands.NewRepository(conn)

	// Services
	guildService := Guild.NewService(guildRepo)
	commandsService := Commands.NewService(commandsRepo)

	// Handlers
	Guild.Register(bot, guildService)
	Commands.Register(bot, commandsService)

	// Open connection to discord and start listening
	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close Discord session
	err = bot.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
