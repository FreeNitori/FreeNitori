// +build !windows

package config

import "os"

func firstRun(_ bool) {
	os.Exit(1)
}
