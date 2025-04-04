FROM golang:1.24 AS builder

WORKDIR /auth

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 80

RUN go build -o /go/bin/auth-app ./cmd/main.go

FROM alpine:3.20

WORKDIR /auth

COPY --from=builder /auth/internal/migration /migration

COPY --from=builder /go/bin/auth-app /go/bin/auth-app

CMD ./go/bin/auth-app