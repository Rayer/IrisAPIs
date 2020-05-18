FROM golang:alpine as build-env
WORKDIR /IrisAPIs
ADD . /IrisAPIs
RUN cd /IrisAPIs/server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server.app

FROM scratch
WORKDIR /app
COPY --from=build-env /IrisAPIs/server/server.app /app
EXPOSE 8080
ENTRYPOINT ["/app/server.app"]
