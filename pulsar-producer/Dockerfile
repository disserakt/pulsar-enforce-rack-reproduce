FROM golang:1.16 AS builder

COPY ./ /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o producer ./


FROM debian
COPY --from=builder /app/producer /app/producer
RUN ls -al /app/producer

CMD ["/app/producer"]