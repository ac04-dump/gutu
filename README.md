
# gutu

[![Release](https://img.shields.io/github/v/release/alexcoder04/gutu)](https://github.com/alexcoder04/gutu/releases/latest)
[![AUR](https://img.shields.io/aur/version/gutu)](https://aur.archlinux.org/packages/gutu)
[![Top language](https://img.shields.io/github/languages/top/alexcoder04/gutu)](https://github.com/alexcoder04/gutu/search?l=go)
[![License](https://img.shields.io/github/license/alexcoder04/gutu)](https://github.com/alexcoder04/gutu/blob/main/LICENSE)
[![Issues](https://img.shields.io/github/issues/alexcoder04/gutu)](https://github.com/alexcoder04/gutu/issues)
[![Pull requests](https://img.shields.io/github/issues-pr/alexcoder04/gutu)](https://github.com/alexcoder04/gutu/pulls)

`gutu` manages your desktop startup applications. Especially on standalone
window managers, it can be used to uniform the processes you need to run (e. g.
notification daemon, keybind handler etc)

## Installation

### Arch Linux (AUR)

gutu is available on the [AUR](https://aur.archlinux.org/packages/gutu)

### Prebuild release (all distros, binary only)

Download the binary from [github.com/alexcoder04/gutu/releases/latest](https://github.com/alexcoder04/gutu/releases/latest)
and copy it to a directory in your `$PATH`.

### Other Distros (Make)

```sh
git clone "https://github.com/alexcoder04/gutu.git"
cd gutu

go build .   # build the binary
go install . # install the executable to your $GOPATH
```

## Usage

Auto-start gutu in your window manager / desktop environment startup: execute
`gutu`. For configuring the services, see below.

## Configuration

Every service file goes into a separate file in `$XDG_CONFIG_HOME/gutu`. These
files can have `.service`, `.yml` and `.yaml` extensions.

### Example service

```yaml
# your name for the service
Name: compositor
Command: picom
Args: ["--experimental-backends"]
# restart if fails, at most 3 times
KeepAlive: true
RetryNumber: 3
# start only on Xorg
When: x11
# kill picom instances that are still running
KillOld: true
```

For more examples, see [`contrib`](https://github.com/alexcoder04/gutu/tree/main/contrib).

### Configuration fields

```yaml
Name        # Name of service
Command     # Command to run
Args        # Arguments for the command (["-c", "arg1", "arg2"])
Interval    # Re-run command periodically in this interval (in seconds, 0=never)
KeepAlive   # Restart the command if it exits (true/false)
RetryNumber # How often try to restart the command (if KeepAlive=true)
When        # wayland/x11/always/never
Delay       # Number of seconds to wait before starting
KillOld     # Kill running "Command" processes
```

