# XGo

[![Yum](https://img.shields.io/badge/-Buy%20me%20a%20cookie-blue?labelColor=grey&logo=cookiecutter&style=for-the-badge)](https://www.buymeacoffee.com/mjwhitta)

[![Go Report Card](https://goreportcard.com/badge/github.com/mjwhitta/xgo?style=for-the-badge)](https://goreportcard.com/report/github.com/mjwhitta/xgo)
![License](https://img.shields.io/github/license/mjwhitta/xgo?style=for-the-badge)

## What is this?

This Go module is a drop-in replacement for `go` as it is a simple
wrapper tool. It will parse your `GOARCH` and `GOOS` environment vars
to determine which cross-compilers to use, if needed. It will then set
the appropriate `CC` and `CXX` environment vars before calling the
requested `go` subcommand. If the subcommand is `build`, `get` or
`install` it will ensure some sane default flags are provided:

- `--buildvcs=false`
- `--ldflags="-s -w"`
- `--trimpath`

If you specify these flags manually, then they will not be added or
modified.

## How to install

Open a terminal and run the following:

```
$ go install github.com/mjwhitta/xgo/cmd/xgo@latest
```

## Usage

You can simply use `xgo` in place of `go` or create a shell alias:

### Bash/Zsh

```
$ xgo build .
$ # OR
$ alias go="xgo" # Add this to your shell rc file
$ go build .
```

### PowerShell

```
PS> xgo build .
PS> # OR
PS> set-alias -option allscope go xgo # Add this to your profile.ps1
PS> go build .
```

Additionally, if you do not want to set environment vars in PowerShell
(because it's not as convenient as bash/zsh), there are two CLI
options that could be useful to you: `--goarch` and `--goos`.

### Scripts

There is a hidden `-d`/`--debug` CLI option that can be used to
display the underlying shell command. The output is bash/zsh
compatible on Darwin and Linux, and PowerShell on Windows. I never
really meant to support that functionality, but I used it for
debugging and figured somebody might find it useful someday.

```
$ GOOS=windows xgo -d build .
```

## Cross-Compilers per host OS

### Darwin hosts

To compile for Linux (386 and amd64):

```
$ brew install musl-cross --with-i468 --without-aarch64
```

To compile for Windows (386 and amd64):

```
$ brew install mingw-w64
```

For Zig support (really *NOT* recommended, doesn't fully work!):

```
$ brew install zig
$ export XGOZIG=1
```

### Linux hosts

To compile for Darwin (amd64 and arm64) you will need to install
[osxcross]. This requires grabbing your own copy of the macOS SDK. I
have found [this script][gen_sdk_package.sh] on a macOS host with
XCode Commandline Tools to be the easiest method. You will need to add
`XCODE_TOOLS=1` to line 2 of that script: Additionally, that repo
appears broken as of 2024-07-05, but you can `git checkout 2389e32b`
to rollback to a working version.

To compile for Linux (386) you will need to install `gcc` with
multilib support.

To compile for Windows (386 and amd64) you will need to install
[MinGW-w64].

For Zig support (really *NOT* recommended, doesn't fully work!):

```
$ # Choose a command below for your package manager:
$ # apk add zig
$ # apt install zig
$ # apt-get install zig
$ # dnf install zig
$ # pacman -S zig
$ # zypper in zig
$ export XGOZIG=1
```

### Windows hosts

To compile for Darwin (amd64 and arm64) you will need to install TODO.

To compile for Linux (386 and amd64) you will need to install TODO.

To compile for Windows (386) you will need to install TODO.

For Zig support (really *NOT* recommended, leaves all sorts of
artifacts in your binaries):

```
PS> choco install zig
PS> $env:XGOZIG="1"
```

**NOTE**: Zig will allow for compiling of 386/amd64 Linux and Windows.
I have not managed to get it to work for Darwin at this time.

## Links

- [Source](https://github.com/mjwhitta/xgo)

## TODO

- Work on Windows (host) support

[gen_sdk_package.sh]: https://github.com/tpoechtrager/osxcross/blob/master/tools/gen_sdk_package.sh
[MinGW-w64]: https://www.mingw-w64.org
[osxcross]: https://github.com/tpoechtrager/osxcross
