package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

var (
	token string
	roles = map[string]string{
		"💛": "918889974399639614",
		"💚": "889077232209895454"}
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	bot.AddHandler(addReaction)
	bot.AddHandler(removeReaction)

	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.")
	initReaction(bot)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	//goland:noinspection GoUnhandledErrorResult
	bot.Close()
}

var guildID = "638760361369010177"
var channelID = "918866377899655238"
var messageID = "918876972468289566"

func initReaction(bot *discordgo.Session) {
	for emoji := range roles {
		err := bot.MessageReactionAdd(channelID, messageID, emoji)
		if err != nil {
			fmt.Println("error init reaction,", err)
			return
		}
	}

}

func addReaction(bot *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if (reaction.GuildID != guildID) && (reaction.ChannelID != channelID) && (reaction.MessageID != messageID) {
		return
	}

	role := roles[reaction.Emoji.Name]
	err := bot.GuildMemberRoleAdd(guildID, reaction.UserID, role)
	if err != nil {
		fmt.Println("error add role,", err)
		return
	}
}

func removeReaction(bot *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	if (reaction.GuildID != guildID) && (reaction.ChannelID != channelID) && (reaction.MessageID != messageID) {
		return
	}

	role := roles[reaction.Emoji.Name]
	err := bot.GuildMemberRoleRemove(guildID, reaction.UserID, role)
	if err != nil {
		fmt.Println("error remove role,", err)
		return
	}
}
