// Event multiplexer.
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

// Context information passed to the handlers
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

// Structure the multiplexer
type Multiplexer struct {
	Routes []*Route
	Prefix string
}

// Function signature for functions that handle commands
type CommandHandler func(*Context)

// Information about a handler
type Route struct {
	Pattern       string
	AliasPatterns []string
	Description   string
	Category      *CommandCategory
	Handler       CommandHandler
}

// Command categories
type CommandCategory struct {
	Routes      []*Route
	Title       string
	Description string
}

// Some structures to save some registering work
type CommandHandlers struct{}

func init() {
	EventHandlers = append(EventHandlers, Router.OnMessageCreate)
	EventHandlers = append(EventHandlers, Router.OnGuildMemberAdd)
	EventHandlers = append(EventHandlers, Router.OnGuildMemberRemove)
	EventHandlers = append(EventHandlers, Router.OnGuildDelete)
}

// Returns a new command category
func NewCategory(name string, description string) *CommandCategory {
	cat := &CommandCategory{
		Title:       name,
		Description: description,
	}
	return cat
}

// Returns a new message route multiplexer
func New() *Multiplexer {
	return &Multiplexer{
		Prefix: config.Config.System.Prefix,
	}
}

// Register a route
func (mux *Multiplexer) Route(route *Route) *Route {
	route.Category.Routes = append(route.Category.Routes, route)
	mux.Routes = append(mux.Routes, route)
	return route
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
			} else {
				// Attempt caching the channel
				err = session.State.GuildAdd(guild)
				if err != nil {
					log.Warnf("Failed to cache channel fetched from API, %s", err)
				}
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
		} else {
			// Attempt caching the channel
			err = session.State.ChannelAdd(channel)
			if err != nil {
				log.Warnf("Failed to cache channel fetched from API, %s", err)
			}
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
	guildPrefix := context.GenerateGuildPrefix()

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
				if function, success := hook.(func(context *Context)); success {
					function(context)
				}
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
	route, fields := mux.MatchRoute(context.Content)
	if route != nil {
		context.Fields = fields
		route.Handler(context)
		return
	}

	// If no command was matched, resort to either being annoyed by the ping or a command not found message
	if context.HasMention {
		_, _ = session.ChannelMessageSend(channel.ID, "<a:KyokoAngryPing:757399059114885180>")
	} else {
		_, _ = session.ChannelMessageSend(channel.ID,
			fmt.Sprintf("This command does not exist! Issue `%sman` for a list of command manuals.",
				guildPrefix))
	}
}

// Event handler that fires when a guild member is added
func (mux *Multiplexer) OnGuildMemberAdd(session *discordgo.Session, add *discordgo.GuildMemberAdd) {
	go func() {
		for _, hook := range GuildMemberAdd {
			if function, success := hook.(func(session *discordgo.Session, add *discordgo.GuildMemberAdd)); success {
				function(session, add)
			}
		}
	}()
	return
}

// Event handler that fires when a guild member is removed
func (mux *Multiplexer) OnGuildMemberRemove(session *discordgo.Session, remove *discordgo.GuildMemberRemove) {
	go func() {
		for _, hook := range GuildMemberRemove {
			if function, success := hook.(func(session *discordgo.Session, remove *discordgo.GuildMemberRemove)); success {
				function(session, remove)
			}
		}
	}()
	return
}

// Event handler that fires when a guild is deleted
func (mux *Multiplexer) OnGuildDelete(session *discordgo.Session, delete *discordgo.GuildDelete) {
	go func() {
		for _, hook := range GuildDelete {
			if function, success := hook.(func(session *discordgo.Session, delete *discordgo.GuildDelete)); success {
				function(session, delete)
			}
		}
	}()
	return
}
