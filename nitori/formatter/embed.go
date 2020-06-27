package formatter

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type Embed struct {
	*discordgo.MessageEmbed
}

// Length limits for any embed
const (
	EmbedLimitTitle       = 256
	EmbedLimitDescription = 2048
	EmbedLimitFieldName   = 256
	EmbedLimitFieldValue  = 1024
	EmbedLimitField       = 25
	EmbedLimitFooter      = 2048
)

// Make a new Embed object with specific things
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

// Append a field to the Embed
func (embed *Embed) AddField(name, value string, inline bool) *Embed {

	if len(embed.Fields) == EmbedLimitField {
		log.Printf("Embed with name \"%s\" exceeded limit!", name)
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

// Set the Embed's footer string and image
func (embed *Embed) SetFooter(args ...string) *Embed {
	var (
		iconURL      string
		text         string
		proxyIconURL string
	)

	switch {
	case len(args) > 2:
		proxyIconURL = args[2]
		fallthrough
	case len(args) > 1:
		iconURL = args[1]
		fallthrough
	case len(args) > 0:
		text = args[0]
	case len(args) == 0:
		return embed
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

// Set image URL of the Embed
func (embed *Embed) SetImage(args ...string) *Embed {
	var (
		URL      string
		proxyURL string
	)

	switch {
	case len(args) > 1:
		proxyURL = args[1]
		fallthrough
	case len(args) > 0:
		URL = args[0]
		fallthrough
	case len(args) == 0:
		return embed
	}

	embed.Image = &discordgo.MessageEmbedImage{
		URL:      URL,
		ProxyURL: proxyURL,
	}

	return embed
}

// Set thumbnail URL of the Embed
func (embed *Embed) SetThumbnail(args ...string) *Embed {
	var (
		URL      string
		proxyURL string
	)

	switch {
	case len(args) > 1:
		proxyURL = args[1]
		fallthrough
	case len(args) > 0:
		URL = args[0]
		fallthrough
	case len(args) == 0:
		return embed
	}

	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:      URL,
		ProxyURL: proxyURL,
	}

	return embed
}

// Set author information of the Embed
func (embed *Embed) SetAuthor(args ...string) *Embed {
	var (
		name     string
		iconURL  string
		URL      string
		proxyURL string
	)

	switch {
	case len(args) > 3:
		proxyURL = args[3]
		fallthrough
	case len(args) > 2:
		URL = args[2]
		fallthrough
	case len(args) > 1:
		iconURL = args[1]
		fallthrough
	case len(args) == 0:
		return embed

	}

	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:         name,
		IconURL:      iconURL,
		URL:          URL,
		ProxyIconURL: proxyURL,
	}

	return embed
}
