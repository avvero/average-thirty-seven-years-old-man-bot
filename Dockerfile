####
# Build image
####
ARG GOLANG_VER=latest
FROM golang:latest AS build
LABEL maintainer avvero

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd/bot/main.go

####
# Runtime image
####
FROM scratch
LABEL maintainer avvero

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/main /app/main

ENTRYPOINT [ "/app/main" ]
