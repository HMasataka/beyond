FROM golang:1.24 AS builder

RUN go install github.com/volatiletech/sqlboiler/v4@v4.18.0 && \
    go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.18.0

FROM gcr.io/distroless/base-debian12

COPY --from=builder /go/bin/sqlboiler /bin/
COPY --from=builder /go/bin/sqlboiler-mysql /bin/

ENTRYPOINT ["sqlboiler", "--no-tests", "--add-soft-deletes"]
