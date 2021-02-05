# build binary
FROM golang:1.15.8 AS builder
WORKDIR /app

# populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
RUN go mod download

# build
COPY . .
RUN make docker-binary


# run image
FROM alpine:3.12.0

# ca-certificates
RUN apk add --no-cache ca-certificates

# add binary
COPY --from=builder /app/github-actions-badge /

# ports
EXPOSE 3000

# run binary
CMD ["/github-actions-badge"]
