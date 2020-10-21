# FreeNitori
[Download latest build](https://git.randomchars.net/RandomChars/FreeNitori/-/jobs/artifacts/master/raw/build/freenitori.tar.gz?job=build)

FreeNitori is a general purpose Discord bot written in Golang.

---
**Project is still in very early stages and documentation is incomplete/nonexistent in most parts, here will be a way to get started if you want to contribute.**

Download that archive from the URL above, extract it somewhere, run `freenitori-supervisor` or `freenitori` (they are both the supervisor program) once, and edit the configuration file, fill in the credentials and replace the default binary paths to the locations you placed FreeNitori binaries, then run the supervisor program again.

When using the Makefile, remember to run `make` at least once then use `make test` in development (which retains the extra debug information, remove an extra downloading stage and automatically starts the program).

---