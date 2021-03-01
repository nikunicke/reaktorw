FROM golang:1.15.5-alpine3.12
WORKDIR /usr/src/app
COPY . .
RUN go mod init github.com/nikunicke/reaktorw
RUN go build ./cmd/reaktorw
CMD ["./reaktor-warehouse"]