package formatter

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"github.com/bwmarrin/discordgo"
)

type Embed struct {
	*discordgo.MessageEmbed
}

const (
	EmbedLimitTitle       = 256
	EmbedLimitDescription = 2048
	EmbedLimitFieldName   = 256
	EmbedLimitFieldValue  = 1024
	EmbedLimitField       = 25
	EmbedLimitFooter      = 2048
)

// NewEmbed makes a new Embed object.
func NewEmbed(title string, description string) *Embed {
	embed := Embed{&discordgo.MessageEmbed{}}

	if len(title) > EmbedLimitTitle {
		title = title[:EmbedLimitTitle]
	}

	if len(description) > EmbedLimitDescription {
		description = description[:EmbedLimitDescription]
	}

	embed.Title = title
	embed.Description = description
	return &embed
}

// AddField adds a field to the embed.
func (embed *Embed) AddField(name, value string, inline bool) *Embed {

	if len(embed.Fields) == EmbedLimitField {
		log.Warnf("Embed with name \"%s\" exceeded limit!", name)
		return embed
	}

	if len(value) > EmbedLimitFieldValue {
		value = value[:EmbedLimitFieldValue]
	}

	if len(name) > EmbedLimitFieldName {
		name = name[:EmbedLimitFieldName]
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})

	return embed
}

// SetFooter sets the footer text and image of an embed.
func (embed *Embed) SetFooter(args ...string) *Embed {
	var (
		iconURL      string
		text         string
		proxyIconURL string
	)

	switch argsLength := len(args); {
	case argsLength > 2:
		proxyIconURL = args[2]
		fallthrough
	case argsLength > 1:
		iconURL = args[1]
		fallthrough
	case argsLength > 0:
		text = args[0]
	case argsLength == 0:
	}

	if len(text) > EmbedLimitFooter {
		text = text[:EmbedLimitFooter]
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		IconURL:      iconURL,
		Text:         text,
		ProxyIconURL: proxyIconURL,
	}

	return embed
}

// SetImage sets the image URL of an embed.
func (embed *Embed) SetImage(args ...string) *Embed {
	var (
		URL      string
		proxyURL string
	)

	switch argsLength := len(args); {
	case argsLength > 1:
		proxyURL = args[1]
		fallthrough
	case argsLength > 0:
		URL = args[0]
		fallthrough
	case argsLength == 0:
	}

	embed.Image = &discordgo.MessageEmbedImage{
		URL:      URL,
		ProxyURL: proxyURL,
	}

	return embed
}

// SetThumbnail sets the thumbnail URL of an embed.
func (embed *Embed) SetThumbnail(args ...string) *Embed {
	var (
		URL      string
		proxyURL string
	)

	switch argsLength := len(args); {
	case argsLength > 1:
		proxyURL = args[1]
		fallthrough
	case argsLength > 0:
		URL = args[0]
		fallthrough
	case argsLength == 0:
	}

	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return embed
}

// SetAuthor sets author name, URL and icon URL of an embed.
func (embed *Embed) SetAuthor(args ...string) *Embed {
	var (
		name     string
		iconURL  string
		URL      string
		proxyURL string
	)

	switch argsLength := len(args); {
	case argsLength > 3:
		proxyURL = args[3]
		fallthrough
	case argsLength > 2:
		URL = args[2]
		fallthrough
	case argsLength > 1:
		iconURL = args[1]
		fallthrough
	case argsLength > 0:
		name = args[0]
		fallthrough
	case argsLength == 0:
	}

	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:         name,
		IconURL:      iconURL,
		URL:          URL,
		ProxyIconURL: proxyURL,
	}

	return embed
}
