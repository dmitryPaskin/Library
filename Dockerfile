FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main cmd/main.go

FROM alpine:latest

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]