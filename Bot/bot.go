package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotToken string

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func Run() {

	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	discord.AddHandler(newMessage)

	discord.Open()
	defer discord.Close()

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	// Prevent bot from sending message to itself
	if message.Author.ID == discord.State.User.ID {
		return
	}

	switch {
	case strings.Contains(message.Content, "help"):
		discord.ChannelMessageSend(message.ChannelID, "Hello World😃")
	case strings.Contains(message.Content, "bye"):
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
	default:
		discord.MessageReactionAdd(message.ChannelID, message.ID, "\U0001F440")
		thread, err := discord.MessageThreadStartComplex(message.ChannelID, message.ID, &discordgo.ThreadStart{
			Name:      message.Content + " -" + message.Author.Username,
			Invitable: false,
		})
		fmt.Println(message.Content)
		if err != nil {
			panic(err)
		}
		_, _ = discord.ChannelMessageSend(thread.ID, "pong")
		message.ChannelID = thread.ID
	}

}
