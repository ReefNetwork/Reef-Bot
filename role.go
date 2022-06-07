package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var roleChannelID = "918866377899655238"
var roleMessageID = "918876972468289566"
var joinRoleMessageID = "983017946827866132"

func initReaction(bot *discordgo.Session) {
	for emoji := range roles {
		err := bot.MessageReactionAdd(roleChannelID, roleMessageID, emoji)
		if err != nil {
			fmt.Println("error init reaction,", err)
			return
		}
	}
	for emoji := range joinRoles {
		err := bot.MessageReactionAdd(roleChannelID, joinRoleMessageID, emoji)
		if err != nil {
			fmt.Println("error init join reaction,", err)
			return
		}
	}
}

func onJoinRoleHandler(bot *discordgo.Session, join *discordgo.GuildMemberAdd) {
	if join.GuildID == guildID {
		for _, role := range joinRoles {
			go addRole(bot, join.User.ID, role, false)
		}
	}
}

func onAddReactionRoleHandler(bot *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if (reaction.GuildID == guildID) && (reaction.ChannelID == roleChannelID) {
		if reaction.MessageID == roleMessageID {
			role := roles[reaction.Emoji.Name]
			go addRole(bot, reaction.UserID, role, true)
		}
		if reaction.MessageID == joinRoleMessageID {
			role := joinRoles[reaction.Emoji.Name]
			go addRole(bot, reaction.UserID, role, true)
		}
	}
}

func onRemoveReactionRoleHandler(bot *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	if (reaction.GuildID == guildID) && (reaction.ChannelID == roleChannelID) {
		if reaction.MessageID == roleMessageID {
			role := roles[reaction.Emoji.Name]
			go removeRole(bot, reaction.UserID, role)
		}
		if reaction.MessageID == joinRoleMessageID {
			role := joinRoles[reaction.Emoji.Name]
			go removeRole(bot, reaction.UserID, role)
		}
	}
}

func addRole(bot *discordgo.Session, userID string, role string, isNotice bool) {
	err := bot.GuildMemberRoleAdd(guildID, userID, role)
	if err != nil {
		fmt.Println("error add role,", err)
		return
	}
	if isNotice {
		sendDM(bot, userID, "ロールを付与しました✅")
	}
}

func removeRole(bot *discordgo.Session, userID string, role string) {
	err := bot.GuildMemberRoleRemove(guildID, userID, role)
	if err != nil {
		fmt.Println("error remove role,", err)
		return
	}
	sendDM(bot, userID, "ロールを削除しました✅")
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
