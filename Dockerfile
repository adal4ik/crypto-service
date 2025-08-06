FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/app

FROM alpine:latest

WORKDIR /

COPY --from=builder /main /main


CMD ["/main"]