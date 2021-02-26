package internals

import (
	"fmt"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/discord/sessioning"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/overrides"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"github.com/bwmarrin/discordgo"
)

func init() {
	state.Multiplexer.MessageDelete = append(state.Multiplexer.MessageDelete, messageDeleteLog)
	state.Multiplexer.MessageUpdate = append(state.Multiplexer.MessageUpdate, messageUpdateLog)
	overrides.RegisterSimpleEntry(overrides.SimpleConfigurationEntry{
		Name:         "logging",
		FriendlyName: "Logging",
		Description:  "Configure logging system.",
		DatabaseKey:  "log_channel",
		Cleanup:      func(context *multiplexer.Context) {},
		Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
			if channel := context.GetChannel(*input); channel != nil {
				*input = channel.ID
				return true, true
			}
			return false, true
		},
		Format: func(context *multiplexer.Context, value string) (string, string, bool) {
			if channel := context.GetChannel(value); channel != nil {
				return channel.Name, channel.ID, true
			}
			return "No channel configured", fmt.Sprintf("Configure it by issuing command `%sconf logging <channel>`.\n"+
				"Setting this will enable message logging.", context.Prefix()), true
		},
	})
}

func messageDeleteLog(context *multiplexer.Context) {
	messageDelete, ok := context.Event.(*discordgo.MessageDelete)
	if !ok {
		return
	}
	if messageDelete.GuildID == "" {
		return
	}
	if messageDelete.BeforeDelete == nil {
		return
	}
	if messageDelete.BeforeDelete.Author.ID == state.RawSession.State.User.ID {
		return
	}
	var embed = embedutil.New("Message Delete", "")
	channelID, err := config.GetGuildConfValue(messageDelete.GuildID, "log_channel")
	if err != nil {
		return
	}
	if channelID == "" {
		return
	}
	if sessioning.FetchChannel(sessioning.FetchGuild(messageDelete.GuildID), "", channelID) == nil {
		return
	}
	embed.Color = multiplexer.KappaColor
	embed.SetAuthor(messageDelete.BeforeDelete.Author.Username+"#"+messageDelete.BeforeDelete.Author.Discriminator, messageDelete.BeforeDelete.Author.AvatarURL("128"))
	embed.SetFooter(fmt.Sprintf("Channel: %s Message: %s Author: %s", messageDelete.ChannelID, messageDelete.BeforeDelete.ID, messageDelete.BeforeDelete.Author.ID))
	if messageDelete.BeforeDelete.Content != "" {
		embed.AddField("Content Pre", messageDelete.BeforeDelete.Content, false)
	}
	for _, attachment := range messageDelete.BeforeDelete.Attachments {
		embed.AddField("Attachment Pre", fmt.Sprintf("[%s](%s)", attachment.Filename, attachment.URL), false)
	}
	embed.AddField("Channel", fmt.Sprintf("<#%s>", messageDelete.ChannelID), false)
	context.Message = &discordgo.Message{ChannelID: channelID}
	context.SendEmbed("", embed)
	for _, e := range messageDelete.BeforeDelete.Embeds {
		context.SendEmbed("Embed included in previously deleted message.", embedutil.Embed{MessageEmbed: e})
	}
}

func messageUpdateLog(context *multiplexer.Context) {
	update, ok := context.Event.(*discordgo.MessageUpdate)
	if !ok {
		return
	}
	if update.GuildID == "" {
		return
	}
	if update.BeforeUpdate == nil {
		return
	}
	if update.BeforeUpdate.Author.ID == state.RawSession.State.User.ID {
		return
	}
	if update.Author == nil {
		return
	}
	var embed = embedutil.New("Message Update",
		fmt.Sprintf("[Message Link](https://discord.com/channels/%s/%s/%s)",
			update.BeforeUpdate.GuildID,
			update.BeforeUpdate.ChannelID,
			update.BeforeUpdate.ID))
	channelID, err := config.GetGuildConfValue(update.GuildID, "log_channel")
	if err != nil {
		return
	}
	if channelID == "" {
		return
	}
	if sessioning.FetchChannel(sessioning.FetchGuild(update.GuildID), "", channelID) == nil {
		return
	}
	embed.Color = multiplexer.KappaColor
	embed.SetAuthor(update.BeforeUpdate.Author.Username+"#"+update.BeforeUpdate.Author.Discriminator, update.BeforeUpdate.Author.AvatarURL("128"))
	embed.SetFooter(fmt.Sprintf("Channel: %s Message: %s Author: %s", update.ChannelID, update.BeforeUpdate.ID, update.BeforeUpdate.Author.ID))
	if update.BeforeUpdate.Content != "" {
		embed.AddField("Content Pre", update.BeforeUpdate.Content, false)
	}
	for _, attachment := range update.BeforeUpdate.Attachments {
		embed.AddField("Attachment Pre", fmt.Sprintf("[%s](%s)", attachment.Filename, attachment.URL), false)
	}
	if update.Message.Content != "" {
		embed.AddField("Content Post", update.Message.Content, false)
	}
	for _, attachment := range update.Message.Attachments {
		embed.AddField("Attachment Post", fmt.Sprintf("[%s](%s)", attachment.Filename, attachment.URL), false)
	}
	embed.AddField("Channel", fmt.Sprintf("<#%s>", update.ChannelID), false)
	context.Message = &discordgo.Message{ChannelID: channelID}
	context.SendEmbed("", embed)
	for _, e := range update.BeforeUpdate.Embeds {
		context.SendEmbed("Embed included in previously updated message.", embedutil.Embed{MessageEmbed: e})
	}
}
