# Health Monitoring

## Build

```shell
go build -ldflags "-X main.version=v0.1.5" -o .hm main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=v0.1.5" -o .hm main.go
```
