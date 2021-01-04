FROM golang:alpine as build-env
RUN apk --no-cache add tzdata
RUN apk add build-base
RUN apk add -U --no-cache ca-certificates
WORKDIR /IrisAPIs
ADD . /IrisAPIs
RUN cd /IrisAPIs/server && go generate
RUN cd /IrisAPIs/server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server.app
RUN cd /IrisAPIs/serverInfo && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o serverInfo.app
RUN cd /IrisAPIs/apikey_cli && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o apikey_cli.app

FROM scratch
WORKDIR /app
COPY --from=build-env /IrisAPIs/server/server.app /app
COPY --from=build-env /IrisAPIs/server/docs/ /app/docs
COPY --from=build-env /IrisAPIs/serverInfo/serverInfo.app /app
COPY --from=build-env /IrisAPIs/apikey_cli/apikey_cli.app /app
COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV TZ=Asia/Taipei
EXPOSE 8080
EXPOSE 8082
ENTRYPOINT ["/app/server.app"]
