FROM golang:1.22 AS builder

WORKDIR /auth

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 80

RUN go build -o /go/bin/app ./cmd/main.go

FROM alpine:3.20

WORKDIR /auth

COPY --from=builder /auth/internal/migration /migration

COPY --from=builder /go/bin/app /go/bin/app

CMD ./go/bin/app