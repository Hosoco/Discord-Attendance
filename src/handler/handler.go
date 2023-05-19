package handler

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"time"

	"attendance/src/attendance"
	"attendance/src/auth"
	"attendance/src/export"

	"github.com/bwmarrin/discordgo"
)

func DiscordReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Println("Bot is up!")
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, response string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func getHours(s *discordgo.Session, i *discordgo.InteractionCreate) (int64, error) {
	return strconv.ParseInt(i.ApplicationCommandData().Options[0].StringValue(), 10, 64)
}

func ClockIn(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if attendance.ClockedIn(i.Member.User.ID) {
		respond(s, i, "You are already clocked in!")
	} else {
		attendance.ClockIn(i.Member.User.ID)
		respond(s, i, "Clocked in at "+time.Now().Format("15:04:05"))
	}
}

func ClockOut(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !attendance.ClockedIn(i.Member.User.ID) {
		respond(s, i, "You are already clocked out!")
	} else {
		attendance.ClockOut(i.Member.User.ID)
		respond(s, i, "Clocked out at "+time.Now().Format("15:04:05"))
	}
}

func ChangeHours(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if !auth.IsAuthenticated(i.Member.User.ID) {
		return
	}

	hours, err := getHours(s, i)
	if err != nil {
		respond(s, i, "Invalid number of hours")
		return
	}
	if i.ApplicationCommandData().Name == "removehours" {
		hours = -hours
	}
	switch len(i.ApplicationCommandData().Options) {
	case 1:
		ownUser := i.Member.User.ID
		attendance.ChangeHours(ownUser, hours)
	case 2:
		targetUser := i.ApplicationCommandData().Options[1].UserValue(s).ID
		attendance.ChangeHours(targetUser, hours)
	default:
		respond(s, i, "Invalid number of arguments")
		return
	}

	respond(s, i, "Successfully changed hours of user!")
}

func Export(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !auth.IsAuthenticated(i.Member.User.ID) {
		return
	}
	var fp string
	if len(i.ApplicationCommandData().Options) == 1 {
		fp = export.CSV(i.ApplicationCommandData().Options[0].UserValue(s).ID)
	} else {
		fp = export.CSV(i.Member.User.ID)
	}
	data, _ := os.ReadFile(fp)
	respond(s, i, "Exported data!")
	s.ChannelFileSend(i.ChannelID, "export.csv", bytes.NewReader(data))
	os.Remove(fp)
}

func NewPeriod(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !auth.IsAuthenticated(i.Member.User.ID) {
		return
	}
	attendance.NewPeriod()
	respond(s, i, "Data reset, now on new period!")
}
