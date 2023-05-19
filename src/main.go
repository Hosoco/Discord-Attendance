package main

import (
	"fmt"

	"attendance/src/attendance"
	"attendance/src/bot"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func loadTokens() {
	err := godotenv.Load("tokens.env")
	if err != nil {
		fmt.Println("Error loading .env file, will use environment variables instead")
	}
}

func cronJob() {
	c := cron.New()
	c.AddFunc("@hourly", func() {
		attendance.Save()
	})
	c.Start()
}

func init() {

	loadTokens()
	attendance.Load()
	cronJob()
	bot.InitSession()
}

func main() {
	bot.OpenSession()
	bot.Shutdown()
}
