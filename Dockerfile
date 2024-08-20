# syntax=docker/dockerfile:1

FROM golang:alpine3.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE 8082

CMD ["/docker-gs-ping"]
