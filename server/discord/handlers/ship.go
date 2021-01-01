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
	user1 := context.GetMember(context.Fields[1]).User
	user2 := context.GetMember(context.Fields[2]).User
	res := int(math.Mod(float64(func() int { id, _ := strconv.Atoi(user1.ID); return id }()^func() int { id, _ := strconv.Atoi(user2.ID); return id }()), 101))
	embed := embedutil.NewEmbed(
		fmt.Sprintf("`%s` ❤️ `%s`", user1.Username, user2.Username),
		fmt.Sprintf("[%s%s] %v", strings.Repeat("=", res/2), strings.Repeat("-", 50-res/2), res)+"%")
	context.SendEmbed(embed)
}
