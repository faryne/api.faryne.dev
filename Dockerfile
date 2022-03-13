FROM golang:1.17.6-alpine3.15 AS builder

RUN mkdir /apiv2 && apk add git
COPY . /apiv2
WORKDIR /apiv2

RUN apk add git && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init && \
    go build -o server main.go

FROM alpine:3.15
WORKDIR /
COPY --from=builder /apiv2/server .
CMD ["/server"]