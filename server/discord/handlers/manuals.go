package handlers

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"strings"
)

func init() {
	multiplexer.ManualsCategory.Register(manuals, "man", []string{"manuals", "help"}, "An interface to the system reference manuals.")
}

func manuals(context *multiplexer.Context) {
	guildPrefix := context.GenerateGuildPrefix()

	switch {
	case len(context.Fields) == 1:
		{
			// Generate a list of all categories
			embed := embedutil.NewEmbed("Manuals",
				fmt.Sprintf("Issue `%sman <category>` for category-specific manuals.", guildPrefix))
			embed.Color = 0x3492c4

			// The block of text with all categories
			var catText string
			for _, category := range multiplexer.Categories {

				// Only display categories with description set
				if category.Description == "" {
					continue
				}

				catText += fmt.Sprintf("%s - %s \n", category.Title, category.Description)

			}

			// Add list of categories to the Embed
			embed.AddField("Categories", catText, false)
			_ = context.SendEmbed(embed)
		}

	case len(context.Fields) == 2:
		{

			// Figure out if the category exist, and fallthrough if it doesn't
			var desiredCat *multiplexer.CommandCategory
			for _, cat := range multiplexer.Categories {
				if strings.EqualFold(cat.Title, context.Fields[1]) {
					desiredCat = cat
					break
				}
			}

			// Break out of the case if no category was matched
			if desiredCat == nil {
				context.SendMessage(vars.InvalidArgument)
				break
			}

			// Generate list of all commands in one specific category
			embed := embedutil.NewEmbed(desiredCat.Title,
				desiredCat.Description)
			embed.Color = vars.KappaColor

			for _, route := range desiredCat.Routes {

				// Only display routes with proper description
				if route.Description == "" {
					continue
				}

				// Just add the stuff as an entry, it will do on this level
				embed.AddField(route.Pattern, route.Description, false)
			}
			context.SendEmbed(embed)

		}

	case len(context.Fields) > 2:
		{
			// Some catch-all case I guess, though there will be a command-specific thing later maybe
			context.SendMessage(vars.InvalidArgument)
		}
	}
}
