FROM golang:1.18-alpine3.16 as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -tags netgo -o main.app ./app

FROM alpine:latest

COPY --from=builder /app/main.app .

CMD ["./main.app"]