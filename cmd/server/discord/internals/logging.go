package internals

import (
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/discord/sessioning"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/embedutil"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/overrides"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
)

func init() {
	multiplexer.MessageDelete = append(multiplexer.MessageDelete, messageDeleteLog)
	multiplexer.MessageUpdate = append(multiplexer.MessageUpdate, messageUpdateLog)
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

func messageDeleteLog(session *discordgo.Session, delete *discordgo.MessageDelete) {
	if delete.BeforeDelete == nil {
		return
	}
	if delete.BeforeDelete.Author.ID == state.RawSession.State.User.ID {
		return
	}
	var embed = embedutil.NewEmbed("Message Delete", "")
	channelID, err := config.GetGuildConfValue(delete.GuildID, "log_channel")
	if err != nil {
		return
	}
	if channelID == "" {
		return
	}
	if sessioning.FetchChannel(sessioning.FetchGuild(delete.GuildID), "", channelID) == nil {
		return
	}
	embed.Color = state.KappaColor
	embed.SetAuthor(delete.BeforeDelete.Author.Username+"#"+delete.BeforeDelete.Author.Discriminator, delete.BeforeDelete.Author.AvatarURL("128"))
	embed.SetFooter(fmt.Sprintf("Channel: %s Message: %s Author: %s", delete.ChannelID, delete.BeforeDelete.ID, delete.BeforeDelete.Author.ID))
	if delete.BeforeDelete.Content != "" {
		embed.AddField("Content Pre", delete.BeforeDelete.Content, false)
	}
	for _, attachment := range delete.BeforeDelete.Attachments {
		embed.AddField("Attachment Pre", fmt.Sprintf("[%s](%s)", attachment.Filename, attachment.URL), false)
	}
	embed.AddField("Channel", fmt.Sprintf("<#%s>", delete.ChannelID), false)
	context := &multiplexer.Context{Message: &discordgo.Message{ChannelID: channelID}, Session: session}
	context.SendEmbed("", embed)
	for _, e := range delete.BeforeDelete.Embeds {
		context.SendEmbed("Embed included in previously deleted message.", embedutil.Embed{MessageEmbed: e})
	}
}

func messageUpdateLog(session *discordgo.Session, update *discordgo.MessageUpdate) {
	if update.BeforeUpdate == nil {
		return
	}
	if update.BeforeUpdate.Author.ID == state.RawSession.State.User.ID {
		return
	}
	if update.Author == nil {
		return
	}
	var embed = embedutil.NewEmbed("Message Update",
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
	embed.Color = state.KappaColor
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
	context := &multiplexer.Context{Message: &discordgo.Message{ChannelID: channelID}, Session: session}
	context.SendEmbed("", embed)
	for _, e := range update.BeforeUpdate.Embeds {
		context.SendEmbed("Embed included in previously updated message.", embedutil.Embed{MessageEmbed: e})
	}
}
