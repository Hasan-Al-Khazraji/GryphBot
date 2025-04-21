package main

import (
	"os"

	bot "github.com/Hasan-Al-Khazraji/GryphBot/Bot"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	bot.BotToken = os.Getenv("DISCORD_TOKEN")
	bot.Run()
}
