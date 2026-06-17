# ListenBrainz RPC

A CLI tool that shows your ListenBrainz "now playing" track as a Discord Rich Presence activity.

Polls the ListenBrainz API every 10 seconds, resolves album art via the Cover Art Archive, and updates your Discord status accordingly.

## Features

- Shows track name, artist, and album as a Discord activity
- Displays album art from the Cover Art Archive
- Shows track progress with start/end timestamps
- Links to your ListenBrainz profile and MusicBrainz release page
- TUI with spinner showing current playback state
- Configurable via CLI flags, environment variables, or YAML config file

## Installation

### From source

```bash
go install github.com/ikanexus/listenbrainz-rpc@latest
```

### From release

Download the latest binary for your platform from the [releases page](https://github.com/ikanexus/listenbrainz-rpc/releases).

## Usage

```bash
listenbrainz-rpc --user <listenbrainz-username>
```

Press any key to quit.

## Configuration

All flags can be set via CLI flags, environment variables, or a YAML config file.

| Flag | Short | Default | Env Var | Description |
|------|-------|---------|---------|-------------|
| `--config` | | `$XDG_CONFIG_HOME/listenbrainz-rpc.yaml` | `LISTENBRAINZ_CONFIG` | Config file path |
| `--app-id` | `-a` | `1232457767726485545` | `LISTENBRAINZ_APP_ID` | Discord Application ID |
| `--user` | `-u` | *(required)* | `LISTENBRAINZ_USER` | ListenBrainz username |
| `--verbose` | `-v` | `false` | `LISTENBRAINZ_VERBOSE` | Enable debug logging |

### YAML config file

Create a config file at `$XDG_CONFIG_HOME/listenbrainz-rpc.yaml` (or specify a path with `--config`):

```yaml
user: your-listenbrainz-username
app-id: "1232457767726485545"
verbose: false
```

## Building

```bash
go build -o listenbrainz-rpc .
```

For a local snapshot release across all platforms:

```bash
goreleaser release --clean --snapshot
```

## How it works

```
Poll (every 10s) → ListenBrainz API: "now playing?"
  ├─ New track    → Login to Discord IPC, fetch album art, set activity
  ├─ Same track   → Continue polling
  └─ No track     → Logout from Discord IPC, wait for next poll
```

The tool uses the [ListenBrainz API](https://listenbrainz.readthedocs.io/) to check what's currently playing, resolves MusicBrainz IDs for album art via the [Cover Art Archive](https://coverartarchive.org/), and updates your Discord Rich Presence through the local IPC socket.

## Requirements

- Discord desktop client running locally
- A ListenBrainz account with scrobbling enabled
