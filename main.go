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

var roleChannelID = "918866377899655238"
var roleMessageID = "918876972468289566"

var questionChannelID = "928542953310392330"

func initReaction(bot *discordgo.Session) {
	for emoji := range roles {
		err := bot.MessageReactionAdd(roleChannelID, roleMessageID, emoji)
		if err != nil {
			fmt.Println("error init reaction,", err)
			return
		}
	}

}

func addReaction(bot *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if (reaction.GuildID == guildID) && (reaction.ChannelID == roleChannelID) && (reaction.MessageID == roleMessageID) {
		role := roles[reaction.Emoji.Name]
		err := bot.GuildMemberRoleAdd(guildID, reaction.UserID, role)
		if err != nil {
			fmt.Println("error add role,", err)
			return
		}
		sendDM(bot, reaction.UserID, "ロールを付与しました✅")
	}

	if (reaction.GuildID == guildID) && (reaction.ChannelID == questionChannelID) {
		message, err := bot.ChannelMessage(reaction.ChannelID, reaction.MessageID)
		if err != nil {
			fmt.Println("error check all reaction,", err)
			return
		}
		for _, messageReaction := range message.Reactions {
			if messageReaction.Emoji.Name == reaction.Emoji.Name {
				break
			}

			users, err := bot.MessageReactions(reaction.ChannelID, reaction.MessageID, messageReaction.Emoji.Name, 100, "", "")
			if err != nil {
				fmt.Println("error check reaction,", err)
				return
			}
			for _, user := range users {
				if user.ID == reaction.UserID {
					err := bot.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, messageReaction.Emoji.Name, user.ID)
					if err != nil {
						fmt.Println("error remove reaction,", err)
						return
					}
				}
			}
		}
	}
}

func removeReaction(bot *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	if (reaction.GuildID == guildID) && (reaction.ChannelID == roleChannelID) && (reaction.MessageID == roleMessageID) {
		role := roles[reaction.Emoji.Name]
		err := bot.GuildMemberRoleRemove(guildID, reaction.UserID, role)
		if err != nil {
			fmt.Println("error remove role,", err)
			return
		}
		sendDM(bot, reaction.UserID, "ロールを削除しました✅")
	}
}

func sendDM(bot *discordgo.Session, userID string, message string) {
	channel, err := bot.UserChannelCreate(userID)
	if err != nil {
		return
	}
	_, err = bot.ChannelMessageSend(channel.ID, message)
	if err != nil {
		fmt.Println("error sending DM,", err)
		return
	}
}
