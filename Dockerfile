FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build


FROM alpine
WORKDIR /app
COPY --from=builder /app/jams-exporter /app/jams-exporter
ENTRYPOINT ["/app/jams-exporter"]
