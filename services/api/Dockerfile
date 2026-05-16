FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/rides-api ./cmd/rides-api

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/rides-api /usr/local/bin/rides-api

USER app

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/rides-api"]

