// +build !windows

package config

import "os"

func first(_ bool) {
	os.Exit(1)
}
