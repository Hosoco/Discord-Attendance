package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"attendance/src/attendance"
	"attendance/src/handler"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

var s *discordgo.Session
var cmdIDs map[string]string
var appCmds = []discordgo.ApplicationCommand{
	{Name: "clock-in", Type: discordgo.UserApplicationCommand},
	{Name: "clock-out", Type: discordgo.UserApplicationCommand},
	{Name: "new-period", Type: discordgo.UserApplicationCommand},
	{
		Name:        "addhours",
		Description: "Adds hours on top of the current period",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "hours",
				Description:  "Amount of hours to add",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: true,
			},
			{
				Name:         "user",
				Description:  "What user to add hours to",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     false,
				Autocomplete: false,
			},
		},
	},
	{
		Name:        "removehours",
		Description: "Removes hours on top of the current period",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "hours",
				Description:  "Amount of hours to remove",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: true,
			},
			{
				Name:         "user",
				Description:  "What user to remove hours to",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     false,
				Autocomplete: false,
			},
		},
	},
	{
		Name:        "export",
		Description: "Exports the current period to a CSV file",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "user",
				Description:  "What user to export",
				Type:         discordgo.ApplicationCommandOptionUser,
				Required:     false,
				Autocomplete: false,
			},
		},
	},
	{
		Name:        "clockin",
		Description: "Clocks-in the user",
		Type:        discordgo.ChatApplicationCommand,
		Options:     []*discordgo.ApplicationCommandOption{},
	},
	{
		Name:        "clockout",
		Description: "Clocks-out the user",
		Type:        discordgo.ChatApplicationCommand,
		Options:     []*discordgo.ApplicationCommandOption{},
	},
}

var commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"clock-in":    handler.ClockIn,
	"clock-out":   handler.ClockOut,
	"new-period":  handler.NewPeriod,
	"clockin":     handler.ClockIn,
	"clockout":    handler.ClockOut,
	"addhours":    handler.AddHours,
	"removehours": handler.RemoveHours,
	"export":      handler.Export,
}

func assignCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}

func appCommands() {
	cmdIDs = make(map[string]string, len(appCmds))
	for _, cmd := range appCmds {
		rcmd, err := s.ApplicationCommandCreate(os.Getenv("APP_ID"), "", &cmd)
		if err != nil {
			log.Fatalf("Cannot create command %q: %v", cmd.Name, err)
		}
		cmdIDs[rcmd.ID] = rcmd.Name
	}
}

func init() {
	err := godotenv.Load("tokens.env")
	if err != nil {
		fmt.Println("Error loading .env file, will use environment variables instead")
	}
	attendance.Load()
	c := cron.New()
	c.AddFunc("@hourly", func() {
		attendance.Save()
	})
	c.Start()
	s, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	s.AddHandler(handler.DiscordReady)
	s.AddHandler(assignCommands)
	appCommands()
}

func main() {
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()
	shutdown()
}

func shutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
	attendance.Save()

	for id, name := range cmdIDs {
		err := s.ApplicationCommandDelete(os.Getenv("APP_ID"), "", id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}
}
