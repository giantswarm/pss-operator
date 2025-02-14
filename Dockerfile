FROM alpine:3.21.3

RUN apk add --no-cache ca-certificates

ADD ./pss-operator /pss-operator

ENTRYPOINT ["/pss-operator"]
