FROM golang:1.24-alpine AS builder
WORKDIR /src

# ensure certificates and git for fetching modules
RUN apk add --no-cache git ca-certificates

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the server binary from the cmd package
RUN go build -v -o /app/server ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /server

EXPOSE 3000
ENV PORT=3000

ENTRYPOINT ["/server"]
