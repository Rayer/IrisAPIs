FROM golang:alpine as build-env
RUN apk --no-cache add tzdata
WORKDIR /IrisAPIs
ADD . /IrisAPIs
RUN cd /IrisAPIs/server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server.app
RUN cd /IrisAPIs/serverInfo && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o serverInfo.app

FROM scratch
WORKDIR /app
COPY --from=build-env /IrisAPIs/server/server.app /app
COPY --from=build-env /IrisAPIs/serverInfo/serverInfo.app /app
COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Taipei
EXPOSE 8080
ENTRYPOINT ["/app/server.app"]
