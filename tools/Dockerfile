FROM golang:1.17

WORKDIR /go/src/app

RUN go install github.com/volatiletech/sqlboiler/v4@v4.8.3 && \
    go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.8.3


CMD sqlboiler mysql
