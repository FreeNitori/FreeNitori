package multiplexer

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var Router = New()

// Context information passed to the handlers
type Context struct {
	Message           *discordgo.Message
	Session           *discordgo.Session
	Guild             *discordgo.Guild
	Author            *discordgo.User
	Fields            []string
	Content           string
	IsPrivate         bool
	IsTargeted        bool
	HasPrefix         bool
	HasMention        bool
	HasLeadingMention bool
}

// Wrapper around some stuff
func (context *Context) GenerateGuildPrefix() (guildID int, guildPrefix string) {
	var gID int
	var gPrefix string
	if !context.IsPrivate {
		gID, _ = strconv.Atoi(context.Guild.ID)
		gPrefix = config.GetPrefix(guildID)
	} else {
		gID = 0
		gPrefix = config.Prefix
	}
	return gID, gPrefix
}

// Function signature for functions that handle commands
type CommandHandler func(*Context)

// Information about a handler
type Route struct {
	Pattern       string
	AliasPatterns []string
	Description   string
	Manuals       string
	Category      CommandCategory
	Handler       CommandHandler
}

// Command categories
type CommandCategory struct {
	Routes      []*Route
	Title       string
	Description string
}

// Returns a new command category
func NewCategory(name string, description string) *CommandCategory {
	cat := &CommandCategory{
		Title:       name,
		Description: description,
	}
	return cat
}

// Structure for all multiplexer things
type Multiplexer struct {
	Routes []*Route
	Prefix string
}

// Returns a new message route multiplexer
func New() *Multiplexer {
	mux := &Multiplexer{
		Prefix: config.Prefix,
	}
	return mux
}

// Register a route
func (mux *Multiplexer) Route(pattern string, aliasPattern []string, description string, handler CommandHandler, category *CommandCategory) *Route {
	route := Route{
		Pattern:       pattern,
		AliasPatterns: aliasPattern,
		Description:   description,
		Handler:       handler,
		Category:      *category,
	}
	category.Routes = append(category.Routes, &route)
	mux.Routes = append(mux.Routes, &route)
	return &route
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
			// If an exact match is found, immediately give that to the kappa
			if routeIter.Pattern == fieldIter {
				return routeIter, fields[fieldIndex:]
			}

			// If an exact match on any alias is found, immediately give that to the kappa
			for _, aliasPattern := range routeIter.AliasPatterns {
				if aliasPattern == fieldIter {
					return routeIter, fields[fieldIndex:]
				}
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

	// Add to the counter if the message is valid
	err = config.AddTotalMessages()
	if err != nil {
		log.Printf("Failed to increase the message counter, %s", err)
	}

	// Figure out the message guild
	var guild *discordgo.Guild
	guild, err = session.State.Guild(create.GuildID)
	if err != nil {
		// Attempt direct API fetching
		guild, err = session.Guild(create.GuildID)
		if err != nil {
			log.Printf("Failed to fetch guild from API or cache, %s", err)
			return
		} else {
			// Attempt caching the channel
			err = session.State.GuildAdd(guild)
			if err != nil {
				log.Printf("Failed to cache channel fetched from API, %s", err)
			}
		}
	}

	// Figure out the message channel
	var channel *discordgo.Channel
	channel, err = session.State.Channel(create.ChannelID)
	if err != nil {
		// Attempt direct API fetching
		channel, err = session.Channel(create.ChannelID)
		if err != nil {
			log.Printf("Failed to fetch channel from API or cache, %s", err)
			return
		} else {
			// Attempt caching the channel
			err = session.State.ChannelAdd(channel)
			if err != nil {
				log.Printf("Failed to cache channel fetched from API, %s", err)
			}
		}
	}

	// Make context info for route
	context := &Context{
		Content: strings.TrimSpace(create.Content),
		Message: create.Message,
		Session: session,
		Author:  create.Author,
		Guild:   guild,
	}

	// Put the channel into context
	if channel != nil {
		if channel.Type == discordgo.ChannelTypeDM {
			context.IsPrivate = true
		}
	}

	// Get guild-specific prefix
	_, guildPrefix := context.GenerateGuildPrefix()

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
	if !context.IsTargeted && len(guildPrefix) > 0 {
		// TODO: Database integration
		if strings.HasPrefix(context.Content, guildPrefix) {
			context.IsTargeted, context.HasPrefix = true, true
			context.Content = strings.TrimPrefix(context.Content, guildPrefix)
		}
	}

	// Get out of the code if no one targeted the Kappa
	if !context.IsTargeted {
		return
	}

	// Tell the supervisor to log the processed message
	var hostName string
	if context.IsPrivate {
		hostName = "Private Messages"
	} else {
		hostName = "\"" + context.Guild.Name + "\""
	}
	WritePacket(IPCConnection,
		IPCPacket{
			IssuerIdentifier:   "ChatBackend",
			ReceiverIdentifier: "Supervisor",
			MessageIdentifier:  "RouteLog",
			Body: []string{
				context.Author.Username + "#" + context.Author.Discriminator,
				hostName,
				context.Message.Content,
				strconv.Itoa(session.ShardID),
			},
		})

	// Figure out the route of the message
	route, fields := mux.MatchRoute(context.Content)
	if route != nil {
		context.Fields = fields
		route.Handler(context)
		return
	}

	// If no command was matched, resort to either being annoyed by the ping or a command not found message
	if context.HasMention {
		_, _ = session.ChannelMessageSend(channel.ID, "<a:KyoukoAngryPing:710413221927976980>")
	} else {
		_, _ = session.ChannelMessageSend(channel.ID,
			fmt.Sprintf("This command does not exist! Issue `%sman` for a list of command manuals.",
				guildPrefix))
	}
}
