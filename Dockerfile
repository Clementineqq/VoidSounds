FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@v0.3.1001 && templ generate
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/voidsounds ./cmd/server

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/voidsounds /app/voidsounds

COPY --from=builder /build/migrations /app/migrations

EXPOSE 8080

CMD ["/app/voidsounds"]