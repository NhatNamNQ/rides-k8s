FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/rides-simulator ./cmd/rides-simulator

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/rides-simulator /usr/local/bin/rides-simulator

USER app

EXPOSE 8081

ENTRYPOINT ["/usr/local/bin/rides-simulator"]
