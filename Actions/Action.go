package actions

import (
	"github.com/bwmarrin/discordgo"

	"github.com/soliel/BulliBot/command"
)

//ChatAction is a command that runs independently of normal command handling when a certain pre-requisite is met..
type ChatAction struct {
	predicate func(*discordgo.MessageCreate) bool
	action    func(command.Context)
}
