# gh-zen 🧘

A zen garden.

Rake sand, navigate around boulders, and find your calm. When the garden is fully raked, receive wisdom.

## Install

```sh
gh extension install bagtoad/gh-zen
```

## Usage

```sh
gh zen
```

## Controls

| Key | Action |
|---|---|
| `↑` `↓` `←` `→` / `h` `j` `k` `l` | Move the rake |
| `r` | New garden |
| `q` / `Ctrl+C` | Quit |

## How It Works

- The garden fills your terminal with **sand** (`░`)
- Move your **rake** to carve directional trails (`─` `│`)
- Sand ahead of the rake gets **pushed** forward
- Sand pushed against a wall or **rock** is **flattened** (`·`)
- Multi-character **boulders** block your path
- Rake all the sand to receive a **zen quote** from `api.github.com/zen`

## Building from Source

```sh
git clone https://github.com/bagtoad/gh-zen.git
cd gh-zen
go build -o gh-zen .
./gh-zen
```
