####
# Build image
####
ARG GOLANG_VER=latest
FROM golang:latest AS build
LABEL maintainer avvero

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test *.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main

####
# Runtime image
####
FROM scratch
LABEL maintainer avvero

COPY --from=build /app/main /app/main

ENTRYPOINT [ "/app/main" ]