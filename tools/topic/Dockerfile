FROM golang:1.24 AS builder

COPY ./ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/topic tools/topic/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /go/bin/topic topic

ENTRYPOINT ["/topic"]
