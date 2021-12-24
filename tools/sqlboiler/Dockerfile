FROM golang:1.17

RUN go install github.com/volatiletech/sqlboiler/v4@v4.8.3 && \
    go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.8.3

ENTRYPOINT ["sqlboiler"]
