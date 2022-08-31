# remote-lucca

## Getting started

### Release

Download the latest binary for your os and architecture from the [GitHub Releases page](https://github.com/ViBiOh/remote-lucca/releases)

```bash
curl \
  --disable \
  --silent \
  --show-error \
  --location \
  --max-time 300 \
  --output "/usr/local/bin/remote-lucca"
  https://github.com/ViBiOh/remote-lucca/releases/download/v0.0.1/remote-lucca_$(uname -s | tr "[:upper:]" "[:lower:]")_amd64
chmod +x "/usr/local/bin/remote-lucca"
```

### Golang

```bash
go install "github.com/ViBiOh/remote-lucca@latest"
```

## Usage

```bash
Usage of remote-lucca:
  -days string
        Days of week, comma separated
  -dry-run
        Dry run
  -end string
        End of repetition, in ISO format
  -leaveType string
        Type of leave request (default "Télétravail")
  -password string
        Password
  -start string
        Start of repetition, in ISO format
  -subdomain string
        Sub domain used
  -username string
        Username
```

### Example

```bash
go run main.go \
  -subdomain company \
  -days Monday,Friday \
  -start 2022-09-01 \
  -end 2022-10-01 \
  -username username@company.com \
  -password "USE_A_PASSWORD_MANAGER_PLEASE" \
  -dry-run
```
