package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"github.com/bwmarrin/discordgo"
)

func init() {
	multiplexer.MessageReactionAdd = append(multiplexer.MessageReactionAdd, addReaction)
	multiplexer.MessageReactionRemove = append(multiplexer.MessageReactionRemove, removeReaction)
}

func addReaction(session *discordgo.Session, add *discordgo.MessageReactionAdd) {
	// TODO: Implement add reaction handling
}

func removeReaction(session *discordgo.Session, remove *discordgo.MessageReactionRemove) {
	// TODO: Implement remove reaction handling
}
