FROM golang:1.23-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o thumbnail cmd/main.go

FROM golang:1.23-alpine as client-builder

WORKDIR /app

COPY . .


RUN go mod tidy
RUN go build -o client internal/proto/client/client.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/thumbnail /usr/local/bin/thumbnail
COPY --from=client-builder /app/client /usr/local/bin/client
COPY --from=builder /app/.env /app/.env

RUN chmod 644 /app/.env


EXPOSE 8080


CMD ["/usr/local/bin/thumbnail"]
