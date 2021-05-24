package main

import (
	log "git.randomchars.net/FreeNitori/Log"
	"io/ioutil"
	"time"
)

type out struct{}

func (out) Write(p []byte) (n int, err error) {
	return len(p), ioutil.WriteFile("/dev/tty1", p, 0644)
}

func main() {
	log.Instance.SetOutput(out{})
	log.Info("Nitori!")
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
