package args

import (
	"flag"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"go/types"
)

var _ = flags()

func flags() *types.Nil {
	flag.StringVar(&config.TokenOverride, "a", "", "Override Discord Authorization Token")
	flag.StringVar(&config.NitoriConfPath, "c", "", "Specify configuration file path")
	flag.Parse()
	return nil
}
