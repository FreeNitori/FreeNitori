package config

import (
	log "git.randomchars.net/FreeNitori/Log"
	"strconv"
	"time"
)

func CreationTime(snowflake string) time.Time {
	id, err := strconv.Atoi(snowflake)
	if err != nil {
		log.Debugf("Unexpected snowflake passed to CreationTime.")
		return time.Unix(0, 0)
	}
	return time.Unix(int64(((id>>22)+1420070400000)/1000), 0)
}

func ValidateSnowflake(snowflake string) bool {
	id, err := strconv.Atoi(snowflake)
	if err != nil {
		return false
	}
	return id>>22 >= 0
}
