FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/ydsterm-server .

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/ydsterm-server .
COPY manifest/config/config.docker.yaml ./manifest/config/config.yaml

RUN mkdir -p /app/data

EXPOSE 7788

VOLUME ["/app/data"]

ENV YDSTERM_SERVER_KEY=""

CMD ["./ydsterm-server"]
