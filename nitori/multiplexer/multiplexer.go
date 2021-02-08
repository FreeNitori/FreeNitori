// Package multiplexer does command related stuff.
package multiplexer

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strconv"
	"strings"
)

// Context holds information of an event.
type Context struct {
	Message           *discordgo.Message
	Session           *discordgo.Session
	Guild             *discordgo.Guild
	Author            *discordgo.User
	Create            *discordgo.MessageCreate
	Fields            []string
	Content           string
	IsPrivate         bool
	IsTargeted        bool
	HasPrefix         bool
	HasMention        bool
	HasLeadingMention bool
}

// Multiplexer represents the command router.
type Multiplexer struct {
	Routes []*Route
	Prefix string
}

// CommandHandler represents the handler function of a Route.
type CommandHandler func(*Context)

// Route represents a command route.
type Route struct {
	Pattern       string
	AliasPatterns []string
	Description   string
	Category      *CommandCategory
	Handler       CommandHandler
}

// CommandCategory represents a category of Route.
type CommandCategory struct {
	Routes      []*Route
	Title       string
	Description string
}

func init() {
	EventHandlers = append(EventHandlers,
		Router.onReady,
		Router.handleMessage,
		Router.onGuildMemberAdd,
		Router.onGuildMemberRemove,
		Router.onGuildDelete,
		Router.onMessageCreate,
		Router.onMessageDelete,
		Router.onMessageReactionAdd,
		Router.onMessageReactionRemove)
}

// NewCategory returns a new command category
func NewCategory(name string, description string) *CommandCategory {
	cat := &CommandCategory{
		Title:       name,
		Description: description,
	}
	return cat
}

// New returns a router.
func New() *Multiplexer {
	return &Multiplexer{
		Prefix: config.Config.System.Prefix,
	}
}

// Route registers a route to the router.
func (mux *Multiplexer) Route(route *Route) *Route {
	route.Category.Routes = append(route.Category.Routes, route)
	mux.Routes = append(mux.Routes, route)
	return route
}

// MatchRoute fuzzy matches a message to a route.
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

func (mux *Multiplexer) handleMessage(session *discordgo.Session, create *discordgo.MessageCreate) {
	var err error

	// Ignore self and bot messages
	if create.Author.ID == session.State.User.ID || create.Author.Bot {
		return
	}

	// Add to the counter if the message is valid
	err = config.AdvanceTotalMessages()
	if err != nil {
		log.Warnf("Failed to increase the message counter, %s", err)
	}

	// Figure out the message guild
	var guild *discordgo.Guild
	if create.GuildID != "" {
		guild, err = session.State.Guild(create.GuildID)
		if err != nil {
			// Attempt direct API fetching
			guild, err = session.Guild(create.GuildID)
			if err != nil {
				log.Errorf("Failed to fetch guild from API or cache, %s", err)
				return
			}
			// Attempt caching the channel
			err = session.State.GuildAdd(guild)
			if err != nil {
				log.Warnf("Failed to cache channel fetched from API, %s", err)
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
			log.Errorf("Failed to fetch channel from API or cache, %s", err)
			return
		}
		// Attempt caching the channel
		err = session.State.ChannelAdd(channel)
		if err != nil {
			log.Warnf("Failed to cache channel fetched from API, %s", err)
		}
	}

	// Make context info for route
	context := &Context{
		Content: strings.TrimSpace(create.Content),
		Message: create.Message,
		Session: session,
		Author:  create.Author,
		Create:  create,
		Guild:   guild,
	}

	// Put the channel into context
	if channel != nil {
		if channel.Type == discordgo.ChannelTypeDM {
			context.IsPrivate = true
		}
	}

	// Get guild-specific prefix
	guildPrefix := context.Prefix()

	// Figure out if the Kappa got pinged
	if !context.IsTargeted {
		for _, mentionedUser := range create.Mentions {
			if mentionedUser.ID == session.State.User.ID {
				context.IsTargeted, context.HasMention = true, true
				mentionRegex := regexp.MustCompile(fmt.Sprintf("<@!?(%s)>", session.State.User.ID))
				location := mentionRegex.FindStringIndex(context.Content)

				// Figure out if the message started with the ping
				if len(location) == 0 {
					context.HasLeadingMention = true
				} else if location[0] == 0 {
					context.HasLeadingMention = true
				}

				context.Content = mentionRegex.ReplaceAllString(context.Content, "")

				break
			}
		}
	}

	// Figure out if a proper command is issued to the Kappa
	if !context.IsTargeted && len(guildPrefix) > 0 {
		if strings.HasPrefix(context.Content, guildPrefix) {
			context.IsTargeted, context.HasPrefix = true, true
			context.Content = strings.TrimPrefix(context.Content, guildPrefix)
		}
	}

	// Get out of the code if no one targeted the Kappa
	if !context.IsTargeted {
		// Run all the not targeted hooks and leave
		go func() {
			for _, hook := range NotTargeted {
				hook(context)
			}
		}()
		return
	}

	// Log the processed message
	var hostName string
	if context.IsPrivate {
		hostName = "Private Messages"
	} else {
		hostName = "\"" + context.Guild.Name + "\""
	}
	log.Infof("(Shard %s) \"%s\"@%s > %s",
		strconv.Itoa(session.ShardID),
		context.Author.Username+"#"+context.Author.Discriminator,
		hostName,
		context.Message.Content)

	// Figure out the route of the message
	if !(context.HasMention && !context.HasLeadingMention) {
		route, fields := mux.MatchRoute(context.Content)
		if route != nil {
			context.Fields = fields
			route.Handler(context)
			return
		}
	}

	// If no command was matched, resort to either being annoyed by the ping or a command not found message
	if context.HasMention {
		context.SendMessage("<a:KyokoAngryPing:757399059114885180>")
	} else {
		context.SendMessage(fmt.Sprintf("This command does not exist! Issue `%sman` for a list of command manuals.",
			guildPrefix))
	}
}

// Event handler that fires when ready
func (mux *Multiplexer) onReady(session *discordgo.Session, ready *discordgo.Ready) {
	go func() {
		for _, hook := range Ready {
			hook(session, ready)
		}
	}()
	return
}

// Event handler that fires when a guild member is added
func (mux *Multiplexer) onGuildMemberAdd(session *discordgo.Session, add *discordgo.GuildMemberAdd) {
	go func() {
		for _, hook := range GuildMemberAdd {
			hook(session, add)
		}
	}()
	return
}

// Event handler that fires when a guild member is removed
func (mux *Multiplexer) onGuildMemberRemove(session *discordgo.Session, remove *discordgo.GuildMemberRemove) {
	go func() {
		for _, hook := range GuildMemberRemove {
			hook(session, remove)
		}
	}()
	return
}

// Event handler that fires when a guild is deleted
func (mux *Multiplexer) onGuildDelete(session *discordgo.Session, delete *discordgo.GuildDelete) {
	go func() {
		for _, hook := range GuildDelete {
			hook(session, delete)
		}
	}()
	return
}

// Event handler that fires when a message is created
func (mux *Multiplexer) onMessageCreate(session *discordgo.Session, create *discordgo.MessageCreate) {
	go func() {
		for _, hook := range MessageCreate {
			hook(session, create)
		}
	}()
	return
}

// Event handler that fires when a message is deleted
func (mux *Multiplexer) onMessageDelete(session *discordgo.Session, delete *discordgo.MessageDelete) {
	go func() {
		for _, hook := range MessageDelete {
			hook(session, delete)
		}
	}()
	return
}

// Event handler that fires when a reaction is added
func (mux *Multiplexer) onMessageReactionAdd(session *discordgo.Session, add *discordgo.MessageReactionAdd) {
	go func() {
		for _, hook := range MessageReactionAdd {
			hook(session, add)
		}
	}()
	return
}

// Event handler that fires when a reaction is removed
func (mux *Multiplexer) onMessageReactionRemove(session *discordgo.Session, remove *discordgo.MessageReactionRemove) {
	go func() {
		for _, hook := range MessageReactionRemove {
			hook(session, remove)
		}
	}()
}
