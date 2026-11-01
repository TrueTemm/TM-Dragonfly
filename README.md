<h1 align="center">TM-Dragonfly</h1>

<p align="center">
  Minecraft: Bedrock server software in Go — one listener for every client from 1.20.60 to 1.26.60.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-1.26+-00ADD8?logo=go&logoColor=white" alt="Go 1.26+">
  <img src="https://img.shields.io/badge/bedrock-1.20.60%20→%201.26.60-9BE564" alt="Bedrock 1.20.60 to 1.26.60">
  <img src="https://img.shields.io/badge/license-MIT-brightgreen" alt="MIT License">
  <a href="https://t.me/TMDragonfly"><img src="https://img.shields.io/badge/telegram-@TMDragonfly-229ED9?logo=telegram&logoColor=white" alt="Telegram"></a>
</p>

## What is this?

A fork of [Dragonfly](https://github.com/df-mc/dragonfly) with full **multiversion** support:
every client from Bedrock **1.20.60 to 1.26.60** joins the same server, same port, no proxy. The
newest version runs natively; older ones are translated per client. Still the full Dragonfly API
underneath — build your game on top of it as a Go library.

## Features

- **Multiversion** — 1.20.60 … 1.26.60 on the same server and game port.
- **RakNet and NetherNet** — UDP gameplay with HTTP/HTTPS signaling for newer clients.
- **Config** — world mode (`current` or `flat`), render distance, difficulty, spawn protection.
- **Operators** — kept in `ops.yml`, with in-game commands: `weather`, `time set`, `gamemode`, `tp`, `give` (tab-complete items), `about`.
- **Protected packs** — drop a pack in `resources/` and it is served encrypted; the client's cached copy is unreadable.

## Getting started

Needs **Go 1.26+**.

```shell
git clone https://github.com/TrueTemm/TM-Dragonfly
cd TM-Dragonfly
go run .
```

A `config.toml` and `ops.yml` are written on first start. Type `stop` or press **ctrl+c** to shut
down. Release build: `go build -o tm-dragonfly .`

## Network ports

Clients use the same game address and port for every supported version. The default is `19132`.

| Protocol | Port | Purpose |
| --- | --- | --- |
| UDP | `19132` | RakNet and NetherNet game traffic |
| TCP | `19132` | NetherNet HTTP/HTTPS signaling |

Open both ports on the host firewall and any provider firewall:

```shell
sudo ufw allow 19132/udp
sudo ufw allow 19132/tcp
```

NetherNet is enabled by default with a shared UDP socket. If `SharedUDP` is disabled in
`config.toml`, its ICE UDP port defaults to `19133` and must also be opened. Set the NetherNet
certificate and private key paths in `config.toml` to enable HTTPS signaling. Port hints are
printed only when a NetherNet listener fails to open or stops with an error.

## Console commands

`help`, `list`, `status`, `version`, `about`, `weather`, `time set`, `op <nick>`, `deop <nick>`,
`say <message>`, `stop`. Operators also get the in-game commands listed above.

## License

MIT — see [LICENSE](LICENSE). Built on [Dragonfly](https://github.com/df-mc/dragonfly) by Dragonfly Tech.

Telegram: [@TMDragonfly](https://t.me/TMDragonfly)
