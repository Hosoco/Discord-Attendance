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

func AddHours(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//If interaction is a slash command
	if i.Type == discordgo.InteractionApplicationCommand {
		if !auth.IsAuthenticated(i.Member.User.ID) {
			return
		}
		hours, err := strconv.ParseInt(i.ApplicationCommandData().Options[0].StringValue(), 10, 64)
		if err != nil {
			respond(s, i, "Invalid number of hours")
			return
		}
		if len(i.ApplicationCommandData().Options) == 1 {
			attendance.ChangeHours(i.Member.User.ID, hours)
		} else if len(i.ApplicationCommandData().Options) == 2 {
			attendance.ChangeHours(i.ApplicationCommandData().Options[1].UserValue(s).ID, hours)
		} else {
			respond(s, i, "Invalid number of arguments")
			return
		}
		respond(s, i, "Successfully added hours to user!")
	}
}

func RemoveHours(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !auth.IsAuthenticated(i.Member.User.ID) {
		return
	}
	hours, err := strconv.ParseInt(i.ApplicationCommandData().Options[0].StringValue(), 10, 64)
	if err != nil {
		respond(s, i, "Invalid number of hours")
		return
	}
	hours = -hours
	if len(i.ApplicationCommandData().Options) == 1 {
		attendance.ChangeHours(i.Member.User.ID, hours)
	} else if len(i.ApplicationCommandData().Options) == 2 {
		attendance.ChangeHours(i.ApplicationCommandData().Options[1].UserValue(s).ID, hours)
	} else {
		respond(s, i, "Invalid number of arguments")
		return
	}
	respond(s, i, "Successfully removed hours of user!")
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
