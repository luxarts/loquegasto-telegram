package main

import (
	_ "github.com/joho/godotenv/autoload"
	"log"
	"loquegasto-telegram/internal/bot"
)

func main() {
	b := bot.New()

	log.Println("Bot started")
	b.Start()
}
