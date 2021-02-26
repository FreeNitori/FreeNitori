package internals

import (
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/overrides"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"github.com/bwmarrin/discordgo"
)

func init() {
	state.Multiplexer.GuildMemberAdd = append(state.Multiplexer.GuildMemberAdd, autoRoleHandler)
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

func autoRoleHandler(context *multiplexer.Context) {

	if len(context.Guild.Channels) == 0 {
		return
	}

	// If Nitori has permission
	permissions, err := context.Session.State.UserChannelPermissions(context.Session.State.User.ID, context.Guild.Channels[0].ID)
	if !(err == nil && (permissions&discordgo.PermissionManageRoles == discordgo.PermissionManageRoles)) {
		return
	}

	role := false
	roleID, err := config.GetGuildConfValue(context.Guild.ID, "role_join")
	if err != nil {
		return
	}
	if roleID == "" {
		return
	}

	for _, r := range context.Guild.Roles {
		if r.ID == roleID {
			role = true
		}
	}

	if !role {
		return
	}
	_ = context.Session.GuildMemberRoleAdd(context.Guild.ID, context.User.ID, roleID)
}
