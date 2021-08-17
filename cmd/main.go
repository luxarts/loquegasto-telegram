package main

import (
	"github.com/joho/godotenv"
	"log"
	"loquegasto-telegram/internal/bot"
)

func main() {
	_ = godotenv.Load()

	b := bot.New()

	log.Println("Bot started")
	b.Start()
}