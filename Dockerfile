FROM golang:1.25.3-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -u 10001 appuser

WORKDIR /app
COPY --from=builder /app/server /app/server

USER appuser
EXPOSE 3030

ENTRYPOINT ["/app/server"]
