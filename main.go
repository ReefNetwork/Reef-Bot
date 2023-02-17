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
	token   string
	guildID = "638760361369010177"

	roles = map[string]string{
		"💛": "918889974399639614",
		"💚": "889077232209895454",
		"💙": "936600142637830184",
		"💜": "1076051118230089808",
	}

	joinRoles = map[string]string{
		"⛏": "982904224532815932",
	}
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

	bot.AddHandler(onJoinRoleHandler)
	bot.AddHandler(onAddReactionRoleHandler)
	bot.AddHandler(onRemoveReactionRoleHandler)

	bot.AddHandler(onAddReactionQuestionHandler)

	bot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	initReaction(bot)
	fmt.Println("Bot is now running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	//goland:noinspection GoUnhandledErrorResult
	bot.Close()
}
