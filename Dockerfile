FROM golang:latest
LABEL maintainer avvero

ADD . /app
WORKDIR /app
RUN go test *.go
RUN go build -o main .
CMD ["/app/main"]