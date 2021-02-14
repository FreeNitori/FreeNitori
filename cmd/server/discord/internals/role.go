package internals

import (
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/overrides"
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
					}
					return false, true
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

	guild, err := session.State.Guild(add.GuildID)
	if err != nil {
		guild, err = session.Guild(add.GuildID)
		if err != nil {
			return
		}
		_ = session.State.GuildAdd(guild)
	}

	if len(guild.Channels) == 0 {
		return
	}

	// If Nitori has permission
	permissions, err := session.State.UserChannelPermissions(session.State.User.ID, guild.Channels[0].ID)
	if !(err == nil && (permissions&discordgo.PermissionManageRoles == discordgo.PermissionManageRoles)) {
		return
	}

	role := false
	roleID, err := config.GetGuildConfValue(add.GuildID, "role_join")
	if err != nil {
		return
	}
	if roleID == "" {
		return
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
