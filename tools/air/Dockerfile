FROM golang:1.24-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@v1.61.1

CMD ["air", "-c", ".air.toml"]
