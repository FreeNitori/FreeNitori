package multiplexer

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/utils"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strings"
)

// Context information passed to the handlers
type Context struct {
	Fields            []string
	Content           string
	IsPrivate         bool
	IsTargeted        bool
	HasPrefix         bool
	HasMention        bool
	HasLeadingMention bool
}

// Function signature for functions that handle commands
type CommandHandler func(*discordgo.Session, *discordgo.Message, *Context)

// Information about a handler
type Route struct {
	Pattern     string
	Description string
	Manuals     string
	Handler     CommandHandler
}

// Structure for all multiplexer things
type Multiplexer struct {
	Routes  []*Route
	Default *Route
	Prefix  string
}

// Returns a new message route multiplexer
func New() *Multiplexer {
	mux := &Multiplexer{}
	mux.Prefix = utils.Prefix
	return mux
}

// Register a route
func (mux *Multiplexer) Route(pattern, description string, handler CommandHandler) (*Route, error) {
	route := Route{}
	route.Pattern = pattern
	route.Description = description
	route.Handler = handler

	return &route, nil
}

// This matches routes for the message
func (mux *Multiplexer) MatchRoute(message string) (*Route, []string) {
	// Make a slice of words out of the message
	fields := strings.Fields(message)

	// Get out if there is nothing
	if len(fields) == 0 {
		return nil, nil
	}

	// Find a route from what we already have
	var route *Route
	var similarityRating int
	var fieldIndex int

	for fieldIndex, fieldIter := range fields {
		for _, routeIter := range mux.Routes {
			// If an exact match was found, immediately give that to the Kappa
			if routeIter.Pattern == fieldIter {
				return routeIter, fields[fieldIndex:]
			}

			// If that's not found just return some random shit
			if strings.HasPrefix(routeIter.Pattern, fieldIter) {
				if len(fieldIter) > similarityRating {
					route = routeIter
					similarityRating = len(fieldIter)
				}
			}
		}
	}
	return route, fields[fieldIndex:]
}

// DiscordGo library event handler registered into the session
func (mux *Multiplexer) OnMessageCreate(session *discordgo.Session, create *discordgo.MessageCreate) {
	var err error

	// Ignore self messages
	if create.Author.ID == session.State.User.ID {
		return
	}

	// Make context info for route
	context := &Context{
		Content: strings.TrimSpace(create.Content),
	}

	// Figure out the message channel
	var channel *discordgo.Channel
	channel, err = session.State.Channel(create.ChannelID)
	if err != nil {
		// Attempt direct API fetching
		channel, err = session.Channel(create.ChannelID)
		if err != nil {
			log.Printf("Failed to fetch channel from API or cache, %s", err)
		} else {
			// Attempt caching the channel
			err = session.State.ChannelAdd(channel)
			if err != nil {
				log.Printf("Failed to cache channel fetched from API, %s", err)
			}
		}
	}

	// Put the channel into context
	if channel != nil {
		if channel.Type == discordgo.ChannelTypeDM {
			context.IsPrivate, context.IsTargeted = true, true
		}
	}

	// Figure out if the Kappa got pinged
	if !context.IsTargeted {
		for _, mentionedUser := range create.Mentions {
			if mentionedUser.ID == session.State.User.ID {
				context.IsTargeted, context.HasMention = true, true
				mentionRegex := regexp.MustCompile(fmt.Sprintf("<@!?(%s)>", session.State.User.ID))

				// Figure out if the message started with the ping
				if mentionRegex.FindStringIndex(context.Content)[0] == 0 {
					context.HasLeadingMention = true
				}

				// Remove the pings so our Kappa doesn't get mad
				context.Content = mentionRegex.ReplaceAllString(context.Content, "")

				break
			}
		}
	}

	// Figure out if a proper command is issued to the Kappa
	if !context.IsTargeted && len(mux.Prefix) > 0 {
		// TODO: Database integration
		if strings.HasPrefix(context.Content, mux.Prefix) {
			context.IsTargeted, context.HasPrefix = true, true
			context.Content = strings.TrimPrefix(context.Content, mux.Prefix)
		}
	}

	// Get out of the code if no one targeted the Kappa
	if !context.IsTargeted {
		return
	}

	// Figure out the route of the message
	route, fields := mux.MatchRoute(context.Content)
	if route != nil {
		context.Fields = fields
		route.Handler(session, create.Message, context)
		return
	}

	// If no command was matched, resort to either being annoyed by the ping or a command not found message
	if context.HasMention {
		_, _ = session.ChannelMessageSend(channel.ID, "<a:KyoukoAngryPing:710413221927976980>")
	} else {
		_, _ = session.ChannelMessageSend(channel.ID, "Got your message, but this command does not exist!")
	}
}
