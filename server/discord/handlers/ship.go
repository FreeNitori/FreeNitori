package handlers

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"math"
	"strconv"
	"strings"
)

func init() {
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "ship",
		AliasPatterns: []string{},
		Description:   "Shipping two users via mathematics.",
		Category:      multiplexer.MediaCategory,
		Handler:       ship,
	})
}

func ship(context *multiplexer.Context) {
	if context.IsPrivate {
		context.SendMessage(vars.GuildOnly)
		return
	}
	if len(context.Fields) != 3 {
		context.SendMessage(vars.InvalidArgument)
		return
	}
	member1 := context.GetMember(context.Fields[1])
	member2 := context.GetMember(context.Fields[2])
	if member1 == nil || member2 == nil {
		context.SendMessage(vars.MissingUser)
		return
	}
	res := int(math.Mod(float64(func() int { id, _ := strconv.Atoi(member1.User.ID); return id }()^func() int { id, _ := strconv.Atoi(member2.User.ID); return id }()), 101))
	embed := embedutil.NewEmbed(
		fmt.Sprintf("`%s` ❤️ `%s`", member1.User.Username, member2.User.Username),
		fmt.Sprintf("[%s%s] %v", strings.Repeat("=", res/2), strings.Repeat("-", 50-res/2), res)+"%")
	context.SendEmbed(embed)
}
