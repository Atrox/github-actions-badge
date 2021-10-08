# build binary
FROM golang:1.17.2 AS builder
WORKDIR /app

# populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
RUN go mod download

# build
COPY . .
RUN make docker-binary


# run image
FROM alpine:3.13.5

# ca-certificates
RUN apk add --no-cache ca-certificates

# add binary
COPY --from=builder /app/github-actions-badge /

# ports
EXPOSE 3000

# run binary
CMD ["/github-actions-badge"]
