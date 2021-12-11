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

var geriraID = "💛"
var seichiID = "💚"

func initReaction(bot *discordgo.Session) {
	err := bot.MessageReactionAdd(channelID, messageID, geriraID)
	err = bot.MessageReactionAdd(channelID, messageID, seichiID)
	if err != nil {
		fmt.Println("error init reaction,", err)
		return
	}
}

func addReaction(bot *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if (reaction.GuildID != guildID) && (reaction.MessageID != messageID) {
		return
	}

	switch reaction.Emoji.Name {
	case geriraID:
		err := bot.GuildMemberRoleAdd(guildID, reaction.UserID, "918889974399639614")
		if err != nil {
			fmt.Println("error add role,", err)
			return
		}
	case seichiID:
		err := bot.GuildMemberRoleAdd(guildID, reaction.UserID, "889077232209895454")
		if err != nil {
			fmt.Println("error add role,", err)
			return
		}
	}
}

func removeReaction(bot *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	if reaction.GuildID != guildID {
		return
	}
	if reaction.MessageID != messageID {
		return
	}

	switch reaction.Emoji.Name {
	case geriraID:
		err := bot.GuildMemberRoleRemove(guildID, reaction.UserID, "918889974399639614")
		if err != nil {
			fmt.Println("error remove role,", err)
			return
		}
	case seichiID:
		err := bot.GuildMemberRoleRemove(guildID, reaction.UserID, "889077232209895454")
		if err != nil {
			fmt.Println("error remove role,", err)
			return
		}
	}
}
