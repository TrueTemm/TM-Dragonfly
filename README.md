<h1 align="center">TM-Dragonfly</h1>

<p align="center">
  A Minecraft: Bedrock Edition server in Go — one listener for every client from 1.21.0 to 1.26.50.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-1.26+-00ADD8?logo=go&logoColor=white" alt="Go 1.26+">
  <img src="https://img.shields.io/badge/bedrock-1.21.0%20→%201.26.50-9BE564" alt="Bedrock 1.21.0 to 1.26.50">
  <img src="https://img.shields.io/badge/license-MIT-brightgreen" alt="MIT License">
  <a href="https://t.me/TMDragonfly"><img src="https://img.shields.io/badge/telegram-@TMDragonfly-229ED9?logo=telegram&logoColor=white" alt="Telegram"></a>
</p>

## What is this?

TM-Dragonfly is a fork of [Dragonfly](https://github.com/df-mc/dragonfly) that adds full
**multiversion** support: players from Bedrock **1.21.0 right up to 1.26.50** all connect to the
same server, on the same port, with no proxy in front of it. Blocks, items and packets are
translated per client; the newest version runs natively with no translation at all.

Under the hood it's still Dragonfly — fast, heavily concurrent, and meant to be used as a Go
library you build your own game on top of.

## Features

- **Multiversion** — Bedrock 1.21.0 … 1.26.50 on a single listener, no external proxy.
- **Native latest version** — newest client runs with zero translation overhead.
- **Clean console** — light-green themed logs and a handful of built-in commands.
- **Full Dragonfly API** — worlds, entities, inventories, blocks and items, all unchanged.
- **No forks** — plain, official gophertunnel and go-raknet; nothing vendored.

## Getting started

Needs **Go 1.26+**.

```shell
git clone https://github.com/TrueTemm/TM-Dragonfly
cd TM-Dragonfly
go run .
```

A `config.toml` is written next to the binary on first start — edit it and restart. Type `stop`
in the console or press **ctrl+c** to shut down cleanly.

Release build:

```shell
go build -o tm-dragonfly .
```

## Console commands

| Command         | What it does                                      |
|-----------------|---------------------------------------------------|
| `help`          | list the commands                                 |
| `list`          | online players and the client version each is on  |
| `status`        | uptime, TPS, load and memory                      |
| `version`       | build info and the supported version range        |
| `say <message>` | broadcast a message to everyone                   |
| `stop`          | shut the server down                              |

## Developer info

The `server` package is the entry point; the subpackages carry the block, item and world APIs.
It's a drop-in Dragonfly fork, so the whole upstream API works as-is.

## License

MIT — see [LICENSE](LICENSE). TM-Dragonfly is built on
[Dragonfly](https://github.com/df-mc/dragonfly) by Dragonfly Tech.

Telegram: [@TMDragonfly](https://t.me/TMDragonfly)
