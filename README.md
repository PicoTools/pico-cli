# Pico CLI

Pico CLI is a operator CLI for Pico C2

## Quick Start

1. Download latest [release](https://github.com/PicoTools/pico-cli/releases) or build it yourself (see [Makefile](https://github.com/PicoTools/pico-cli/blob/master/Makefile)).
2. Connect to the Pico C2 server

```sh
$ ./pico-cli -H <server_ip>:<port> -t <operator_token>
pico > whoami
pico-operator
```

## Features

- [x] Command completion and sugestion
- [x] [PLAN](https://github.com/PicoTools/plan) scripts support
