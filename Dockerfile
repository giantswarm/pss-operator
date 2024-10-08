FROM alpine:3.20.3

RUN apk add --no-cache ca-certificates

ADD ./pss-operator /pss-operator

ENTRYPOINT ["/pss-operator"]
