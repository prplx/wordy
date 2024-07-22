# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . /app/.

RUN CGO_ENABLED=0 GOOS=linux go build -o /wordy ./cmd/app/main.go

EXPOSE 3000

CMD ["/wordy"]
