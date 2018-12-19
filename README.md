# steamutil

## What is it?
A Go library for working with [Steam](https://steampowered.com) - Valve
Software's game distribution application.

## API
The following subsections will provide a high level overview of the APIs
provided by this library. Examples will be provided for the more notable
use cases in the [examples directory](examples/).

#### locations
Package locations provides functionality for locating notable Steam application
files and directories.

- [Checking if Steam is installed](examples/is-steam-installed/main.go)
- [Finding where Steam stores its data](examples/steam-data/main.go)

#### shortcuts
Package shortcuts provides functionality for working with Steam's custom
game shortcuts.

- [Read a shortcuts file](examples/read-shortcuts-file/main.go)
- [Create or update a shortcut file](examples/update-shortcuts-file/main.go)

#### vdf
Package vdf provides functionality for working with Steam's .vdf file format.
