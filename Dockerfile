FROM golang:1.15 as builder

ADD . /${GOPATH}/src/github.com/Drakkar-Software/Metrics-Forwarder
WORKDIR /${GOPATH}/src/github.com/Drakkar-Software/Metrics-Forwarder

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o server .

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/Drakkar-Software/Metrics-Forwarder/server /server/
COPY --from=builder /go/src/github.com/Drakkar-Software/Metrics-Forwarder/docker /docker/
WORKDIR /server

EXPOSE 8080

# Start up
CMD ["/docker/startup.sh"]
