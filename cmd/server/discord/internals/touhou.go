package internals

import (
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/embedutil"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	imagefetch "git.randomchars.net/FreeNitori/ImageFetch"
	"math/rand"
	"strings"
	"time"
)

var sessions = map[string]chan [2]string{}

func init() {
	multiplexer.NotTargeted = append(multiplexer.NotTargeted, guessResponse)
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "touhou",
		AliasPatterns: []string{"t", "th"},
		Description:   "Finds picture of requested character.",
		Category:      multiplexer.MediaCategory,
		Handler:       touhou,
	})
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "guess",
		AliasPatterns: []string{},
		Description:   "Guess character based on artwork.",
		Category:      multiplexer.MediaCategory,
		Handler:       guess,
	})
}

func touhou(context *multiplexer.Context) {
	var name string
	var char imagefetch.CharacterInfo
	var text string

	if len(context.Fields) > 1 {
		name = strings.ToLower(context.Fields[1])
	}

	for _, character := range imagefetch.Touhou.Characters() {
		if name == strings.ToLower(character.FriendlyName) {
			char = character
		}
	}

	if char.SearchString == "" {
		text = "No character matched your query, displaying random character."
		char = imagefetch.Touhou.Characters()[rand.Intn(len(imagefetch.Touhou.Characters()))]
	}

	var art imagefetch.CharacterArt
	art, err = imagefetch.Fetch(char)
	if err == imagefetch.ErrNoArtAvailable {
		context.SendMessage("No art available for this character.")
		return
	}
	if !context.HandleError(err) {
		return
	}
	embed := embedutil.NewEmbed("", "")
	embed.Color = char.Color
	embed.SetImage(art.ImageURL)
	embed.SetAuthor(char.FriendlyName)
	embed.SetFooter("Source URL: " + art.SourceURL)
	context.SendEmbed(text, embed)
}

func guessResponse(context *multiplexer.Context) {
	channel, ok := sessions[context.Message.ChannelID]
	if !ok {
		return
	}
	channel <- [2]string{context.Content, context.Author.Mention()}
}

func guess(context *multiplexer.Context) {
	if context.IsPrivate {
		context.SendMessage(state.GuildOnly)
		return
	}

	_, ok := sessions[context.Message.ChannelID]
	if ok {
		context.SendMessage("A guessing session already exists in this channel.")
		return
	}

	message := make(chan [2]string, 1)
	sessions[context.Message.ChannelID] = message
	defer func() { delete(sessions, context.Message.ChannelID) }()

	var art imagefetch.CharacterArt
	char := imagefetch.Touhou.Characters()[rand.Intn(len(imagefetch.Touhou.Characters()))]
	art, err = imagefetch.Fetch(char)
	if err == imagefetch.ErrNoArtAvailable {
		context.SendMessage("No art available for this character.")
		return
	}
	if !context.HandleError(err) {
		return
	}

	embed := embedutil.NewEmbed("Guess Character", "You have 15 seconds to decide.")
	embed.Color = char.Color
	embed.SetImage(art.ImageURL)
	context.SendEmbed("", embed)

	end := make(chan bool, 1)
	go func() { time.Sleep(15 * time.Second); end <- true }()

	for {
		select {
		case <-end:
			context.SendMessage(fmt.Sprintf("Time's up, the character is %s.", char.FriendlyName))
			return
		case msg := <-message:
			if strings.ToLower(msg[0]) == strings.ToLower(char.FriendlyName) ||
				strings.ToLower(msg[0]) == strings.Replace(strings.ToLower(char.SearchString), "_", " ", -1) {
				context.SendMessage(fmt.Sprintf("%s correct! The character is %s.", msg[1], char.FriendlyName))
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
