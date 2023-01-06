package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"loquegasto-telegram/internal/bot"
	"net/http"
)

func main() {
	b := bot.New()

	log.Println("Bot started")
	go b.Start()

	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})

	if err := router.Run(); err != nil {
		log.Fatalln(err)
	}
}
