FROM golang:1.23

WORKDIR /go/app
RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]
