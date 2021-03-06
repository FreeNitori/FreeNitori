package routes

import (
	"fmt"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"strings"
)

func init() {
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "man",
		AliasPatterns: []string{"manuals", "help"},
		Description:   "An interface to the system reference manuals.",
		Category:      multiplexer.ManualsCategory,
		Handler:       manuals,
	})
}

func manuals(context *multiplexer.Context) {
	guildPrefix := context.Prefix()

	switch {
	case len(context.Fields) == 1:
		{
			// Generate a list of all categories
			embed := embedutil.New("Manuals",
				fmt.Sprintf("Issue `%sman <category>` for category-specific manuals.", guildPrefix))
			embed.Color = 0x3492c4

			// The block of text with all categories
			var catText string
			for _, category := range state.Multiplexer.Categories {

				// Only display categories with description set
				if category.Description == "" {
					continue
				}

				catText += fmt.Sprintf("%s - %s \n", category.Title, category.Description)

			}

			// Add list of categories to the Embed
			embed.AddField("Categories", catText, false)
			_ = context.SendEmbed("", embed)
		}

	case len(context.Fields) == 2:
		{

			// Figure out if the category exist, and fallthrough if it doesn't
			var desiredCat *multiplexer.CommandCategory
			for _, cat := range state.Multiplexer.Categories {
				if strings.EqualFold(cat.Title, context.Fields[1]) {
					desiredCat = cat
					break
				}
			}

			// Break out of the case if no category was matched
			if desiredCat == nil {
				context.SendMessage(multiplexer.InvalidArgument)
				break
			}

			// Generate list of all commands in one specific category
			embed := embedutil.New(desiredCat.Title,
				desiredCat.Description)
			embed.Color = multiplexer.KappaColor

			for _, route := range desiredCat.Routes {

				// Only display routes with non-empty description
				if route.Description == "" {
					continue
				}

				var aliases string
				if len(route.AliasPatterns) > 0 {
					aliases = " (alias patterns:"
					for _, alias := range route.AliasPatterns {
						aliases += " `" + alias + "`"
					}
					aliases += ")"
				}
				embed.AddField(route.Pattern+aliases, route.Description, false)
			}
			context.SendEmbed("", embed)

		}

	case len(context.Fields) > 2:
		{
			// Some catch-all case I guess, though there will be a command-specific thing later maybe
			context.SendMessage(multiplexer.InvalidArgument)
		}
	}
}
