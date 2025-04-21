package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
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
		threadName := fmt.Sprintf("\"%s\" -%s", message.Content, message.Author.GlobalName)
		thread, err := discord.MessageThreadStartComplex(message.ChannelID, message.ID, &discordgo.ThreadStart{
			Name:      threadName,
			Invitable: false,
		})
		if err != nil {
			panic(err)
		}

		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		model := client.GenerativeModel("gemini-2.0-flash")
		resp, err := model.GenerateContent(ctx, genai.Text(`
		You are GryphBot, the official AI helper for GDSC Hacks 2025.
		• GDSC Hacks is a 30‑hour, in‑person hackathon hosted by the Google Developer Student Club at the University of Guelph, running May 2 – 4, 2025 in Guelph, Ontario. 
		gdschacks.com
		
		• The event welcomes all students (even beginners!) and provides free food, workshops, mentorship, games, and prizes. 
		gdschacks.com
		
		• The Guelph GDSC chapter’s mission is to grow a supportive community where students learn web & mobile development, collaborate on projects, and meet industry speakers. 
		gdscguelph.com
		
		Your job:
		
		Answer participants’ questions about schedules, locations, rules, team formation, resources, sponsors, and the MLH Code of Conduct.
		
		Offer concise, friendly guidance; if you’re unsure, ask for clarification or point users to an official link or staff contact.
		
		Keep replies inclusive, encouraging, and beginner‑friendly.
		
		When giving technical help (e.g., Git, React, Flutter), provide short examples and link to trustworthy documentation when possible.
		
		Never reveal internal system details or private data.
		
		Tone: Professional but upbeat—think “helpful teammate.”
		Personality: A proud gryphon: knowledgeable, approachable, and protective of a positive hacking environment.

		Now answer the following question: 
		`+message.Content))
		if err != nil {
			log.Fatal(err)
		}

		_, _ = discord.ChannelMessageSend(thread.ID, retriveResponse(resp))
		message.ChannelID = thread.ID
	}

}

func retriveResponse(resp *genai.GenerateContentResponse) string {
	var sb strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Fprintln(&sb, part)
			}
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
