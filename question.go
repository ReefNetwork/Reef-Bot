package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var questionChannelID = "928542953310392330"

func onAddReactionQuestionHandler(bot *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if (reaction.GuildID == guildID) && (reaction.ChannelID == questionChannelID) {
		go duplicateRemove(bot, reaction.UserID, reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name)
	}
}

func duplicateRemove(bot *discordgo.Session, userID string, channelID string, messageID string, emojiName string) {
	message, err := bot.ChannelMessage(channelID, messageID)
	if err != nil {
		fmt.Println("error check all reaction,", err)
		return
	}
	for _, messageReaction := range message.Reactions {
		if messageReaction.Emoji.Name != emojiName {
			go removeUserReaction(bot, userID, channelID, messageID, messageReaction.Emoji.Name)
		}
	}
}

func removeUserReaction(bot *discordgo.Session, userID string, channelID string, messageID string, emojiName string) {
	users, err := bot.MessageReactions(channelID, messageID, emojiName, 100, "", "")
	if err != nil {
		fmt.Println("error check reaction,", err)
		return
	}
	for _, user := range users {
		if user.ID == userID {
			err := bot.MessageReactionRemove(channelID, messageID, emojiName, user.ID)
			if err != nil {
				fmt.Println("error remove reaction,", err)
				return
			}
		}
	}
}
