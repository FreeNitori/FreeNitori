<div align="center">

![banner](https://git.randomchars.net/RandomChars/FreeNitori/-/raw/master/assets/web/static/banner.png "FreeNitori")

An open source, general purpose modular Discord bot written in go.

</div>

---
**Project is still in very early stages and documentation is incomplete/nonexistent in most parts, here will be a way to get started if you want to contribute.**

---

Building
---
You need GNU Make, a POSIX-compliant shell and of course, go. The latest version from your package manager should work as long as you aren't using something ancient. 

To build the project, just run `make` in the root of the project, binaries are produced in `build/`. After building for the first time you can skip the dependency part and have make automatically start Nitori by doing `make run`.

Running
---
Nitori produces a configuration file if not present and exits, make sure to edit that file and fill in your stuff before running again.

Discussion
---
We currently have a [Discord guild](https://discord.com/invite/Tap77D3) for discussions on development of this project, however any topic is OK.