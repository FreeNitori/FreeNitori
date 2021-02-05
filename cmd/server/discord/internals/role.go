package internals

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/overrides"
	"github.com/bwmarrin/discordgo"
)

func init() {
	multiplexer.GuildMemberAdd = append(multiplexer.GuildMemberAdd, autoRoleHandler)
	overrides.RegisterComplexEntry(overrides.ComplexConfigurationEntry{
		Name:         "role",
		FriendlyName: "Role Assignment",
		Description:  "Configure role assignment related utilities.",
		Entries: []overrides.SimpleConfigurationEntry{
			{
				Name:         "join",
				FriendlyName: "Automatic Role Assignment",
				Description:  "Role automatically assigned on join.",
				DatabaseKey:  "role_join",
				Cleanup:      func(context *multiplexer.Context) {},
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					if role := context.GetRole(*input); role != nil {
						*input = role.ID
						return true, true
					} else {
						return false, true
					}
				},
				Format: func(context *multiplexer.Context, value string) (string, string, bool) {
					if role := context.GetRole(value); role != nil {
						return role.Name, role.ID, true
					}
					return "No role configured", fmt.Sprintf("Configure it by issuing command `%sconf role join <role>`.", context.Prefix()), true
				},
			},
		},
		CustomEntries: nil,
	})
}

func autoRoleHandler(session *discordgo.Session, add *discordgo.GuildMemberAdd) {
	role := false
	roleID, err := config.GetGuildConfValue(add.GuildID, "role_join")
	if err != nil {
		return
	}
	if roleID == "" {
		return
	}
	guild, err := session.State.Guild(add.GuildID)
	if err != nil {
		guild, err = session.Guild(add.GuildID)
		if err != nil {
			return
		}
		_ = session.State.GuildAdd(guild)
	}

	for _, r := range guild.Roles {
		if r.ID == roleID {
			role = true
		}
	}

	if !role {
		return
	}
	_ = session.GuildMemberRoleAdd(add.GuildID, add.User.ID, roleID)
}
