# binary names
BINARY_NAME=github-actions-badge

# go commands
GOCMD=go
GOBUILD=$(GOCMD) build -v -o $(BINARY_NAME)

docker-binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "-extldflags '-static'"
