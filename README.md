# lucca

[![Build](https://github.com/ViBiOh/lucca/workflows/Build/badge.svg)](https://github.com/ViBiOh/lucca/actions)
[![codecov](https://codecov.io/gh/ViBiOh/lucca/branch/main/graph/badge.svg)](https://codecov.io/gh/ViBiOh/lucca)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_lucca&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_lucca)

## Getting started

### Release

Download the latest binary for your os and architecture from the [GitHub Releases page](https://github.com/ViBiOh/lucca/releases)

```bash
curl \
  --disable \
  --silent \
  --show-error \
  --location \
  --max-time 300 \
  --output "/usr/local/bin/lucca"
  https://github.com/ViBiOh/lucca/releases/download/v0.1.0/lucca_$(uname -s | tr "[:upper:]" "[:lower:]")_amd64
chmod +x "/usr/local/bin/lucca"
```

### Golang

```bash
go install "github.com/ViBiOh/lucca@latest"
```

## Usage

```bash
Run Lucca action fro the CLI

Usage:
  lucca [flags]
  lucca [command]

Available Commands:
  birthdays   Birthdays of the day
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  leave       Create a leave request

Flags:
      --dry-run            Dry run
  -h, --help               help for lucca
      --password string    Password
      --subdomain string   Subdomain
      --username string    Username

Use "lucca [command] --help" for more information about a command.
```

### Example

> Creating recurring remote work on Monday and Friday in September.

```bash
go run main.go \
  leave \
  --subdomain company \
  --username username@company.com \
  --password "USE_A_PASSWORD_MANAGER_PLEASE" \
  --days Monday \
  --days Friday \
  --start 2022-09-01 \
  --end 2022-10-01 \
  --dry-run
```

> Get birthdays of the day

```bash
go run main.go \
  birthdays \
  --subdomain company \
  --username username@company.com \
  --password "USE_A_PASSWORD_MANAGER_PLEASE"
```
