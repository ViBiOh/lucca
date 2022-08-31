# remote-lucca

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
