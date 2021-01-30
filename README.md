<div align="center">

![banner](https://freenitori.jp/img/banner.png "FreeNitori")

An open source, general purpose modular Discord bot written in go.

</div>

---
**Project is still in very early stages and documentation is incomplete/nonexistent in most parts, here will be a way to get started if you want to contribute.**

---

Special Thanks
---
[Rogue](https://twitter.com/RogueDono) (artist)

Building
---
GNU make is required, FreeBSD's make implementation simply does not work. You also need the go compiler.

To build the project, just run `make` in the root of the project, binaries are produced in `build/`. After building for the first time you can skip the dependency part and have make automatically start Nitori by doing `make run`.

Running
---
Nitori produces a configuration file if not present and exits, make sure to edit that file and fill in your stuff before running again.

Discussion
---
We currently have a [Discord guild](https://discord.com/invite/Tap77D3) for discussions on development of this project, however any topic is OK.