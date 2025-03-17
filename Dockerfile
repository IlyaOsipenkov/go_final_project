FROM golang:1.24.1 AS builder

WORKDIR /app

COPY . .

RUN go build -o scheduler .

FROM ubuntu:latest
WORKDIR /root/

RUN mkdir -p /root/pkg/db

COPY --from=builder /app/scheduler .
COPY --from=builder /app/web ./web

EXPOSE 7540

CMD ["./scheduler"]
