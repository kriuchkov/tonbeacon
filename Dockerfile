FROM golang:1.24-alpine as gomod
WORKDIR /gomod
ADD go.mod go.sum ./
RUN go mod download

FROM golang:1.24-alpine as builder
WORKDIR /builder
COPY --from=gomod /go/pkg /go/pkg
COPY . .

RUN GOOS=linux go build -a -o /app/.build/ ./cmd/*

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/.build/ /app/