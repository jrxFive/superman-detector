FROM golang:1.13-stretch as tester

WORKDIR /go/src/github.com/jrxfive/superman-detector/
COPY . .

# go-sqlite3 requires cgo to work
# go race with alpine issue https://github.com/golang/go/issues/14481
RUN go test -race -cover ./...

FROM golang:1.13-alpine as builder
WORKDIR /go/src/github.com/jrxfive/superman-detector/
COPY --from=tester go/src/github.com/jrxfive/superman-detector .

RUN apk update \
    && apk add --no-cache gcc musl-dev \
    && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o superman-detector

FROM golang:1.13-alpine
ARG SERVICE_PORT=8080

ENV DETECTOR_API_SERVICE_PORT=$SERVICE_PORT
EXPOSE $SERVICE_PORT

RUN apk update \
    && apk add --no-cache ca-certificates curl \
    && rm -rf /var/cache/apk/*

WORKDIR /app/

COPY --from=builder /go/src/github.com/jrxfive/superman-detector/superman-detector .
COPY --from=builder /go/src/github.com/jrxfive/superman-detector/GeoLite2-City.mmdb .

CMD ["./superman-detector"]