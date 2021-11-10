FROM golang:1-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o /build/faktory-cron

FROM alpine
COPY --from=builder /build/faktory-cron /bin/faktory-cron
CMD ["/bin/faktory-cron"]
