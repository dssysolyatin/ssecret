GOCMD						= go
GOBUILD						= $(GOCMD) build
SSECRET_CLI_CMD_SRC			= cmd/ssecret/ssecret.go
SSECRET_CLI_CMD_BINARY		= ssecret

fmt:
	find . -name *.go -not -path "./vendor/*" -exec goimports -w {} +

.PHONY: fmt