# gr8 - a Chip-8 emulator

[![Unit tests](https://github.com/aricodes-oss/gr8/actions/workflows/tests.yml/badge.svg)](https://github.com/aricodes-oss/gr8/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/aricodes-oss/gr8)](https://goreportcard.com/report/github.com/aricodes-oss/gr8)

A small [Chip-8 emulator](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM) written in Go.

![Preview](https://raw.githubusercontent.com/aricodes-oss/gr8/main/gr8.gif)

## Installing

### via Go

```sh
go install github.com/aricodes-oss/gr8@latest
```

### from source

First, [clone the repository](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository) and `cd` into it. Then,

```sh
go install .
```

Alternatively, you can use `go run` to run directly from source:

```sh
go run main.go /path/to/rom.ch8
```

## Running

Simply invoke the command with the path to a valid ROM file:

```sh
gr8 /path/to/rom.ch8
```

For convenience I have included the 8 roms from [Timendus' chip8-test-suite](https://github.com/Timendus/chip8-test-suite) in the `roms/` directory in this repository.

For more options, see the usage page:

```
Usage:
  gr8 [flags]

Flags:
  -h, --help        help for gr8
  -s, --scale int   screen scaling factor (default 16)
```
