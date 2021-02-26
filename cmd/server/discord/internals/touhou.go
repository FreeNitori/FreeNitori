package internals

import (
	"fmt"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	imagefetch "git.randomchars.net/FreeNitori/ImageFetch"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"strings"
	"time"
)

const (
	noCharacterMatch = "No character matched your query, displaying random character."
	noArtAvailable   = "No art available for this character."
)

var sessions = map[string]chan [2]string{}

func init() {
	state.Multiplexer.NotTargeted = append(state.Multiplexer.NotTargeted, guessResponse)
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "touhou",
		AliasPatterns: []string{"t", "th"},
		Description:   "Finds picture of requested character.",
		Category:      multiplexer.MediaCategory,
		Handler:       touhou,
	})
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "guess",
		AliasPatterns: []string{},
		Description:   "Guess character based on artwork.",
		Category:      multiplexer.MediaCategory,
		Handler:       guess,
	})
}

func touhou(context *multiplexer.Context) {
	var character imagefetch.CharacterInfo
	var art imagefetch.CharacterArt
	var ok bool
	var text string
	var randomArt = func() {
		character, art, err = imagefetch.FetchRandom(imagefetch.SafeTouhouQuery, imagefetch.Touhou)
		if err == imagefetch.ErrNoArtAvailable {
			context.SendMessage(noArtAvailable)
			return
		}
		if !context.HandleError(err) {
			return
		}
	}

	if len(context.Fields) > 1 {
		character, ok = imagefetch.Touhou.CharacterFriendly(context.Fields[1])
		if !ok {
			context.SendMessage(noCharacterMatch)
			randomArt()
		} else {
			art, err = imagefetch.Fetch(imagefetch.SafeTouhouQuery, character)
			if err == imagefetch.ErrNoArtAvailable {
				context.SendMessage(noArtAvailable)
				return
			}
			if !context.HandleError(err) {
				return
			}
		}
	} else {
		randomArt()
	}

	embed := embedutil.New("", "")
	embed.Color = character.Color
	embed.SetImage(art.ImageURL)
	embed.SetAuthor(character.FriendlyName)
	embed.SetFooter("Source URL: " + art.SourceURL)
	context.SendEmbed(text, embed)
}

func guessResponse(context *multiplexer.Context) {
	channel, ok := sessions[context.Message.ChannelID]
	if !ok {
		return
	}
	channel <- [2]string{context.Text, context.User.Mention()}
}

func guess(context *multiplexer.Context) {
	if context.IsPrivate {
		context.SendMessage(multiplexer.GuildOnly)
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
	var character imagefetch.CharacterInfo
	character, art, err = imagefetch.FetchRandom(imagefetch.SafeTouhouQuery, imagefetch.Touhou)
	if err == imagefetch.ErrNoArtAvailable {
		context.SendMessage("No art available for this character.")
		return
	}
	if !context.HandleError(err) {
		return
	}

	embed := embedutil.New("Guess Character", "You have 15 seconds to decide.")
	embed.Color = character.Color
	embed.SetImage(art.ImageURL)
	context.SendEmbed("", embed)

	end := make(chan bool, 1)
	go func() { time.Sleep(15 * time.Second); end <- true }()

	for {
		select {
		case <-end:
			context.SendMessage(fmt.Sprintf("Time's up, the character is %s.", character.FriendlyName))
			return
		case msg := <-message:
			if strings.ToLower(msg[0]) == strings.ToLower(character.FriendlyName) ||
				strings.ToLower(msg[0]) == strings.Replace(strings.ToLower(character.SearchString), "_", " ", -1) {
				context.SendMessage(fmt.Sprintf("%s correct! The character is %s.", msg[1], character.FriendlyName))
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
