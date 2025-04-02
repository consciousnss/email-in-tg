FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /email ./cmd

FROM alpine:3.21 AS image
COPY --from=builder /email /email
CMD ["/email"]