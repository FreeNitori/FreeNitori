package multiplexer

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"math/rand"
	"strconv"
)

func ProcessMessageExperience(context *Context) {
	var err error

	// Not do anything if private or bot
	if context.IsPrivate {
		return
	}

	// Also don't do anything if experience system is disabled
	expEnabled, err := config.ExpEnabled(context.Guild.ID)
	if err != nil {
		return
	}
	if !expEnabled {
		return
	}

	previousExp, err := config.GetMemberExp(context.Member)
	if err != nil {
		Logger.Error(fmt.Sprintf("Database error on user experience advancing, %s", err))
		return
	}
	advancedExp := previousExp + rand.Intn(10) + 5
	err = config.SetMemberExp(context.Member, advancedExp)
	if err != nil {
		Logger.Error(fmt.Sprintf("Database error on user experience advancing, %s", err))
		return
	}
	advancedLevel := config.ExpToLevel(advancedExp)
	if advancedLevel > config.ExpToLevel(previousExp) {
		context.SendMessage(fmt.Sprintf("Level up message, %s", strconv.Itoa(advancedLevel)), fmt.Sprintf("generating level up message for %s", context.Guild.ID))
	}
}
