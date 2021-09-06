FROM golang:1.17

WORKDIR /go/src/gitlab.com/tokend/traefik-cop
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/traefik-cop gitlab.com/tokend/traefik-cop

###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/traefik-cop /usr/local/bin/traefik-cop
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["traefik-cop"]
