FROM golang:1.14.14-stretch
LABEL maintainer avvero

ADD . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main", "-httpPort=8080"]