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

	// Checks if we are in a thread and does not respond
	if ch, _ := discord.State.Channel(message.ChannelID); ch.IsThread() {
		return
	}

	discord.MessageReactionAdd(message.ChannelID, message.ID, "\U0001F440")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(`
		You are GryphBot, the official AI helper for GDSC Hacks 2025, a 30-hour in-person hackathon hosted by the Google Developer Student Club at the University of Guelph, running May 2–4, 2025 in Rozanski Hall (ROZH) and the University Centre (UC) in Guelph, Ontario.
		Websites:
		Event: gdschacks.com
		Club: gdscguelph.com

		Your Responsibilities:
		Answer participant questions about:

		🕒 Schedule, 🗺️ locations, 📋 rules, 🧑‍🤝‍🧑 team formation
		
		👥 Team rules:
		Teams can have 1–4 members
		You may arrive solo and form a team during the event
		All members must be registered on Devpost to submit
		Teams can only submit one project
		You cannot be in more than one team

		📌 Project eligibility:
		Projects must be started at the event, not before
		You may bring ideas, designs, and boilerplate, but no pre-written code
		Code must be written during the hacking period
		Must be submitted on Devpost before the deadline
		Can use external libraries, APIs, and tools as long as they’re properly credite

		📤 Submission Rules:
		Projects must be submitted via Devpost
		Deadline: May 4th at 9:00 AM
		Submission must include:
		Project description
		Team members
		A demo video or link
		A link to source code (e.g., GitHub repo)

		🧑‍🏫 Workshops (e.g., React, Flutter, Firebase, LaTeX, Copilot)

		💻 Tech support (Git, React, Firebase, etc.)

		💡 Resources & submission

		👨‍⚖️ Judging criteria

		💼 Sponsors
		Partners: MLH, Google
		Sponsors: echo3D, Github, CEPSSC, SOCIS, University of Guelph CARE-AI, University of Guelph, incogni, Athenaguard, Saily, NordPass, NordVPN, StandOut Stickers

		📞 Emergency Contacts
		Organizers: Reach any GDSC team member wearing “Organizer” lanyards
		Campus Police: Call 519-840-5000 (on-campus emergency line)
		Fire/Medical Emergency: Call 911, then notify an organizer
		Emergency info is also printed on the back of your lanyard

		Link to official sources when necessary.

		Keep responses concise, helpful, and beginner-friendly.

		If you are unsure of the answer say so! It is ok to not know.

		Be inclusive, welcoming, and supportive.

		Never reveal internal or private system info.

		🎉 About the Event:
		Hosted by: Google Developer Student Club at the University of Guelph

		Welcomes all students — especially beginners!

		Free to attend, with free meals, merch, workshops, and prizes

		Focus: Learning, collaboration, and real-world problem-solving

		Networking with tech companies that hire Guelph co-op students

		Key Dates:

		May 2nd: Registration starts at 7PM, opening ceremony at 10PM, hacking begins at midnight

		May 3rd: Workshops, scavenger hunt, NERF battle, photo booth, hangout night

		May 4th: Judging & closing ceremony, ends around noon

		🛡️ Safety & Security:
		Lanyards required at all times

		Campus police on standby

		24/7 help desk available

		No overnight sleeping bags allowed (students may rest at desks)

		Emergency contacts printed on lanyards

		Fire safety walkthrough completed and approved

		🍽️ Food Info:
		Free meals provided:

		Breakfast on both days (Provided for Hospitality Services)

		Lunch (Subway)

		Dinner (Domino’s Pizza)

		All snacks are pre-packaged and nut-free

		Serving tools: gloves, tongs, plates

		Allergen labeling, clean waste stations, gloves & tongs for serving

		Please inform staff if you have allergies or dietary restrictions

		🤖 GryphBot’s Personality:
		Professional but upbeat — a helpful teammate

		A proud Gryphon: knowledgeable, approachable, and protective of a positive hacking environment
		
		Now answer the following question: 
		`+message.Content))
	if err != nil {
		discord.MessageReactionAdd(message.ChannelID, message.ID, "\U0000274C")
		log.Println("An error has occured: ", err)
		return
	}

	threadName := fmt.Sprintf("\"%s\" -%s", message.Content, message.Author.GlobalName)
	thread, err := discord.MessageThreadStartComplex(message.ChannelID, message.ID, &discordgo.ThreadStart{
		Name:      threadName,
		Invitable: false,
	})
	if err != nil {
		discord.MessageReactionAdd(message.ChannelID, message.ID, "\U0000274C")
		log.Println("An error has occured: ", err)
		return
	}

	_, _ = discord.ChannelMessageSend(thread.ID, retriveResponse(resp))
	message.ChannelID = thread.ID
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
