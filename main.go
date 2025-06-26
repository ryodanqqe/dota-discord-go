package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Println("DISCORD_TOKEN environment variable not set")
		return
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	session.AddHandler(messageCreate)
	session.AddHandler(interactionCreate)

	if err = session.Open(); err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	registerSlashCommands(session)

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	_ = session.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping!" {
		_, err := s.ChannelMessageSend(m.ChannelID, "pong!")
		if err != nil {
			fmt.Println("error sending message,", err)
		}
	}

}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	switch i.ApplicationCommandData().Name {
	case "ping":
		handlePingCommand(s, i)
	default:
		fmt.Println("I don't know that command")

	}
}

func handlePingCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong",
		},
	})

	if err != nil {
		fmt.Println(err)
	}
}

func registerSlashCommands(s *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Responds with pong",
		},
	}

	for _, command := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
		if err != nil {
			fmt.Println(err)
		}
	}
}
