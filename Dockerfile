FROM golang:1.23-alpine as builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o thumbnail cmd/main.go

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/thumbnail /usr/local/bin/thumbnail
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/internal /app/internal

RUN chmod 644 /app/.env

EXPOSE 8080

CMD ["/usr/local/bin/thumbnail"]